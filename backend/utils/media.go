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
	shouldIndex  bool
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
				defer file.Close()
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
					shouldIndex:  true,
				}
			}
			return nil
		})
	if err != nil {
		return err
	}

	numProcessed := len(fsMedia)
	numIndexed := 0 // Number of media files that are going to be indexed
	numDeleted := 0 // Number of media files that were deleted from the filesystem

	// Iterate through all media files in the database and compare their last
	// modified tomes with the ones in the filesystem. If they are not equal,
	// mark those files as "should index" to index them
	var results []models.Media
	err = context.DB.Select("id", "path", "last_modified").Find(&results).Error
	if err != nil {
		return err
	}
	for _, result := range results {
		if media, ok := fsMedia[result.Path]; ok {
			if media.lastModified == result.LastModified {
				media.shouldIndex = false
				fsMedia[result.Path] = media
			}
		}
	}

	// Iterate through all the media files in the filesystem map
	// and index them if needed
	for path, media := range fsMedia {
		if !media.shouldIndex {
			continue
		}

		// TODO: Proprely index file (generate thumbnails, etc.)

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

	// Iterate through all the media files in the database and check if they don't
	// exist in the filesystem. If so, delete them
	for _, result := range results {
		if _, ok := fsMedia[result.Path]; !ok {
			err = context.DB.Unscoped().Delete(&models.Media{}, result.ID).Error
			if err != nil {
				return err
			}

			// TODO: Properly delete remains (thumbnails, etc.)

			numDeleted += 1
		}
	}

	log.Printf("    Processed: %v\n", numProcessed)
	log.Printf("    Indexed: %v\n", numIndexed)
	log.Printf("    Deleted: %v\n", numDeleted)
	log.Printf("Processing done in %v\n", time.Since(start))

	return nil
}
