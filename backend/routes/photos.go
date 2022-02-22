package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RegisterPhotosRouter(router *mux.Router) {
	router.HandleFunc("/{id}", GetPhoto)
}

func GetPhoto(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, "Invalid id")
		return
	}
	fmt.Fprintf(w, "Get photo %d", id)
}
