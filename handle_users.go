package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type parameters struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      string `json:"age"`
	Password string `json:"password"`
}

func (a apiConfig) HandleCreateUser(c *gin.Context) {
	user := parameters{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usr, err := a.dbClient.CreateUser(user.Email, user.Password, user.Name, user.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("user created ", usr)
	c.JSON(http.StatusOK, usr)
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

func (a apiConfig) HandleGetUser(c *gin.Context) {
	email := c.Param("email")
	user, err := a.dbClient.GetUser(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a apiConfig) HandleDeleteUser(c *gin.Context) {
	email := c.Param("email")
	err := a.dbClient.DeleteUser(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}
