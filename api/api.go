package api

import (
	"GemDB/api/auth"
	"GemDB/storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	Router *mux.Router
	Store  *storage.Storage
	Auth   *auth.Auth
}

func New(store *storage.Storage) *API {
	a := &API{
		Router: mux.NewRouter(),
		Store:  store,
		Auth:   auth.New(),
	}

	a.routes()
	return a
}

func (a *API) routes() {
	a.Router.Use(a.Auth.Middleware)
	a.Router.HandleFunc("/createTable", a.createTable).Methods("POST")
	a.Router.HandleFunc("/deleteTable", a.deleteTable).Methods("POST")
	a.Router.HandleFunc("/set", a.set).Methods("POST")
	a.Router.HandleFunc("/get", a.get).Methods("GET")
	a.Router.HandleFunc("/delete", a.delete).Methods("POST")
	a.Router.HandleFunc("/exportToFile", a.exportToFile).Methods("POST")
	a.Router.HandleFunc("/listTables", a.listTablesHandler).Methods("GET")
	a.Router.HandleFunc("/listKeys", a.listKeysHandler).Methods("GET")
}

type createTableRequest struct {
	Name string `json:"name"`
}

func (a *API) createTable(w http.ResponseWriter, r *http.Request) {
	var req createTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := a.Store.CreateTable(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type deleteTableRequest struct {
	Name string `json:"name"`
}

func (a *API) deleteTable(w http.ResponseWriter, r *http.Request) {
	var req deleteTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := a.Store.DeleteTable(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type setRequest struct {
	TableName string `json:"table_name"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (a *API) set(w http.ResponseWriter, r *http.Request) {
	var req setRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := a.Store.Set(req.TableName, req.Key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) get(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("table_name")
	key := r.URL.Query().Get("key")

	value, err := a.Store.Get(name, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"value": value})
}

type deleteRequest struct {
	TableName string `json:"table_name"`
	Key       string `json:"key"`
}

func (a *API) delete(w http.ResponseWriter, r *http.Request) {
	var req deleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := a.Store.Delete(req.TableName, req.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) exportToFile(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Filename string `json:"filename"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.Store.ExportToFile(requestData.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) listTablesHandler(w http.ResponseWriter, r *http.Request) {
	tables, err := a.Store.ListTables()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Tables []string `json:"tables"`
	}{
		Tables: tables,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (a *API) listKeysHandler(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table_name")
	if tableName == "" {
		http.Error(w, "Invalid request: table_name is required", http.StatusBadRequest)
		return
	}

	keys, err := a.Store.ListKeys(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Keys []string `json:"keys"`
	}{
		Keys: keys,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response) // TODO: неккоректная работа при большом количестве ключей
}
