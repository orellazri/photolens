package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/orellazri/photolens/database"
	"github.com/orellazri/photolens/routes"
	"github.com/orellazri/photolens/utils"
)

func main() {
	// Environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}
	port := os.Getenv("PORT")
	dsn := os.Getenv("POSTGRES_URL")

	// Database
	db, err := database.SetupDatabase(dsn)
	if err != nil {
		log.Fatal("Could not setup databse")
	}
	err = database.MigrateDatabase(db)
	if err != nil {
		log.Fatal("Could not migrate database")
	}

	// Context
	context := utils.Context{
		DB: db,
	}

	// Routes
	router := mux.NewRouter()
	router.HandleFunc("/", routes.IndexHandler).Methods("GET")
	photosRouter := router.PathPrefix("/photos").Subrouter()
	routes.RegisterPhotosRouter(&context, photosRouter)

	// HTTP
	http.Handle("/", router)
	log.Printf("Photolens Server started on port %s...", port)
	http.ListenAndServe("localhost:"+port, router)
}
