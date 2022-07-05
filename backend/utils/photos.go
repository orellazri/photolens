package utils

import (
	"log"
	"os"
	"time"

	"github.com/iafan/cwalk"
	"github.com/orellazri/photolens/models"
)

func IndexPhotos(context *Context) error {
	log.Print("Starting to index photos...")
	start := time.Now()

	// Walk the photos directory and get all photo names
	photos := make(map[string]time.Time, 0) // Map photo path to last modified time
	err := cwalk.Walk(context.RootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if this is a file and not a directory
			if !info.IsDir() {
				photos[path] = info.ModTime().UTC()
			}
			return nil
		})
	if err != nil {
		return err
	}

	numPhotos := len(photos)

	// Iterate through all photos in the database and remove from
	// the filesystem map if their last modified time is equal
	var results []models.Photo
	context.DB.Select("path", "last_modified").Find(&results)
	for _, result := range results {
		if _, ok := photos[result.Path]; ok {
			if photos[result.Path] == result.LastModified {
				delete(photos, result.Path)
			}
		}
	}

	// Now the filesystem map contains photos that are either not in
	// the database, or their last modified times are different.
	// So we need to sync them
	for path, lastModified := range photos {
		// TODO: Sync photo (Generate thumbnails, etc.)

		// Try to create photo in database, or update last modified time if
		// it already exists
		photo := models.Photo{
			Path:         path,
			LastModified: lastModified,
		}
		if context.DB.Model(&photo).Where("path = ?", path).Updates(&photo).RowsAffected == 0 {
			context.DB.Create(&photo)
		}
	}

	// TODO: Compare photos left in database that are not in the filesystem
	// and remove them (from database, thumbnails, etc.)

	log.Printf("Indexed %v photos in %v\n", numPhotos, time.Since(start))

	return nil
}
