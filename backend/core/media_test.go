package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/orellazri/photolens/database"
	"github.com/orellazri/photolens/models"
)

func setup(t *testing.T) (*Context, string) {
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
	return &context, tempDir
}

func cleanup(t *testing.T, context *Context) {
	sqlDB, err := context.DB.DB()
	err = sqlDB.Close()
	if err != nil {
		t.Error(err)
	}
}

func copySampleToTestDir(t *testing.T, context *Context, testDir string, file string) {
	// Copy sample photos to test directory
	input, err := ioutil.ReadFile(filepath.Join("../", "samples", file))
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(filepath.Join(context.RootPath, file), input, 0644)
	if err != nil {
		t.Error(err)
	}
}

func TestProcessMediaAddFiles(t *testing.T) {
	context, testDir := setup(t)

	// Copy sample files
	sampleFiles := []string{"sample1.png", "sample2.jpg", "sample3.bmp", "sample4.gif"}
	for _, file := range sampleFiles {
		copySampleToTestDir(t, context, testDir, file)
	}

	// Process media files
	err := ProcessMedia(context)
	if err != nil {
		t.Error(err)
	}

	// Check database for media files
	var results []models.Media
	err = context.DB.Select("id").Find(&results).Error
	if err != nil {
		t.Error(err)
	}
	if len(results) != 4 {
		t.Errorf("%d media files found in database, want 4", len(results))
	}

	cleanup(t, context)
}

func TestProcessMediaRemoveFiles(t *testing.T) {
	context, testDir := setup(t)

	// Copy sample files
	sampleFiles := []string{"sample1.png", "sample2.jpg"}
	for _, file := range sampleFiles {
		copySampleToTestDir(t, context, testDir, file)
	}

	// Process media files
	err := ProcessMedia(context)
	if err != nil {
		t.Error(err)
	}

	// Check database for media files
	var results []models.Media
	err = context.DB.Select("id").Find(&results).Error
	if err != nil {
		t.Error(err)
	}
	if len(results) != 2 {
		t.Errorf("%d media files found in database, want 2", len(results))
	}

	// Remove one file to see that the database updates
	err = os.Remove(filepath.Join(context.RootPath, "sample2.jpg"))
	if err != nil {
		t.Error(err)
	}

	// Process media files
	err = ProcessMedia(context)
	if err != nil {
		t.Error(err)
	}

	err = context.DB.Select("id").Find(&results).Error
	if err != nil {
		t.Error(err)
	}
	if len(results) != 1 {
		t.Errorf("%d media files found in database, want 1", len(results))
	}

	cleanup(t, context)
}
