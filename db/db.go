package db

import (
	"go-web-app/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	// Migrate schema
	DB.AutoMigrate(&models.Note{}, &models.User{})

	// Insert test user if not exists
	var count int64
	DB.Model(&models.User{}).Where("email = ?", "user@example.com").Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
		DB.Create(&models.User{
			Email:    "user@example.com",
			Password: string(hashedPassword),
		})
	}
}
