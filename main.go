package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/helango-seabiscuit/go-socialmedia/internal/database"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, database.User{
		Email: "test@example.com",
		Name:  "test",
		Age:   "16",
	})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 200, errors.New("testing error"))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response, err := json.Marshal(payload)
	if err != nil {
		response = []byte("{ error: error marshalling}")
	}
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	errorBdy := errorBody{
		Error: err.Error(),
	}
	respondWithJSON(w, 200, errorBdy)
}

type errorBody struct {
	Error string `json:"error"`
}

func main() {
	dbClient := database.NewClient("db.json")
	dbClient.EnsureDB()
	config := apiConfig{
		dbClient: dbClient,
	}
	m := http.NewServeMux()
	m.HandleFunc("/", testHandler)
	m.HandleFunc("/err", testErrHandler)
	m.HandleFunc("/users", config.endpointUsersHandler)
	m.HandleFunc("/users/", config.endpointUsersHandler)
	m.HandleFunc("/posts", config.endpointPostsHandler)
	m.HandleFunc("/posts/", config.endpointPostsHandler)
	const addr = "localhost:8001"
	srv := http.Server{
		Handler:      m,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	fmt.Println("server started pn ", addr)
	err := srv.ListenAndServe()
	log.Fatal(err)

}
