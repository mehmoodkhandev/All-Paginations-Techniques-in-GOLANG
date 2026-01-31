package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectingMySQLDb(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUsername, config.DBPassword, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	log.Println("Connected successfully to the database (MySQL)")
	db.Logger = logger.Default.LogMode(logger.Info)
	return db, nil
}
