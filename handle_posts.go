package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/helango-seabiscuit/go-socialmedia/internal/database"
)

type parameter struct {
	UserEmail string `json:"userEmail"`
	Text      string `json:"text"`
}

func (a apiConfig) HandleCreatePost(c *gin.Context) {
	post := parameter{}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	pst := database.Post{
		ID:        uuid.NewString(),
		UserEmail: post.UserEmail,
		Text:      post.Text,
		CreatedAt: time.Now(),
	}
	pst, err := a.dbSqlClient.CreatePost(pst)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println("post created ", pst)
	c.JSON(http.StatusCreated, pst)
}

func (a apiConfig) HandleRetrievePosts(c *gin.Context) {
	email := c.Param("email")
	posts, err := a.dbSqlClient.RetrievePosts(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, posts)
}

func (a apiConfig) handleDeletePost(c *gin.Context) {
	pid := c.Param("id")
	err := a.dbSqlClient.DeletePost(pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}
