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
	response := MessageResponse{
		Message: "Welcome to Photolens API",
		Error:   false,
	}
	jsonResposne, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(jsonResposne)
}

func SendError(w http.ResponseWriter, msg string) {
	response := MessageResponse{
		Message: msg,
		Error:   true,
	}
	jsonResposne, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(jsonResposne)
}
