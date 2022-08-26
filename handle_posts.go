package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type parameter struct {
	UserEmail string `json:"userEmail"`
	Text      string `json:"text"`
}

func (a apiConfig) endpointPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.HandleRetrievePosts(w, r)
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		post := parameter{}
		err := decoder.Decode(&post)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
		}
		pst, err := a.dbClient.CreatePost(post.UserEmail, post.Text)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
		}
		fmt.Println("post created ", pst)
		respondWithJSON(w, http.StatusCreated, pst)
	case http.MethodPut:
	case http.MethodDelete:
		a.handleDeletePost(w, r)
	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}

func (a apiConfig) HandleRetrievePosts(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/posts/")
	posts, err := a.dbClient.GetPosts(email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (a apiConfig) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	pid := r.URL.Path
	pid = strings.TrimPrefix(pid, "/posts/")
	err := a.dbClient.DeletePost(pid)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
