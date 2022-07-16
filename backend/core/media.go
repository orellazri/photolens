package core

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/orellazri/photolens/models"
	"golang.org/x/exp/slices"
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

	// Get last process time from database
	var lastProcessTimeResult models.Meta
	var lastProcessTime time.Time
	err := context.DB.Where("key = ?", "last_process_time").First(&lastProcessTimeResult).Error
	if err != nil {
		if err.Error() != "record not found" {
			return err
		}

		// We never processed before, so set the last process time to 100 years ago
		log.Println("First processing detected. Adding process time to database")
		lastProcessTime = time.Now().AddDate(-100, 0, 0).UTC()
		err = context.DB.Create(&models.Meta{
			Key:   "last_process_time",
			Value: lastProcessTime.String(),
		}).Error
		if err != nil {
			return err
		}
	} else {
		// Set last process time from database if it exists there already
		lastProcessTime, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", lastProcessTimeResult.Value)
		if err != nil {
			return err
		}
	}

	// Walk the fsMedia directory and get all photo names
	fsMedia := make(map[string]media, 0) // Map photo path to last modified time
	// TODO: Look into walking in parallel
	err = filepath.Walk(context.RootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				// This is a directory.
				// Check if its last modified time is greather than the last processing time
				// from the database
				if info.ModTime().UTC().After(lastProcessTime) {
					// This directory was modified after the last process time, so we need to
					// process it
					// TODO: Only process directories by last modified time
				}
			} else {
				// This is a file.

				// Check if we should ignore this file
				if slices.Contains(FilesToIgnore, filepath.Base(path)) {
					return nil
				}

				// Check content type by reading the first 512 bytes of the file
				file, err := os.Open(path)
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

		// Generate thumbnail in the background
		// TODO: Error handling (wait for all thumbnails to finish at the end of processing
		// 		 and check for errors)
		go GetThumbnail(context, &models.Media{Path: path})

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

			err = os.Remove(getThumbnailPath(context, &result))
			if err != nil && !strings.Contains(err.Error(), "no such file") {
				return err
			}

			numDeleted += 1
		}
	}

	// Set last process time
	lastProcessTime = time.Now().UTC()
	err = context.DB.Model(&models.Meta{}).Where("key = ?", "last_process_time").Update("value", lastProcessTime.String()).Error
	if err != nil {
		return err
	}

	log.Printf("    Processed: %v\n", numProcessed)
	log.Printf("    Indexed: %v\n", numIndexed)
	log.Printf("    Deleted: %v\n", numDeleted)
	log.Printf("Processing done in %v\n", time.Since(start))

	return nil
}

func GetMediaFromID(id int, context *Context) (*models.Media, error) {
	// Query media from id in database
	var media models.Media
	err := context.DB.First(&media, id).Error
	if err != nil {
		return nil, err
	}

	return &media, nil
}

// Returns the path of the thumbnail file for the given media file
func getThumbnailPath(context *Context, media *models.Media) string {
	return filepath.Join(context.CachePath, "thumbnails", fmt.Sprintf("%s.png", media.Path))
}

// Fetch an existing thumbnail or generate a new one if it doesn't exist already
// Returns the thumnail image in a base64 encoded string
func GetThumbnail(context *Context, media *models.Media) (string, error) {
	thumbnailPath := getThumbnailPath(context, media)

	// Check if thumbnail already exists before generating new one
	if _, err := os.Stat(thumbnailPath); err == nil {
		content, err := ioutil.ReadFile(thumbnailPath)
		if err != nil {
			return "", err
		}
		encoded := base64.StdEncoding.EncodeToString(content)
		return encoded, nil
	}

	// If we got here, the thumbnail doesn't exist already
	// so we generate one and return it
	return generateThumbnail(context, media, thumbnailPath)
}

// Generate a thumbnail for a given image
// Returns the thumbnail image in a base64 encoded string
func generateThumbnail(context *Context, media *models.Media, thumbnailPath string) (string, error) {
	// Open original image
	file, err := os.Open(media.Path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode original image
	image, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Create directories for thumbnail according to original media file's path
	err = os.MkdirAll(filepath.Join(context.CachePath, "thumbnails", filepath.Dir(media.Path)), os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create thumbnail file
	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "", err
	}

	// Resize original image and write to thumbnail file
	resizedImage := imaging.Fill(image, 128, 128, imaging.Center, imaging.Lanczos)
	err = png.Encode(thumbnailFile, resizedImage)
	if err != nil {
		return "", err
	}
	thumbnailFile.Close()

	content, err := ioutil.ReadFile(thumbnailPath)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}
