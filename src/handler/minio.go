package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"room-planner/common/constant"

	"github.com/labstack/echo/v4"
)

type MinioNotification struct {
	Records []Record `json:"Records"`
}

type Record struct {
	EventName string `json:"eventName"`
	S3        S3     `json:"s3"`
}

type S3 struct {
	Bucket Bucket `json:"bucket"`
	Object Object `json:"object"`
}

type Bucket struct {
	Name string `json:"name"`
}

type Object struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
}

type UploadNotification struct {
	BucketName string
	ObjectKey  string
	RetryCount int
}

func (h Handler) UploadFurnitureFile(c echo.Context) error {
	var notification MinioNotification
	if err := c.Bind(&notification); err != nil {
		// TODO: STORY-201 - ERROR Handler
		log.Fatal(err)
		return err
	}
	fmt.Println("HOOK", notification, len(notification.Records))
	if len(notification.Records) < 1 {
		// TODO: STORY-201 - ERROR Handler
		return fmt.Errorf("")
	}

	for _, record := range notification.Records {

		uploadNotification := UploadNotification{
			BucketName: record.S3.Bucket.Name,
			ObjectKey:  record.S3.Object.Key,
		}

		message, err := json.Marshal(uploadNotification)
		if err != nil {
			// TODO: STORY-201 - ERROR Handler
			fmt.Println(err)
		}

		err = h.redisClient.LPush(context.Background(), constant.UPLOAD_QUEUE_NAME, message).Err()
		if err != nil {
			// TODO: STORY-201 - ERROR Handler
			fmt.Println(err)
		}
	}

	return c.NoContent(http.StatusOK)
}
