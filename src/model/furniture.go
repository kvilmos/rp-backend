package model

import (
	"time"

	"github.com/google/uuid"
)

type Furniture struct {
	Id         int64
	Name       string
	UserId     int64
	CategoryId int8
	FileName   uuid.UUID
	CreatedAt  time.Time
}

func (Furniture) TableName() string {
	return "furniture_t"
}

type FurnitureUploadStatus struct {
	Furniture           Furniture
	IsThumbnailUploaded bool
	IsObjectUploaded    bool
}
