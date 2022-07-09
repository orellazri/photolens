package database

import (
	"github.com/orellazri/photolens/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("photolens.db"), &gorm.Config{})
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

	db.AutoMigrate(models...)

	return nil
}
