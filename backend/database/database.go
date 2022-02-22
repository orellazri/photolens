package database

import (
	"github.com/orellazri/photolens/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, err
}

func MigrateDatabase(db *gorm.DB) error {
	var models []interface{} = []interface{}{
		&models.Photo{},
	}

	db.AutoMigrate(models...)

	return nil
}
