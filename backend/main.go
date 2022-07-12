package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/orellazri/photolens/core"
	"github.com/orellazri/photolens/database"
	"github.com/orellazri/photolens/routes"
)

func main() {
	// Database
	log.Println("Setting up database...")
	db, err := database.SetupDatabase("photolens.db")
	if err != nil {
		log.Fatalf("Could not setup databse! %v", err)
	}
	log.Println("Migrating database...")
	err = database.MigrateDatabase(db)
	if err != nil {
		log.Fatalf("Could not migrate database! %v", err)
	}

	// Context
	context := core.Context{
		DB:        db,
		RootPath:  "./media",
		CachePath: "./cache",
	}

	// TODO: Create root directory, cache directory and cache/thumbnails directory if they don't exist yet

	// Process media
	err = core.ProcessMedia(&context)
	if err != nil {
		log.Fatalf("Could not process media! %v", err)
	}

	// Routes
	router := mux.NewRouter()
	router.HandleFunc("/", routes.IndexHandler).Methods("GET")
	mediaRouter := router.PathPrefix("/media").Subrouter()
	mediaRouter.Use(routes.CorsMiddleware)
	routes.RegisterMediaRouter(&context, mediaRouter)

	// HTTP
	http.Handle("/", router)
	port := "5000"
	log.Printf("Photolens Server listening on port %s", port)
	http.ListenAndServe("localhost:"+port, router)
}
