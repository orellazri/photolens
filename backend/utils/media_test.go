package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/photolens/database"
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

	cleanup(t, context)
}
