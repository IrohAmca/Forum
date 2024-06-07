package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
				Likes INTEGER,
				Dislike INTEGER,
				CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (UserID) REFERENCES Users(UserID),
				FOREIGN KEY (ThreadID) REFERENCES Threads(ThreadID)
			);`,
		`CREATE TABLE IF NOT EXISTS LikesDislikes(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER NOT NULL,
			PostID INTEGER NOT NULL,
			IsLike BOOLEAN NOT NULL,
			CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (UserID) REFERENCES Users(UserID),
			FOREIGN KEY (PostID) REFERENCES Posts(PostID),
			UNIQUE(UserID, PostID)
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

	// Check if username or email already exists.
	var exists bool
	err := user_db.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE Nickname = ? OR Email = ?)", username, email).Scan(&exists)
	if err != nil {
		log.Println("Error checking for existing user:", err)
		return err
	}

	if exists {
		fmt.Println("Error: Username or email already exists.")
		return fmt.Errorf("username or email already exists")
	}

	statement, err := user_db.Prepare("INSERT INTO Users (UserLevel, Nickname, Email, Password) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(0, username, email, hashedPassword) // user lever for zoro
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User added successfully.")
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

func insertThread(userID, categoryID int, title string) (int, error) {
	statement, err := user_db.Prepare("INSERT INTO Threads (UserID, CategoryID, Title) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return 0, err
	}
	defer statement.Close()

	_, err = statement.Exec(userID, categoryID, title)
	if err != nil {
		log.Println("Error executing statement:", err)
		return 0, err
	}
	statement, err = user_db.Prepare("SELECT ThreadID FROM Threads WHERE Title = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return 0, err
	}
	defer statement.Close()

	var id int
	row := statement.QueryRow(title)
	err = row.Scan(&id)
	if err != nil {
		log.Println("Error scanning row:", err)
		return 0, err
	}
	return id, nil
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

type Post struct {
	PostID    int
	ThreadID  int
	Title     string
	UserID    int
	Username  string
	Content   string
	CreatedAt string
}

func getAllPosts() ([]Post, error) {
	posts := []Post{}
	rows, err := user_db.Query("SELECT PostID, ThreadID, UserID, Content, CreatedAt FROM Posts")
	if err != nil {
		log.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID, threadID, userID int
		var content string
		var createdAt time.Time
		err := rows.Scan(&postID, &threadID, &userID, &content, &createdAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		post := Post{
			PostID:    postID,
			ThreadID:  threadID,
			UserID:    userID,
			Content:   content,
			CreatedAt: createdAt.Format("2006-01-02 15:04:05"), // Format the time as desired
		}
		posts = append(posts, post)
	}

	for i, post := range posts {
		rows, err = user_db.Query("SELECT Nickname FROM Users WHERE UserID = ?", post.UserID)
		if err != nil {
			log.Println("Error querying data:", err)
		}
		defer rows.Close()
		rows.Next()
		var username string
		err := rows.Scan(&username)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		posts[i].Username = username
	}

	for i, post := range posts {
		rows, err = user_db.Query("SELECT Title FROM Threads WHERE ThreadID = ?", post.ThreadID)
		if err != nil {
			log.Println("Error querying data:", err)
		}
		defer rows.Close()
		rows.Next()
		var title string
		err := rows.Scan(&title)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		posts[i].Title = title
	}
	return posts, nil
}