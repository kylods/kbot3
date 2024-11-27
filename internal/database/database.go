package database

import (
	"log"

	"github.com/kylods/kbot-backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(connectionString string) (*gorm.DB, error) {
	log.Println("Opening database connection")
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.AudioFile{}, &models.Guild{})
	if err != nil {
		return nil, err
	}

	db.Create(&models.Guild{Name: "Midtest Devout", DjRoles: "", LoopEnabled: false})

	return db, nil
}
