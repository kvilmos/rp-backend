package storage

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQLConnection() (*gorm.DB, error) {
	connectionStr := os.Getenv("MYSQL_STR")

	db, err := gorm.Open(mysql.Open(connectionStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
