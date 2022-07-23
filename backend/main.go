package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// Create cache directory
	err = os.Mkdir(context.CachePath, os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "exists") {
		log.Fatalf("Could not create cache directory! %v", err)
	}

	// Create thumbnails cache directory
	err = os.Mkdir(filepath.Join(context.CachePath, "thumbnails"), os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "exists") {
		log.Fatalf("Could not create thumbnails cache directory! %v", err)
	}

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
