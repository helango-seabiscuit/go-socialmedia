package database

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createUserTable string = `CREATE TABLE IF NOT EXISTS "users" (
      "id" INTEGER,
	  "name" TEXT NOT NULL,
	  "email" TEXT NOT NULL,
	  "age" TEXT NOT NULL,
	  "password" TEXT NOT NULL,
	  "created_date" DATETIME NOT NULL,
	  PRIMARY KEY(email)
	);`

	createPostTable string = `CREATE TABLE IF NOT EXISTS "posts" (
       "id" TEXT,
	   "userEmail" TEXT NOT NULL,
	   "text" TEXT NOT NULL,
	   "created_date" DATETIME NOT NULL,
	   PRIMARY KEY(id)
	);`
)

type DBClient struct {
	db *sql.DB
	sync.RWMutex
}

func NewSQLite3Repo(dbfile string) (*DBClient, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createUserTable); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createPostTable); err != nil {
		return nil, err
	}

	return &DBClient{
		db: db,
	}, nil

}

type Repository interface {
	RetrieveUser(email string) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	RetrievePosts(email string) ([]Post, error)
	DeletePost(id string) error
	CreatePost(post Post) (Post, error)
}

func (client *DBClient) RetrieveUser(email string) (User, error) {
	client.RLock()
	defer client.RUnlock()
	row := client.db.QueryRow("SELECT name,email,age,password,created_date  from users where email=?", email)
	usr := User{}
	err := row.Scan(&usr.Name, &usr.Email, &usr.Age, &usr.Password, &usr.CreatedAt)
	if err != nil {
		fmt.Println("error in retrieving ", err)
	}
	return usr, err
}

func (client *DBClient) CreateUser(user User) (User, error) {
	insStmt, err := client.db.Prepare("INSERT INTO users VALUES(NULL, ?,?,?,?,?)")
	if err != nil {
		fmt.Println("error while preparing statment ", err)
		return User{}, err
	}
	defer insStmt.Close()
	res, err := insStmt.Exec(user.Name, user.Email, user.Age, user.Password, time.Now())
	if err != nil {
		fmt.Println("error while storing data ", err)
		return User{}, err
	}
	fmt.Println(res.LastInsertId())
	return user, nil
}

func (client *DBClient) UpdateUser(user User) (User, error) {
	updStmt, err := client.db.Prepare("UPDATE users SET name=?, email=?,age=?,password=? where email=?")
	if err != nil {
		return User{}, err
	}

	defer updStmt.Close()
	res, err := updStmt.Exec(user.Name, user.Email, user.Age, user.Password, user.Email)
	if err != nil {
		return User{}, err
	}

	_, err = res.RowsAffected()
	return user, err
}
