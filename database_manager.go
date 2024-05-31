package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var userID int
var user_db *sql.DB

func createDatabase() {
	var err error
	user_db, err = sql.Open("sqlite3", "databases/user.db")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	creationQueries := []string{
		`CREATE TABLE IF NOT EXISTS Users (
			UserID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT,
			Lastname TEXT,
			Nickname TEXT NOT NULL UNIQUE,
			Email TEXT NOT NULL UNIQUE,
			UserBirthdate DATE,
			Password TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS Posts (
			PostID INTEGER PRIMARY KEY AUTOINCREMENT,
			ThreadID INTEGER,
			UserID INTEGER,
			Content TEXT NOT NULL,
			CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (UserID) REFERENCES Users(UserID),
			FOREIGN KEY (ThreadID) REFERENCES Threads(ThreadID)
		);`,
		`CREATE TABLE IF NOT EXISTS Likes(
			LikeID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER,
			PostID INTEGER,
			FOREIGN KEY (UserID) REFERENCES Users(UserID),
			FOREIGN KEY (PostID) REFERENCES Posts(PostID)
		);`,
		`CREATE TABLE IF NOT EXISTS Threads(
			ThreadID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER,
			CategoryID INTEGER,
			CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (UserID) REFERENCES Users(UserID),
			FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID)
		);`,
		`CREATE TABLE IF NOT EXISTS Categories(
			CategoryID INTEGER PRIMARY KEY AUTOINCREMENT,
			CategoryName TEXT NOT NULL UNIQUE,
			CategoryDescription TEXT 
		);`,
		`CREATE TABLE IF NOT EXISTS Comments (
			CommentID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID  INTEGER,
			PostID INTEGER,
			Content TEXT NOT NULL,
			CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (UserID) REFERENCES Users(UserID),
			FOREIGN KEY (PostID) REFERENCES Posts(PostID)
		);`}

	for _, query := range creationQueries {
		_, err := user_db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func updateUserID() {
	rows, err := user_db.Query("SELECT MAX(id) FROM people")
	if err != nil {
		log.Fatal("Error querying data:", err)
	}
	defer rows.Close()

	var maxID sql.NullInt64
	if rows.Next() {
		err := rows.Scan(&maxID)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
	}

	if maxID.Valid {
		userID = int(maxID.Int64) + 1
	} else {
		userID = 0
	}
}

func Query_email(email string) (string, error) {
	var password string
	row := user_db.QueryRow("SELECT email, password FROM people WHERE email = ?", email)
	err := row.Scan(&email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with email %s", email)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return password, nil
}
func check_email(email string) bool {
	row := user_db.QueryRow("SELECT email FROM people WHERE email = ?", email)
	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		}
		return false
	}
	return true
}
func check_username(username string) bool {
	row := user_db.QueryRow("SELECT username FROM people WHERE username = ?", username)
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}
	return true
}

func Query_username(email string) (string, error) {
	var username string
	row := user_db.QueryRow("SELECT username FROM people WHERE email = ?", email)
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with email %s", email)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return username, nil
}

func Query(id int) (string, string, error) {
	var email, password string
	row := user_db.QueryRow("SELECT email, password FROM people WHERE id = ?", id)
	err := row.Scan(&email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user with id %d", id)
		}
		return "", "", fmt.Errorf("error scanning row: %v", err)
	}
	return email, password, nil
}

func insertData(id int, username, email, password string) error {
	hashedPassword := hashPassword(password)
	statement, err := user_db.Prepare("INSERT OR IGNORE INTO people (id, username, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id, username, email, hashedPassword)
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}

	fmt.Println("Data inserted successfully.")
	return nil
}

func WriteAllData() {
	rows, err := user_db.Query("SELECT id, username, email, password FROM people")
	if err != nil {
		log.Println("Error querying data:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email, password string
		err := rows.Scan(&id, &username, &email, &password)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("ID: %d, Username: %s, Email: %s, Password: %s\n", id, username, email, password)
	}
}
