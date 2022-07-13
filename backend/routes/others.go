package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendJsonResponse(w http.ResponseWriter, response interface{}) {
	jsonResposne, err := json.Marshal(response)
	if err != nil {
		log.Printf("Could not send json response! %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResposne)
}

type MessageResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	SendJsonResponse(w, MessageResponse{
		Message: "Welcome to Photolens API",
		Error:   false,
	})
}

func SendError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	SendJsonResponse(w, MessageResponse{
		Message: msg,
		Error:   true,
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		next.ServeHTTP(w, r)
	})
}
