package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

func createTable(db *sql.DB) error {
	query := `
	        CREATE TABLE users (
	            id INT AUTO_INCREMENT,
	            username TEXT NOT NULL,
	            password TEXT NOT NULL,
	            created_at DATETIME,
	            PRIMARY KEY (id)
	        );`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func createUser(db *sql.DB) error {
	username := "johndoe"
	password := "secret"
	createdAt := time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	fmt.Println(id)
	if err != nil {
		return err
	}
	return nil
}

func deleteUser(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1)
	if err != nil {
		return err
	}
	return nil
}

func allUsers(db *sql.DB) error {
	type user struct {
		id        int
		username  string
		password  string
		createdAt time.Time
	}

	rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user

		err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
		if err != nil {
			return err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	fmt.Printf("%#v", users)
	return nil
}

func connection(config Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", config.User, config.Password, config.ServerName, config.DB)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
    	log.Fatalf("Error loading .env file")
	}
	config := Config{
		ServerName: os.Getenv("ServerName"),
		User:       os.Getenv("User"),
		Password:   os.Getenv("Password"),
		DB:         os.Getenv("DB"),
	}
	fmt.Println(os.Getenv("ServerName"))
	conn, err := connection(config)
	if err != nil {
		log.Fatal(err)
	}
	// createTable(conn)
	// createUser(conn)
	// deleteUser(conn)
	allUsers(conn)
}
