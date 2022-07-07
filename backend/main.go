package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/orellazri/photolens/database"
	"github.com/orellazri/photolens/routes"
	"github.com/orellazri/photolens/utils"
)

func main() {
	// Database
	log.Println("Setting up database...")
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("Could not setup databse! %s", err)
	}
	log.Println("Migrating database...")
	err = database.MigrateDatabase(db)
	if err != nil {
		log.Fatalf("Could not migrate database! %s", err)
	}

	// Context
	context := utils.Context{
		DB:       db,
		RootPath: "./media",
	}

	// Process media
	err = utils.ProcessMedia(&context)
	if err != nil {
		log.Fatalf("Could not process media! %s", err)
	}

	// Routes
	router := mux.NewRouter()
	router.HandleFunc("/", routes.IndexHandler).Methods("GET")
	mediaRouter := router.PathPrefix("/media").Subrouter()
	routes.RegisterMediaRouter(&context, mediaRouter)

	// HTTP
	http.Handle("/", router)
	port := "5000"
	log.Printf("Photolens Server listening on port %s", port)
	http.ListenAndServe("localhost:"+port, router)
}
