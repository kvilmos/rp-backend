package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var instance DatabaseContext

type DatabaseContext struct {
	db *gorm.DB
}

func Connect() error {
	connection := "develop:develop@tcp(127.0.0.1:3306)/room-planner?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if err != nil {
		return err
	}
	instance = DatabaseContext{db: db}

	return nil
}

func GetTransaction() *gorm.DB {
	return instance.db.Begin()
}
