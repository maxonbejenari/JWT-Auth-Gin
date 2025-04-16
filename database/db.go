package database

import (
	"fmt"
	"github.com/maxonbejenari/gin-auth/models"
	"log"

	"github.com/maxonbejenari/gin-auth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can not connect to db %v", err)
	}

	DB = db
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to run migrations")
	}

	log.Println("Successful connected to Database")
}
