package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/orellazri/photolens/utils"
)

func RegisterPhotosRouter(context *utils.Context, router *mux.Router) {
	router.HandleFunc("/{id}", GetPhoto)
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
