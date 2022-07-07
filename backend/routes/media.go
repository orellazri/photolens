package routes

import (
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

func RegisterMediaRouter(context *utils.Context, router *mux.Router) {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { GetMedia(w, r, context) }).Methods("GET")
}

func GetMedia(w http.ResponseWriter, r *http.Request, context *utils.Context) {
	idStr := mux.Vars(r)["id"]

	// Check that the id is a number
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, "Invalid id")
		return
	}

	// Query path from id in database
	var media models.Media
	err = context.DB.First(&media, id).Error
	if err != nil {
		log.Printf("Could not find image in database! %v", err)
		SendError(w, "Invalid id")
		return
	}

	img, err := os.Open(fmt.Sprintf("%s/%s", context.RootPath, media.Path))
	if err != nil {
		log.Printf("Could not load media %d! %v", id, err)
		SendError(w, "Could not load media")
		return
	}

	w.Header().Set("Content-Type", media.ContentType)
	io.Copy(w, img)
}
