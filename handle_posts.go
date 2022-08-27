package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	pst, err := a.dbClient.CreatePost(post.UserEmail, post.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println("post created ", pst)
	c.JSON(http.StatusCreated, pst)
}

func (a apiConfig) HandleRetrievePosts(c *gin.Context) {
	email := c.Param("email")
	posts, err := a.dbClient.GetPosts(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, posts)
}

func (a apiConfig) handleDeletePost(c *gin.Context) {
	pid := c.Param("id")
	err := a.dbClient.DeletePost(pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}
