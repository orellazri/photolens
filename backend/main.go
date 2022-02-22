package main

import (
	"fmt"
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

	http.Handle("/", router)
	fmt.Println("Photolens Server started on port " + port + "...")
	http.ListenAndServe("localhost:"+port, router)
}
