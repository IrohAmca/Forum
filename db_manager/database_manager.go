package db_manager

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var User_db *sql.DB

func CreateDatabase() {
	var err error
	User_db, err = sql.Open("sqlite3", "databases/database.db")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	creationQueries := []string{
		`CREATE TABLE IF NOT EXISTS Users (
				UserID INTEGER PRIMARY KEY AUTOINCREMENT,
				Token TEXT,
				UserLevel INTEGER NOT NULL,
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
				Likes INTEGER,
				Dislikes INTEGER,
				CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (UserID) REFERENCES Users(UserID),
				FOREIGN KEY (ThreadID) REFERENCES Threads(ThreadID)
			);`,
		`CREATE TABLE IF NOT EXISTS PostLikesDislikes(
				ID INTEGER PRIMARY KEY AUTOINCREMENT,
				UserID INTEGER NOT NULL,
				PostID INTEGER,
				IsLike BOOLEAN NOT NULL,
				CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (UserID) REFERENCES Users(UserID),
				FOREIGN KEY (PostID) REFERENCES Posts(PostID),
				UNIQUE(UserID, PostID)
				);`,
		`CREATE TABLE IF NOT EXISTS CommentLikesDislikes(
				ID INTEGER PRIMARY KEY AUTOINCREMENT,
				UserID INTEGER NOT NULL,
				PostID INTEGER,
				CommentID INTEGER,
				IsLike BOOLEAN NOT NULL,
				CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (UserID) REFERENCES Users(UserID),
				FOREIGN KEY (CommentID) REFERENCES Comments(CommentID),
				UNIQUE(UserID, CommentID)
				);`,
		`CREATE TABLE IF NOT EXISTS Threads(
				ThreadID INTEGER PRIMARY KEY AUTOINCREMENT,
				Title TEXT NOT NULL,
				UserID INTEGER,
				CategoryIDs TEXT NOT NULL,
				CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (UserID) REFERENCES Users(UserID),
				FOREIGN KEY (CategoryIDs) REFERENCES Categories(CategoryIDs)
			);`,
		`CREATE TABLE IF NOT EXISTS Categories(
				CategoryIDs TEXT PRIMARY KEY,
				CategoryNames TEXT NOT NULL UNIQUE,
				CategoryDescription TEXT 
			);`,
		`CREATE TABLE IF NOT EXIST Session(
			Token TEXT,
			Cookie TEXT
			);`,
		`CREATE TABLE IF NOT EXISTS Comments (
				CommentID INTEGER PRIMARY KEY AUTOINCREMENT,
				UserID  INTEGER,
				PostID INTEGER,
				Content TEXT NOT NULL,
				Likes INTEGER,
				Dislikes INTEGER,
				CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (UserID) REFERENCES Users(UserID),
				FOREIGN KEY (PostID) REFERENCES Posts(PostID)
			);`}

	for _, query := range creationQueries {
		_, err := User_db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

}
func WriteAllData() {
	rows, err := User_db.Query("SELECT UserID, Nickname, Email, Password FROM Users")
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
