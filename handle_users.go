package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helango-seabiscuit/go-socialmedia/internal/database"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const (
	MIN_ENTROPY_BITS = 68
)

type parameters struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      string `json:"age"`
	Password string `json:"password"`
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a apiConfig) HandleLogin(c *gin.Context) {
	var login LoginForm
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := a.dbClient.GetUser(login.Email)
	if err != nil || user.Password != login.Password {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	c.JSON(http.StatusOK, "login successful")
}

func (a apiConfig) HandleCreateUser(c *gin.Context) {
	user := parameters{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := passwordvalidator.Validate(user.Password, MIN_ENTROPY_BITS); err != nil {

		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	u := database.User{
		Name:     user.Name,
		Email:    user.Email,
		Age:      user.Age,
		Password: user.Password,
	}
	usr, err := a.dbSqlClient.CreateUser(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("user created ", usr)
	c.JSON(http.StatusOK, usr)
}

func (a apiConfig) HandleGetUser(c *gin.Context) {
	email := c.Param("email")
	user, err := a.dbSqlClient.RetrieveUser(email)
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
