package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"room-planner/common/constant"
	"room-planner/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type NewFurnitureRequest struct {
	Name string
}

type FurnitureUploadStatus struct {
	Furniture           NewFurnitureRequest
	IsThumbnailUploaded bool
	IsObjectUploaded    bool
}

func (h Handler) CreateFurniture(c echo.Context) error {
	var request NewFurnitureRequest
	err := c.Bind(&request)
	if err != nil {
		// TODO: STORY-201 - ERROR Handler
		log.Fatal(err)
		return err
	}

	fileName := uuid.New()
	fileKey := fileName.String()

	objectUrl, err := h.minioClient.PresignedPutObject(context.Background(), constant.FURNITURE_BUCKET, fileKey, constant.SIGNATURE_TTL)
	if err != nil {
		// TODO: STORY-201 - ERROR Handler
		return err
	}
	thumbnailUrl, err := h.minioClient.PresignedPutObject(context.Background(), constant.THUMBNAIL_BUCKET, fileKey, constant.SIGNATURE_TTL)
	if err != nil {
		// TODO: STORY-201 - ERROR Handler
		return err
	}

	furniture := NewFurnitureRequest{
		Name: request.Name,
	}
	uploadStatus := FurnitureUploadStatus{
		Furniture: furniture,
	}
	furnitureJson, err := json.Marshal(uploadStatus)
	if err != nil {
		// TODO: STORY-201 - ERROR Handler
		return err
	}
	res := h.redisClient.Set(context.Background(), fileKey, furnitureJson, 0)
	fmt.Println(res)

	response := dto.UploadUrlsDto{
		ObjectUrl:    objectUrl.String(),
		ThumbnailUrl: thumbnailUrl.String(),
	}

	return c.JSON(http.StatusOK, response)
}
