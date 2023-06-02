package main

import (
	"fmt"
	"log"
	"net/http"

	"GemDB/api"
	"GemDB/storage"
)

func main() {
	// Инициализация хранилища
	store, err := storage.New("data.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Инициализация API
	apiServer := api.New(store) //6784

	fmt.Println("Starting HTTP server...")
	err = http.ListenAndServe(":8080", apiServer.Router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
