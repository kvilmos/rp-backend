package observer

import (
	"context"
	"encoding/json"
	"log"
	"room-planner/common/constant"
	"room-planner/service"
	"time"

	"github.com/redis/go-redis/v9"
)

type Observer struct {
	RedisClient      *redis.Client
	FurnitureService service.FurnitureService
}

func New(client *redis.Client, fs service.FurnitureService) *Observer {
	return &Observer{
		RedisClient:      client,
		FurnitureService: fs,
	}
}

func (o *Observer) Start(ctx context.Context) {
	go o.processFileUploadQueue(ctx)
}

func (o *Observer) processFileUploadQueue(ctx context.Context) {
	for {
		result, err := o.RedisClient.BRPop(ctx, 0, constant.UPLOAD_QUEUE_NAME).Result()
		if err != nil {
			// TODO STORY-201 ERROR Handler
			time.Sleep(5 * time.Second)
			continue
		}
		var notification service.UploadNotification
		if err := json.Unmarshal([]byte(result[1]), &notification); err != nil {
			// TODO STORY-201 ERROR Handler
			continue
		}

		if err := o.FurnitureService.FinalizeUpload(ctx, notification); err != nil {
			log.Printf("ERROR: failed to process upload for key %s: %v", notification.ObjectKey, err)
			// Itt lehet dönteni az újrapróbálkozásról (requeue) a service-től kapott hiba alapján.
		}
	}
}
