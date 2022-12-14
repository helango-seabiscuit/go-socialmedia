package main

import (
	"encoding/json"
	"log"
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

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
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
	jsonClient := database.NewClient("db.json")
	jsonClient.EnsureDB()

	dbClient, err := database.NewSQLite3Repo("media.db")
	if err != nil {
		log.Fatal(err)
	}
	config := apiConfig{
		dbClient:    jsonClient,
		dbSqlClient: dbClient,
	}

	route := gin.Default()
	route.GET("/ping", ping)
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
	// PORT env for specific port to listen to
	route.Run()

}
