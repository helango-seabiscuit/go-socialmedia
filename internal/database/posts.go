package database

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
}

func (c Client) CreatePost(userEmail, text string) (Post, error) {
	db, err := c.readDB()
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := db.Users[userEmail]; ok {
		id := uuid.New().String()
		post := Post{
			UserEmail: userEmail,
			Text:      text,
			CreatedAt: time.Now().UTC(),
			ID:        id,
		}
		db.Posts[id] = post
		err = c.updateDB(db)
		return post, err
	}
	return Post{}, err
}

func (c Client) DeletePost(id string) error {
	db, err := c.readDB()
	if err != nil {
		log.Fatal(err)
	}
	delete(db.Posts, id)
	err = c.updateDB(db)
	return err
}

func (c Client) GetPosts(userEmail string) ([]Post, error) {
	db, err := c.readDB()
	if err != nil {
		log.Fatal(err)
	}
	posts := []Post{}
	for _, v := range db.Posts {
		if v.UserEmail == userEmail {
			posts = append(posts, v)
		}
	}
	return posts, nil
}
