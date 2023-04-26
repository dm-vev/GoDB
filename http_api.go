package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	db := make(map[string]*Collection)

	http.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Create a new collection
			var name string
			err := json.NewDecoder(r.Body).Decode(&name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if _, ok := db[name]; ok {
				http.Error(w, "collection already exists", http.StatusBadRequest)
				return
			}
			db[name] = &Collection{documents: make(map[string]*Document)}
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/collections/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// Insert a new document into a collection
			collectionName := r.URL.Path[len("/collections/"):]
			col, ok := db[collectionName]
			if !ok {
				http.Error(w, "collection not found", http.StatusNotFound)
				return
			}
			var doc Document
			err := json.NewDecoder(r.Body).Decode(&doc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = col.InsertDocument(&doc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)

		case http.MethodGet:
			// Search documents in a collection by value
			collectionName := r.URL.Path[len("/collections/"):]
			col, ok := db[collectionName]
			if !ok {
				http.Error(w, "collection not found", http.StatusNotFound)
				return
			}
			value := r.URL.Query().Get("value")
			if value == "" {
				http.Error(w, "missing value parameter", http.StatusBadRequest)
				return
			}
			results := col.SearchByValue(value)
			json.NewEncoder(w).Encode(results)

		default:
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
