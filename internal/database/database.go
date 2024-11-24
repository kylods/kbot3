package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

