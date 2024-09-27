package db

import (
	"keno/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DbKey = "db"

func SetupDatabase(file string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Game{}, &models.Card{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
