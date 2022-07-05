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
	db, err := database.SetupDatabase()
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
	port := "5000"
	log.Printf("Photolens Server started on port %s...", port)
	http.ListenAndServe("localhost:"+port, router)
}
