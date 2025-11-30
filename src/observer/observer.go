package observer

import (
	"context"
	"encoding/json"
	"log/slog"
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
			time.Sleep(5 * time.Second)
			continue
		}
		var notification service.UploadNotification
		if err := json.Unmarshal([]byte(result[1]), &notification); err != nil {
			continue
		}

		if err := o.FurnitureService.FinalizeUpload(ctx, notification); err != nil {
			slog.Error("ERROR: failed to process upload for key %s: %v", notification.ObjectKey, err)
		}
	}
}
