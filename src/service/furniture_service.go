package service

import (
	"context"
	"fmt"
	"log"
	"room-planner/common/constant"
	"room-planner/model"
	"room-planner/repository"
	"room-planner/request"

	"github.com/google/uuid"
)

type FurnitureService struct {
	FurnitureRepo repository.FurnitureRepository
	FileStore     repository.FileStore
	Cache         repository.Cache
	Queue         repository.Queue
	Locker        repository.Locker
}

func NewFurnitureService(fr repository.FurnitureRepository, fs repository.FileStore, c repository.Cache, q repository.Queue, l repository.Locker) *FurnitureService {
	return &FurnitureService{
		FurnitureRepo: fr,
		FileStore:     fs,
		Cache:         c,
		Queue:         q,
		Locker:        l,
	}
}

type FurnitureUploadStatus struct {
	FurnitureName       string `json:"furnitureName"`
	IsThumbnailUploaded bool   `json:"isThumbnailUploaded"`
	IsObjectUploaded    bool   `json:"isObjectUploaded"`
}

type FurnitureUploadLinks struct {
	ObjectUrl    string `json:"objectUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type UploadNotification struct {
	BucketName string
	ObjectKey  string
	RetryCount int
}

func (s FurnitureService) PrepareUpload(ctx context.Context, furnitureName string) (*FurnitureUploadLinks, error) {
	fileId := uuid.New().String()
	objectUrl, err := s.FileStore.GenerateSignedUploadUrl(ctx, constant.FURNITURE_BUCKET, fileId, constant.SIGNATURE_TTL)
	if err != nil {
		return nil, err
	}
	thumbnailUrl, err := s.FileStore.GenerateSignedUploadUrl(ctx, constant.THUMBNAIL_BUCKET, fileId, constant.SIGNATURE_TTL)
	if err != nil {
		return nil, err
	}

	uploadStatus := FurnitureUploadStatus{
		FurnitureName: furnitureName,
	}

	if err := s.Cache.Set(ctx, fileId, uploadStatus, constant.UPLOAD_TTL); err != nil {
		return nil, err
	}

	return &FurnitureUploadLinks{
		ObjectUrl:    objectUrl.String(),
		ThumbnailUrl: thumbnailUrl.String(),
	}, nil

}

func (s FurnitureService) QueueUploadNotification(ctx context.Context, notification request.MinioNotification) error {
	fmt.Println("IN SERVICE - QUEUE")
	if len(notification.Records) < 1 {
		return fmt.Errorf("webhook notification contained no records")
	}

	for _, record := range notification.Records {
		internalNotification := UploadNotification{
			BucketName: record.S3.Bucket.Name,
			ObjectKey:  record.S3.Object.Key,
		}
		fmt.Println("IN SERVICE - ", internalNotification)

		err := s.Queue.Enqueue(ctx, constant.UPLOAD_QUEUE_NAME, internalNotification)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s FurnitureService) FinalizeUpload(ctx context.Context, notification UploadNotification) error {
	fileKey := notification.ObjectKey
	lockKey := constant.DISTRIBUTED_LOCK_PREFIX + fileKey
	lockValue, err := s.Locker.Acquire(ctx, lockKey, constant.LOCK_TTL)
	if err != nil {
		// REQUEUE
		// TODO STORY-201 ERROR HANDLER
		return err
	}

	defer func() {
		err := s.Locker.Release(ctx, lockKey, lockValue)
		if err != nil {
			// REQUEUE
			// TODO STORY-201 ERROR HANDLER
		}
	}()

	uploadStatus := FurnitureUploadStatus{}
	err = s.Cache.Get(ctx, fileKey, &uploadStatus)
	if err != nil {
		// REQUEUE
		// TODO STORY-201 ERROR HANDLER
		return err
	}

	switch notification.BucketName {
	case constant.FURNITURE_BUCKET:
		uploadStatus.IsObjectUploaded = true
	case constant.THUMBNAIL_BUCKET:
		uploadStatus.IsThumbnailUploaded = true
	}

	if uploadStatus.IsObjectUploaded && uploadStatus.IsThumbnailUploaded {
		fileId, err := uuid.Parse(fileKey)
		if err != nil {
			// TERMINATE
			// TODO STORY-201 ERROR HANDLER
			return err
		}

		exists, err := s.FurnitureRepo.IsExistByFileId(fileId)
		if err != nil {
			// RETRY DB INSERT
			// TODO STORY-201 ERROR HANDLER
			return err
		}

		if exists {
			log.Printf("INFO: Furniture with key %s already exists in DB. Skipping insert.", fileId)
		} else {
			newFurniture := model.Furniture{
				Name:     uploadStatus.FurnitureName,
				FileName: fileId,
			}
			err := s.FurnitureRepo.Create(ctx, &newFurniture)
			if err != nil {
				// RETRY DB INSERT
				// TODO STORY-201 ERROR HANDLER
				return err
			}
		}

		err = s.Cache.Delete(ctx, fileKey)
		if err != nil {
			// RETRY TERMINATE
			// TODO STORY-201 ERROR HANDLER
			return err
		}

		return nil
	}

	err = s.Cache.Set(ctx, fileKey, uploadStatus, constant.UPLOAD_TTL)
	if err != nil {
		// REQUEUE
		// TODO STORY-201 ERROR HANDLER
		return err
	}

	return nil
}
