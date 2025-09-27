package observer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"room-planner/app"
	"room-planner/common/constant"
	"room-planner/domain/furniture"
	"room-planner/handler"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Observer struct {
	app *app.Application
}

func NewObserver(app *app.Application) *Observer {
	return &Observer{
		app: app,
	}
}

func (o *Observer) Start(ctx context.Context) {
	go o.processFileUploadQueue(ctx)
}

func (o *Observer) processFileUploadQueue(ctx context.Context) {
	for {
		result, err := o.app.RedisClient.BRPop(ctx, 0, constant.UPLOAD_QUEUE_NAME).Result()
		if err != nil {
			// TODO STORY-201 ERROR Handler
			time.Sleep(5 * time.Second)
			continue
		}

		payload := result[1]
		o.handleFurnitureFileUpload(ctx, payload)
	}
}

func (o *Observer) handleFurnitureFileUpload(ctx context.Context, payload string) {
	notification := handler.UploadNotification{}
	s3ResponseBytes := []byte(payload)
	err := json.Unmarshal(s3ResponseBytes, &notification)
	if err != nil {
		// TODO STORY-201 ERROR HANDLER
		return
	}

	furnitureKey := notification.ObjectKey
	lockKey := constant.DISTRIBUTED_LOCK_PREFIX + furnitureKey
	lockValue, err := acquireLock(ctx, o.app.RedisClient, lockKey, constant.LOCK_TTL)

	if err != nil {
		if err != ErrLockNotAcquired {
			// TODO STORY-201 ERROR HANDLER
			// RETRY
			return
		}

		errRequeue := o.app.RedisClient.LPush(context.Background(), constant.UPLOAD_QUEUE_NAME, payload).Err()
		if errRequeue != nil {
			//TODO STORY-201 ERROR HANDLER
			// RETRY
			return
		}

		return
	}

	defer func() {
		err := releaseLock(ctx, o.app.RedisClient, lockKey, lockValue)
		if err != nil {
			//TODO STORY-201 ERROR HANDLER
		}
	}()

	uploadStatus := handler.FurnitureUploadStatus{}
	val, err := o.app.RedisClient.Get(ctx, furnitureKey).Result()
	if err == redis.Nil {
		// TODO STORY-201 ERROR HANDLER
		// TERMINATE
		return
	}

	if err != nil {
		// TODO STORY-201 ERROR HANDLER
		// RETRY
		return
	}

	err = json.Unmarshal([]byte(val), &uploadStatus)
	if err != nil {
		// TODO STORY-201 ERROR HANDLER
		// TERMINATE
		return
	}

	switch notification.BucketName {
	case constant.FURNITURE_BUCKET:
		uploadStatus.IsObjectUploaded = true
	case constant.THUMBNAIL_BUCKET:
		uploadStatus.IsThumbnailUploaded = true
	}

	if uploadStatus.IsObjectUploaded && uploadStatus.IsThumbnailUploaded {
		exists, err := o.isFurnitureInDB(ctx, furnitureKey)
		if err != nil {
			// TODO STORY-201 ERROR HANDLER
			// RETRY DB INSERT
			return
		}

		if exists {
			log.Printf("INFO: Furniture with key %s already exists in DB. Skipping insert.", furnitureKey)
		} else {
			fileName, _ := uuid.Parse(furnitureKey)

			newFurniture := furniture.NewFurniture{
				Name:     uploadStatus.Furniture.Name,
				FileName: fileName,
			}

			tx := o.app.Db.Begin()

			err := furniture.Create(tx, newFurniture)
			if err != nil {
				// TODO STORY-201 ERROR HANDLER
				// RETRY DB INSERT
				tx.Rollback()
				return
			}
			tx.Commit()
		}

		err = o.app.RedisClient.Del(ctx, furnitureKey).Err()
		if err != nil {
			// TODO STORY-201 ERROR HANDLER
			// TERMINATE
			return
		}

		fmt.Println("FINALE", uploadStatus)
		return
	}

	uploadStatusJson, err := json.Marshal(&uploadStatus)
	if err != nil {
		// TODO STORY-201 ERROR HANDLER
		// TERMINATE
		return
	}

	err = o.app.RedisClient.Set(ctx, furnitureKey, uploadStatusJson, constant.UPLOAD_TTL).Err()
	if err != nil {
		// TODO STORY-201 ERROR HANDLER
		// RETRY
		return
	}

	fmt.Println("FINALE", uploadStatus)
}

func (o *Observer) isFurnitureInDB(ctx context.Context, key string) (bool, error) {
	return false, nil
}
