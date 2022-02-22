package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/orellazri/photolens/models"
	"github.com/orellazri/photolens/utils"
)

func RegisterPhotosRouter(context *utils.Context, router *mux.Router) {
	router.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) { NewPhoto(w, r, context) }).Methods("POST")
	router.HandleFunc("/{id}", GetPhoto).Methods("GET")
}

func GetPhoto(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, "Invalid id")
		return
	}

	img, err := os.Open(fmt.Sprintf("%d.jpg", id))
	if err != nil {
		log.Printf("Could not load image %d", id)
		SendError(w, "Could not load image")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}

type NewPhotoRequest struct {
	Path string `json:"path"`
}

func NewPhoto(w http.ResponseWriter, r *http.Request, context *utils.Context) {
	decoder := json.NewDecoder(r.Body)
	var request NewPhotoRequest
	err := decoder.Decode(&request)
	if err != nil {
		SendError(w, "Invalid request")
		return
	}

	context.DB.Create(&models.Photo{
		Path: request.Path,
	})

	fmt.Fprintf(w, "Photo created")
}
