package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/orellazri/photolens/core"
	"github.com/orellazri/photolens/models"
)

func RegisterMediaRouter(context *core.Context, router *mux.Router) {
	router.HandleFunc("/meta", func(w http.ResponseWriter, r *http.Request) { getMetadata(w, r, context) }).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { getMedia(w, r, context) }).Methods("GET")
	router.HandleFunc("/thumbnail/all", func(w http.ResponseWriter, r *http.Request) { getAllThumbnails(w, r, context) }).Methods("GET")
	router.HandleFunc("/thumbnail/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { getThumbnail(w, r, context) }).Methods("GET")
}

func getMetadata(w http.ResponseWriter, r *http.Request, context *core.Context) {
	type mediaMetadataResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	type response struct {
		Data []mediaMetadataResponse `json:"data"`
	}

	// Get all media files from database
	var results []models.Media
	err := context.DB.Select("id", "created_at").Find(&results).Error
	if err != nil {
		log.Printf("Could not get media from database! %v", err)
		SendError(w, "Could not get metadata")
		return
	}

	var metadatas []mediaMetadataResponse
	for _, result := range results {
		metadatas = append(metadatas, mediaMetadataResponse{
			ID:        result.ID,
			CreatedAt: result.CreatedAt,
		})
	}

	SendJsonResponse(w, response{
		Data: metadatas,
	})
}

func getMedia(w http.ResponseWriter, r *http.Request, context *core.Context) {
	// Convert id to number
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, fmt.Sprintf("Invalid id %v", id))
		return
	}

	// Get media from ID
	media, err := core.GetMediaFromID(id, context)
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

func getThumbnail(w http.ResponseWriter, r *http.Request, context *core.Context) {
	type thumbnailResponse struct {
		ID        uint      `json:"id"`
		Thumbnail string    `json:"thumbnail"`
		CreatedAt time.Time `json:"created_at"`
	}

	type response struct {
		Data thumbnailResponse `json:"data"`
	}

	// Convert id to number
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, fmt.Sprintf("Invalid id %v", id))
		return
	}

	// Get media from ID
	media, err := core.GetMediaFromID(id, context)
	if err != nil {
		SendError(w, fmt.Sprintf("Could not find media for id %v", id))
		return
	}

	// Generate/get thumbnail
	thumbnailString, err := core.GetThumbnail(context, media)
	if err != nil {
		log.Printf("Could not generate thumbnail! %v", err)
		SendError(w, "Could not generate thumbnail")
		return
	}

	// Send thumbnail base64 encoded string
	SendJsonResponse(w, response{
		Data: thumbnailResponse{
			ID:        media.ID,
			Thumbnail: thumbnailString,
			CreatedAt: media.CreatedAt,
		},
	})
}

func getAllThumbnails(w http.ResponseWriter, r *http.Request, context *core.Context) {
	type thumbnailResponse struct {
		ID        uint      `json:"id"`
		Thumbnail string    `json:"thumbnail"`
		CreatedAt time.Time `json:"created_at"`
	}

	type response struct {
		Data []thumbnailResponse `json:"data"`
	}

	// Get all media files from database
	var results []models.Media
	err := context.DB.Select("id", "created_at").Find(&results).Error
	if err != nil {
		log.Printf("Could not get media from database! %v", err)
		SendError(w, "Could not get all thumbnails")
		return
	}

	var thumbnails []thumbnailResponse
	// TODO: Look into doing this in parallel
	for _, result := range results {
		media, err := core.GetMediaFromID(int(result.ID), context)
		if err != nil {
			log.Printf("Could not get media from ID! %v", err)
			SendError(w, "Could not get thumbnails")
			return
		}
		thumbnailString, err := core.GetThumbnail(context, media)
		if err != nil {
			log.Printf("Could not get thumbnail for image! %v", err)
			SendError(w, "Could not get thumbnails")
			return
		}
		thumbnails = append(thumbnails, thumbnailResponse{
			ID:        result.ID,
			Thumbnail: thumbnailString,
			CreatedAt: result.CreatedAt,
		})
	}

	SendJsonResponse(w, response{
		Data: thumbnails,
	})
}
