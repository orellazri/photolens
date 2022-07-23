package database

import (
	"github.com/orellazri/photolens/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err
}

func MigrateDatabase(db *gorm.DB) error {
	var models []interface{} = []interface{}{
		&models.Media{},
		&models.Meta{},
	}

	return db.AutoMigrate(models...)
}
