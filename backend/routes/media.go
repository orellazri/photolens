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

func RegisterMediaRouter(context *utils.Context, router *mux.Router) {
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) { GetMedia(w, r, context) }).Methods("GET")
}

func GetMedia(w http.ResponseWriter, r *http.Request, context *utils.Context) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, "Invalid id")
		return
	}

	// TODO: Look for path from id in DB

	img, err := os.Open(fmt.Sprintf("%s/%d.jpg", context.RootPath, id))
	if err != nil {
		log.Printf("Could not load media %d", id)
		SendError(w, "Could not load media")
		return
	}

	// TODO: Set correct content type
	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}
