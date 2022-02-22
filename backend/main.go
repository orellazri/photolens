package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/orellazri/photolens/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}
	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.HandleFunc("/", routes.IndexHandler).Methods("GET")
	photosRouter := router.PathPrefix("/photos").Subrouter()
	routes.RegisterPhotosRouter(photosRouter)

	http.Handle("/", router)
	log.Printf("Photolens Server started on port %s...", port)
	http.ListenAndServe("localhost:"+port, router)
}
