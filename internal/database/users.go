package database

import (
	"errors"
	"log"
	"time"
)

type User struct {
	Id        string
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       string    `json:"age"`
}

func (c Client) CreateUser(email, password, name string, age string) (User, error) {
	user := User{
		Email:     email,
		Password:  password,
		Name:      name,
		Age:       age,
		CreatedAt: time.Now().UTC(),
	}
	db, err := c.readDB()
	if err != nil {
		log.Fatal("error in reading db", err)
	}
	if _, ok := db.Users[email]; ok {
		return User{}, errors.New("existing user")
	}
	db.Users[email] = user
	err = c.updateDB(db)
	return user, err
}

func (c Client) UpdateUser(email, password, name string, age string) (User, error) {
	db, err := c.readDB()
	if err != nil {
		log.Fatal("db read failed")
	}
	if val, ok := db.Users[email]; ok {
		val.Email = email
		val.Password = password
		val.Name = name
		val.Age = age
		err = c.updateDB(db)
		return val, err
	}
	return User{}, errors.New("user doesn't exist")
}

func (c Client) GetUser(email string) (User, error) {
	db, err := c.readDB()
	if err != nil {
		return User{}, err
	}
	if usr, found := db.Users[email]; found {
		return usr, nil
	}
	return User{}, errors.New("User not found")

}

func (c Client) DeleteUser(email string) error {
	db, err := c.readDB()
	if err != nil {
		return err
	}
	delete(db.Users, email)
	return c.updateDB(db)
}
