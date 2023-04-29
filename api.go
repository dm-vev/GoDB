package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var database = make(map[string]interface{})
var file = "C:\\Users\\igovn\\GolandProjects\\GoDB\\database.json"
var mutex = &sync.Mutex{}

func execute(action string, data map[string]interface{}) interface{} {
	key := data["key"].(string)
	if action == "get" {
		return database[key]
	} else if action == "insert" {
		value := data["value"]
		database[key] = value
		return "OK"
	} else if action == "create" {
		database[key] = make(map[string]interface{})
		return "OK"
	} else if action == "delete" {
		delete(database, key)
		return "OK"
	}
	return nil
}

func save() {
	mutex.Lock()
	defer mutex.Unlock()

	bytes, err := json.Marshal(database)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(file, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	bytes, err := ioutil.ReadFile(file)
	if err == nil {
		err = json.Unmarshal(bytes, &database)
		if err != nil {
			panic(err)
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
		res := execute(data["action"].(string), data["data"].(map[string]interface{}))
		if res != nil && data["action"].(string) != "get" {
			save()
		}
		json.NewEncoder(w).Encode(res)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
