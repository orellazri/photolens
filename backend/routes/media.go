package routes

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/orellazri/photolens/utils"
)

func RegisterMediaRouter(context *utils.Context, router *mux.Router) {
	router.HandleFunc("/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { getMedia(w, r, context) }).Methods("GET")
	router.HandleFunc("/thumbnail/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { getThumbnail(w, r, context) }).Methods("GET")
}

func getMedia(w http.ResponseWriter, r *http.Request, context *utils.Context) {
	// Convert id to number
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, fmt.Sprintf("Invalid id %v", id))
		return
	}

	// Get media from ID
	media, err := utils.GetMediaFromID(id, context)
	if err != nil {
		SendError(w, fmt.Sprintf("Could not find media for id %v", id))
		return
	}

	// Open file
	file, err := os.Open(fmt.Sprintf("%s/%s", context.RootPath, media.Path))
	if err != nil {
		log.Printf("Could not load media %d! %v", id, err)
		SendError(w, "Could not load media")
		return
	}
	defer file.Close()

	// Send file
	w.Header().Set("Content-Type", media.ContentType)
	io.Copy(w, file)
}

func getThumbnail(w http.ResponseWriter, r *http.Request, context *utils.Context) {
	// Convert id to number
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, fmt.Sprintf("Invalid id %v", id))
		return
	}

	// Get media from ID
	media, err := utils.GetMediaFromID(id, context)
	if err != nil {
		SendError(w, fmt.Sprintf("Could not find media for id %v", id))
		return
	}

	// TODO: Check if thumbnail already exists before generating new one

	// TODO: Move thumbnail generation to separate function in utils

	// Open image
	file, err := os.Open(filepath.Join(context.RootPath, media.Path))
	if err != nil {
		log.Printf("Could not load media %d! %v", id, err)
		SendError(w, "Could not load media")
		return
	}
	defer file.Close()

	// Decode image
	image, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Could not decode image! %v", err)
		SendError(w, "Could not decode image")
		return
	}

	// Create thumbnail file
	thumbnailFile, err := os.Create(filepath.Join(context.CachePath, "thumbnails", media.Path))
	if err != nil {
		log.Printf("Could not create thumbnail image file! %v", err)
		SendError(w, "Could not create thumbnail image file")
		return
	}

	// Resize image and write to thumbnail file
	resizedImage := imaging.Resize(image, 128, 128, imaging.Lanczos)
	err = png.Encode(thumbnailFile, resizedImage)
	if err != nil {
		log.Printf("Could not encode thumbnail! %v", err)
		SendError(w, "Could not encode thumbnail")
		return
	}
	thumbnailFile.Close()

	thumbnailFile, err = os.Open(filepath.Join(filepath.Join(context.CachePath, "thumbnails", media.Path)))
	if err != nil {
		log.Printf("Could not load thumbnail! %v", err)
		SendError(w, "Could not load thumbnail")
		return
	}
	defer thumbnailFile.Close()

	// Send thumbnail file
	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, thumbnailFile)
}