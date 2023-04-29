package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var database = make(map[string]interface{})
var file = "database.json"
var mutex = &sync.Mutex{}
var logger *log.Logger

func execute(action string, data map[string]interface{}) interface{} {
	key := data["key"].(string)
	if action == "get" {
		logger.Print("get request for key:", key)
		return database[key]
	} else if action == "insert" {
		value := data["value"]
		database[key] = value
		logger.Print("insert request for key:", key)
		return "OK"
	} else if action == "create" {
		database[key] = make(map[string]interface{})
		logger.Print("create request for key:", key)
		return "OK"
	} else if action == "delete" {
		delete(database, key)
		logger.Print("delete request for key:", key)
		return "OK"
	}
	return nil
}

func save() {
	mutex.Lock()
	defer mutex.Unlock()
	bytes, err := json.Marshal(database)
	if err != nil {
		logger.Println("Error marshaling database:", err)
		return
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		logger.Println("Error writing database file:", err)
		return
	}
	logger.Println("Database saved to file")
}

func main() {
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		logger.Fatal("Error reading config file:", err)
	}

	config := make(map[string]string)
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		logger.Fatal("Error parsing config file:", err)
	}

	access_token := config["access_token"]

	// Create logger that writes to file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()
	logger = log.New(logFile, "", log.LstdFlags)

	bytes, err := ioutil.ReadFile(file)
	if err == nil {
		err = json.Unmarshal(bytes, &database)
		if err != nil {
			logger.Println("Error unmarshaling database:", err)
		}
	}

	http.HandleFunc("/call/", func(w http.ResponseWriter, r *http.Request) {
		event := r.URL.Path[len("/call/"):]
		var data map[string]interface{}
		err := json.Unmarshal([]byte(event), &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.Header.Get("access-token") != access_token {
			logger.Println("Unauthorized user: ", r.RemoteAddr, "! Added to blacklist.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		res := execute(data["action"].(string), data["data"].(map[string]interface{}))
		if res != nil && data["action"].(string) != "get" {
			save()
		}
		json.NewEncoder(w).Encode(res)
	})

	logger.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
