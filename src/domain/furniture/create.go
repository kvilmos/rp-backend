package furniture

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewFurniture struct {
	Id       int
	Name     string
	FileName uuid.UUID
}

func (NewFurniture) TableName() string {
	return "furniture_t"
}

func Create(tx *gorm.DB, furniture NewFurniture) error {
	err := tx.Omit("id", "thumbnail_path", "object_path").Create(&furniture).Error
	if err != nil {
		return err
	}

	return nil
}
