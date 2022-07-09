package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/photolens/database"
	"github.com/orellazri/photolens/models"
)

func setup(t *testing.T) *Context {
	tempDir := t.TempDir()
	os.Mkdir(filepath.Join(tempDir, "media"), os.ModePerm)
	os.Mkdir(filepath.Join(tempDir, "cache"), os.ModePerm)

	db, err := database.SetupDatabase(filepath.Join(tempDir, "photolens.db"))
	if err != nil {
		t.Error(err)
	}
	err = database.MigrateDatabase(db)
	if err != nil {
		t.Error(err)
	}
	context := Context{
		DB:        db,
		RootPath:  filepath.Join(tempDir, "media"),
		CachePath: filepath.Join(tempDir, "cache"),
	}
	return &context
}

func cleanup(t *testing.T, context *Context) {
	sqlDB, err := context.DB.DB()
	err = sqlDB.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestProcessMediaAddFiles(t *testing.T) {
	context := setup(t)

	err := ProcessMedia(context)
	if err != nil {
		t.Error(err)
	}

	var results []models.Media
	err = context.DB.Select("id").Find(&results).Error
	if err != nil {
		t.Error(err)
	}

	if len(results) != 0 {
		t.Errorf("%d media file found in database, want 0", len(results))
	}

	cleanup(t, context)
}
