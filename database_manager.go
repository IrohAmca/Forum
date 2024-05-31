package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var user_db *sql.DB

func createDatabase() {
	var err error
	user_db, err = sql.Open("sqlite3", "databases/database.db")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	creationQueries := []string{
		`CREATE TABLE IF NOT EXISTS Users (
				UserID INTEGER PRIMARY KEY AUTOINCREMENT,
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
				Title TEXT NOT NULL,
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

func Query_email(email string) (string, error) {
	var password string
	row := user_db.QueryRow("SELECT Email, Password FROM Users WHERE Email = ?", email)
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
	row := user_db.QueryRow("SELECT Email FROM Users WHERE Email = ?", email)
	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		}
	}
	if email == "" {
		return true
	}
	return false
}

func check_username(username string) bool {
	row := user_db.QueryRow("SELECT Nickname FROM Users WHERE Nickname = ?", username)
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	if username == "" {
		return false
	}
	return true
}

func Query_username(email string) (string, error) {
	var username string
	row := user_db.QueryRow("SELECT Nickname FROM Users WHERE Email = ?", email)
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
	row := user_db.QueryRow("SELECT Email, Password FROM Users WHERE UserID = ?", id)
	err := row.Scan(&email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user with id %d", id)
		}
		return "", "", fmt.Errorf("error scanning row: %v", err)
	}
	return email, password, nil
}

func SetMod(id int) error {
	statement, err := user_db.Prepare("UPDATE Users SET UserLevel = 1 WHERE UserID = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func SetAdmin(id int) error {
	statement, err := user_db.Prepare("UPDATE Users SET UserLevel = 2 WHERE UserID = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func insertUser(username, email, password string) error {
	hashedPassword := hashPassword(password)
	statement, err := user_db.Prepare("INSERT OR IGNORE INTO Users (UserLevel, Nickname, Email, Password) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(0, username, email, hashedPassword)
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}

	fmt.Println("Data inserted successfully.")
	return nil
}
func getUserID(email string) (int, error) {
	var id int
	row := user_db.QueryRow("SELECT UserID FROM Users WHERE Email = ?", email)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user with email %s", email)
		}
		return 0, fmt.Errorf("error scanning row: %v", err)
	}
	return id, nil
}
func insertPost(threadID, userID int, content string) error {
	statement, err := user_db.Prepare("INSERT INTO Posts (ThreadID, UserID, Content) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(threadID, userID, content)
	if err != nil {
		log.Println("Error executing statement,", err)
		return err
	}

	fmt.Println("Data inserted successfully.")
	return nil
}

func insertThread(userID, categoryID int, title string) (int,error) {
	statement, err := user_db.Prepare("INSERT INTO Threads (UserID, CategoryID, Title) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return 0,err
	}
	defer statement.Close()

	_, err = statement.Exec(userID, categoryID, title)
	if err != nil {
		log.Println("Error executing statement:",err)
		return 0,err
	}
	statement, err = user_db.Prepare("SELECT ThreadID FROM Threads WHERE Title = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return 0,err
	}
	defer statement.Close()

	var id int
	row := statement.QueryRow(title)
	err = row.Scan(&id)
	if err != nil {
		log.Println("Error scanning row:", err)
		return 0,err
	}
	return id,nil
}

func WriteAllData() {
	rows, err := user_db.Query("SELECT UserID, Nickname, Email, Password FROM Users")
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
