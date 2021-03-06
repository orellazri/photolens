package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orellazri/photolens/core"
	"github.com/orellazri/photolens/models"
	"golang.org/x/exp/slices"
)

func RegisterMediaRouter(context *core.Context, router *mux.Router) {
	router.HandleFunc("/meta", func(w http.ResponseWriter, r *http.Request) { getMetadata(w, r, context) }).Methods("GET")
	router.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) { getProcessMedia(w, r, context) }).Methods("GET")
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { getMedia(w, r, context) }).Methods("GET")
	// router.HandleFunc("/thumbnail/all", func(w http.ResponseWriter, r *http.Request) { getAllThumbnails(w, r, context) }).Methods("GET")
	router.HandleFunc("/thumbnail/{id}", func(w http.ResponseWriter, r *http.Request) { getThumbnail(w, r, context) }).Methods("GET")
}

func getMetadata(w http.ResponseWriter, r *http.Request, context *core.Context) {
	type thumbnailResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		LastModified time.Time `json:"last_modified"`
	}

	type response struct {
		Data []thumbnailResponse `json:"data"`
	}

	// Get parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	sortByParam := r.URL.Query().Get("sortby")
	sortBy := "created_at"
	if slices.Contains([]string{"created_at", "last_modified"}, sortByParam) {
		sortBy = sortByParam
	}

	sortDirParam := r.URL.Query().Get("sortdir")
	sortDir := "desc"
	if slices.Contains([]string{"desc", "asc"}, sortDirParam) {
		sortDir = sortDirParam
	}

	// Get all media files from database
	var results []models.Media
	err := context.DB.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sortBy, sortDir)).
		Select("id", "created_at", "last_modified").
		Find(&results).
		Error
	if err != nil {
		log.Printf("Could not get media from database! %v", err)
		SendError(w, "Could not get metadata")
		return
	}

	var metadatas []thumbnailResponse
	for _, result := range results {
		metadatas = append(metadatas, thumbnailResponse{
			ID:           result.ID,
			CreatedAt:    result.CreatedAt,
			LastModified: result.LastModified,
		})
	}

	SendJsonResponse(w, response{
		Data: metadatas,
	})
}

func getMedia(w http.ResponseWriter, r *http.Request, context *core.Context) {
	// Convert id string to uuid
	id := mux.Vars(r)["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Could not find parse UUID for id %v! %v", id, err)
		SendError(w, fmt.Sprintf("Could not find media for id %v", id))
		return
	}

	// Get media from ID
	media, err := core.GetMediaFromID(uuid, context)
	if err != nil {
		log.Printf("Could not find media for id %v! %v", id, err)
		SendError(w, fmt.Sprintf("Could not find media for id %v", id))
		return
	}

	// Open file
	file, err := os.Open(media.Path)
	if err != nil {
		log.Printf("Could not load media %v! %v", id, err)
		SendError(w, "Could not load media")
		return
	}
	defer file.Close()

	// Send file
	w.Header().Set("Content-Type", media.ContentType)
	io.Copy(w, file)
}

func getThumbnail(w http.ResponseWriter, r *http.Request, context *core.Context) {
	// Convert id string to uuid
	id := mux.Vars(r)["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Could not find parse UUID for id %v! %v", id, err)
		SendError(w, fmt.Sprintf("Could not find media for id %v", id))
		return
	}

	// Get media from ID
	media, err := core.GetMediaFromID(uuid, context)
	if err != nil {
		log.Printf("Could not find media for id %v! %v", id, err)
		SendError(w, fmt.Sprintf("Could not find media for id %v", id))
		return
	}

	// Generate thumbnail (or get existing)
	thumbnailPath, err := core.GetThumbnail(context, media)
	if err != nil {
		log.Printf("Could not generate thumbnail! %v", err)
		SendError(w, "Could not generate thumbnail")
		return
	}

	// Open file
	file, err := os.Open(thumbnailPath)
	if err != nil {
		log.Printf("Could not load thumbnail! %v", err)
		SendError(w, "Could not load thumbnail")
		return
	}
	defer file.Close()

	// Send file
	w.Header().Set("Content-Type", media.ContentType)
	io.Copy(w, file)
}

func getAllThumbnails(w http.ResponseWriter, r *http.Request, context *core.Context) {
	type thumbnailResponse struct {
		ID        uuid.UUID `json:"id"`
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
		media, err := core.GetMediaFromID(result.ID, context)
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

func getProcessMedia(w http.ResponseWriter, r *http.Request, context *core.Context) {
	err := core.ProcessMedia(context)
	if err != nil {
		log.Printf("Could not process media! %v", err)
		SendError(w, "Could not get process media")
		return
	}

	SendJsonResponse(w, MessageResponse{
		Message: "OK",
		Error:   false,
	})
}
