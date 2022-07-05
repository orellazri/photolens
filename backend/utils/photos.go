package utils

import (
	"log"
	"os"
	"time"

	"github.com/iafan/cwalk"
)

func IndexPhotos(context *Context) error {
	log.Print("Starting to index photos...")
	start := time.Now()

	// Walk the photos directory and get all photo names
	photos := make([]string, 0)
	err := cwalk.Walk(context.RootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if this is a file and not a directory
			if !info.IsDir() {
				photos = append(photos, path)
			}
			return nil
		})
	if err != nil {
		return err
	}

	log.Printf("Indexed %v photos in %v\n", len(photos), time.Since(start))

	return nil
}
