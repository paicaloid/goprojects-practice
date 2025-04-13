package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name string `json:"name"`
}

var userTable = make(map[int]User)

var cacheMutex sync.RWMutex

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/users", handleUser)
	mux.HandleFunc("/users/{id}", getUser)

	fmt.Println("Server listening to :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheMutex.RLock()
	user, ok := userTable[id]
	cacheMutex.RUnlock()
	if !ok {
		http.Error(w, "user not found.", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/users" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		cacheMutex.Lock()
		userTable[len(userTable)+1] = user
		cacheMutex.Unlock()

		w.WriteHeader(http.StatusCreated)

	} else {
		http.Error(w, "Method is not  supported.", http.StatusMethodNotAllowed)
		return
	}

}
