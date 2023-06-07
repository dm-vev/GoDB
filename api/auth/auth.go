package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Auth struct {
	APIKey string
}

func New() *Auth {
	a := &Auth{}
	a.loadAPIKey()
	return a
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != a.APIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Auth) loadAPIKey() {
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		a.APIKey = apiKey
		return
	}

	if _, err := os.Stat("api.json"); err == nil {
		data, err := ioutil.ReadFile("api.json")
		if err != nil {
			panic("Failed to read api.json: " + err.Error())
		}

		var jsonData map[string]string
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			panic("Failed to parse api.json: " + err.Error())
		}

		apiKey, ok := jsonData["api_key"]
		if !ok {
			panic("API key not found in api.json")
		}

		a.APIKey = apiKey
		return
	}

	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		panic("Failed to generate API key: " + err.Error())
	}

	a.APIKey = hex.EncodeToString(buf)
	err = ioutil.WriteFile("api.json", []byte(`{"api_key": "`+a.APIKey+`"}`), 0644)
	if err != nil {
		panic("Failed to write API key to api.json: " + err.Error())
	}
}
