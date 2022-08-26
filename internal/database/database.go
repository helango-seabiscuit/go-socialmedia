package database

import (
	"encoding/json"
	"errors"
	"os"
)

type Client struct {
	pathToFile string
}

func NewClient(path string) Client {
	return Client{
		pathToFile: path,
	}
}

func (c Client) createDb() error {
	data, err := json.Marshal(databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	})
	if err != nil {
		return err
	}
	return os.WriteFile(c.pathToFile, data, 0600)
}

func (c Client) EnsureDB() error {
	_, err := os.ReadFile(c.pathToFile)
	if errors.Is(err, os.ErrNotExist) {
		return c.createDb()
	}

	return err
}

type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

func (c Client) updateDB(db databaseSchema) error {
	data, err := json.Marshal(databaseSchema{
		Users: db.Users,
		Posts: db.Posts,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(c.pathToFile, data, 0600)
}

func (c Client) readDB() (databaseSchema, error) {
	data, err := os.ReadFile(c.pathToFile)
	db := &databaseSchema{}
	if err != nil {
		return databaseSchema{}, err
	}
	err = json.Unmarshal(data, db)
	if err != nil {
		return databaseSchema{}, err
	}
	return *db, nil
}
