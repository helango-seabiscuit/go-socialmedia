package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type parameters struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      string `json:"age"`
	Password string `json:"password"`
}

func (a apiConfig) endpointUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.HandleGetUser(w, r)
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		user := parameters{}
		err := decoder.Decode(&user)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
		}
		usr, err := a.dbClient.CreateUser(user.Email, user.Password, user.Name, user.Age)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err)
		}
		fmt.Println("user created ", usr)
		respondWithJSON(w, http.StatusCreated, usr)
	case http.MethodPut:
	case http.MethodDelete:
		a.handleDeleteUser(w, r)
	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}

func (a apiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/users/")

	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}
	usr, err := a.dbClient.UpdateUser(email, user.Password, user.Name, user.Age)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}
	fmt.Println("user updated ", usr)
	respondWithJSON(w, http.StatusOK, usr)

}

func (a apiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/users/")
	user, err := a.dbClient.GetUser(email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (a apiConfig) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Path
	email = strings.TrimPrefix(email, "/users/")
	err := a.dbClient.DeleteUser(email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
