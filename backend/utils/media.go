package utils

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iafan/cwalk"
	"github.com/orellazri/photolens/models"
)

type media struct {
	contentType  string
	isPhoto      bool
	lastModified time.Time
}

func ProcessMedia(context *Context) error {
	log.Print("Starting to process media files...")
	start := time.Now()

	// Walk the fsMedia directory and get all photo names
	fsMedia := make(map[string]media, 0) // Map photo path to last modified time
	err := cwalk.Walk(context.RootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if this is a file and not a directory
			if !info.IsDir() {
				// Check content type by reading the first 512 bytes of the file
				file, err := os.Open(filepath.Join(context.RootPath, path))
				if err != nil {
					return err
				}
				buff := make([]byte, 512)
				if _, err := file.Read(buff); err != nil {
					return err
				}
				contentType := http.DetectContentType(buff)

				// Add to filesystem map
				fsMedia[path] = media{
					contentType:  contentType,
					isPhoto:      strings.HasPrefix(contentType, "image/"),
					lastModified: info.ModTime().UTC(),
				}
			}
			return nil
		})
	if err != nil {
		return err
	}

	numProcessed := len(fsMedia)
	numIndexed := 0

	// Iterate through all photos in the database and remove from
	// the filesystem map if their last modified time is equal
	var results []models.Media
	context.DB.Select("path", "last_modified").Find(&results)
	for _, result := range results {
		if _, ok := fsMedia[result.Path]; ok {
			if fsMedia[result.Path].lastModified == result.LastModified {
				delete(fsMedia, result.Path)
			}
		}
	}

	// Now the filesystem map contains photos that are either not in
	// the database, or their last modified times are different.
	// So we need to sync them
	for path, media := range fsMedia {
		// TODO: Index media file (Generate thumbnails, etc.)

		// Try to create photo in database, or update last modified time if
		// it already exists
		photo := models.Media{
			Path:         path,
			IsPhoto:      media.isPhoto,
			ContentType:  media.contentType,
			LastModified: media.lastModified,
		}
		if context.DB.Model(&photo).Where("path = ?", path).Updates(&photo).RowsAffected == 0 {
			context.DB.Create(&photo)
		}

		numIndexed += 1
	}

	// TODO: Compare photos left in database that are not in the filesystem
	// and remove them (from database, thumbnails, etc.)

	log.Printf("Processing done in %v\n", time.Since(start))
	log.Printf("    Processed: %v\n", numProcessed)
	log.Printf("    Indexed: %v\n", numIndexed)

	return nil
}
