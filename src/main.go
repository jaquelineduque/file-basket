package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	port := GetPort()
	log.Println("listening on", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/files", uploadFile).
		Methods("POST")
	router.HandleFunc("/v1/files/{id:[A-z0-9-]+}", retrieveFilePaths).
		Methods("GET")
	router.PathPrefix("/v1/files/").Handler(retrieveFile()).Methods("GET")
	http.ListenAndServe(":"+port, router)
}
