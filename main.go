package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helango-seabiscuit/go-socialmedia/internal/database"
)

func testHandler(c *gin.Context) {
	c.JSON(http.StatusOK, database.User{
		Email: "test@example.com",
		Name:  "test",
		Age:   "16",
	})
}

func testErrHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"err": "testing error"})

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

	route := gin.Default()
	route.GET("/", testHandler)
	route.GET("/err", testErrHandler)
	route.POST("/login", config.HandleLogin)
	route.POST("/users", config.HandleCreateUser)
	route.GET("/users/:email", config.HandleGetUser)
	route.DELETE("/users/:email", config.HandleDeleteUser)
	route.GET("/posts/:email", config.HandleRetrievePosts)
	route.POST("/posts", config.HandleCreatePost)
	route.DELETE("/posts/:id", config.handleDeletePost)

	// m.HandleFunc("/posts", config.endpointPostsHandler)
	// m.HandleFunc("/posts/", config.endpointPostsHandler)

	route.Run()

}
