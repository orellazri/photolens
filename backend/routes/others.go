package routes

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	msg := MessageResponse{
		Message: "Welcome to Photolens API",
		Error:   false,
	}
	jsonResposne, err := json.Marshal(msg)
	if err != nil {
		return
	}
	w.Write(jsonResposne)
}
