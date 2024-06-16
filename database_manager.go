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

func QueryToken(username string) (string, error) {
	var token string
	row := user_db.QueryRow("SELECT Token FROM Users WHERE Nickname = ?", username)
	err := row.Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with username %s", username)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return token, nil
}
func QueryTokenID(id int) (string, string, error) {
	var token, username string
	row := user_db.QueryRow("SELECT Token, Nickname FROM Users WHERE UserID = ?", id)
	err := row.Scan(&token, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user with id %d", id)
		}
		return "", "", fmt.Errorf("error scanning row: %v", err)
	}
	return token, username, nil
}

func Query_ID(token string) (int, error) {
	var id int
	row := user_db.QueryRow("SELECT UserID FROM Users WHERE Token = ?", token)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user with token %s", token)
		}
		return 0, fmt.Errorf("error scanning row: %v", err)
	}
	return id, nil
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

func insertUser(username, email, password, token string) error {
	hashedPassword := hashPassword(password)
	statement, err := user_db.Prepare("INSERT OR IGNORE INTO Users (UserLevel, Nickname, Token, Email, Password) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(0, username, token, email, hashedPassword)
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
func category2ID(categories []string) ([]string, error) {
	var ids []string
	dict := map[string]string{
		`Gündem`:             "1",
		`Ev&Yaşam`:           "2",
		`Para&Ekonomi`:       "3",
		`Moda&Stil`:          "4",
		`İnternet&Teknoloji`: "5",
		`Eğitim&Kariyer`:     "6",
	}
	for _, category := range categories {
		id, ok := dict[category]
		if !ok {
			return nil, fmt.Errorf("no category with name %s", category)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
func insertPost(threadID, userID int, content string) error {
	statement, err := user_db.Prepare("INSERT INTO Posts (ThreadID, UserID, Content, Likes, Dislikes) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(threadID, userID, content, 0, 0)
	if err != nil {
		log.Println("Error executing statement,", err)
		return err
	}

	fmt.Println("Data inserted successfully.")
	return nil
}

func insertThread(userID int, title string, categories []string) (int, error) {
	ids, err := category2ID(categories)
	if err != nil {
		return 0, err
	}
	categoryIDs := ""
	for _, id := range ids {
		categoryIDs += id + ","
	}
	categoryIDs = categoryIDs[:len(categoryIDs)-1]
	statement, err := user_db.Prepare("INSERT INTO Threads (UserID, Title, CategoryIDs) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return 0, err
	}
	defer statement.Close()

	_, err = statement.Exec(userID, title, categoryIDs)
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

type Comment struct {
	CommentID int
	UserID    int
	PostID    int
	Content   string
	Username  string
	CreatedAt string
	UpdatedAt string
}

type Post struct {
	PostID         int
	ThreadID       int
	Title          string
	UserToken      string
	Username       string
	Content        string
	CreatedAt      string
	LikeCounter    int
	DislikeCounter int
	Comment        []Comment
}

func getCommentsByPostID(postID int) ([]Comment, error) {
	comments := []Comment{}

	rows, err := user_db.Query("SELECT CommentID, UserID, PostID, Content, CreatedAt, UpdatedAt FROM Comments WHERE PostID=?", postID)
	if err != nil {
		log.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		var createdAt, updatedAt time.Time
		err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.PostID, &comment.Content, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		comment.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		comment.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		comments = append(comments, comment)

	}

	for i, comment := range comments {
		rows, err = user_db.Query("SELECT Nickname FROM Users WHERE UserID=?", comment.UserID)
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
		comments[i].Username = username

	}

	return comments, nil

}
func getAllPosts() ([]Post, error) {
	posts := []Post{}
	rows, err := user_db.Query("SELECT PostID, ThreadID, UserID, Content, CreatedAt, Likes, Dislikes FROM Posts")
	if err != nil {
		log.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID, threadID, userID int
		var content string
		var createdAt time.Time
		var likes, dislikes int
		err := rows.Scan(&postID, &threadID, &userID, &content, &createdAt, &likes, &dislikes)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		userToken, username, err := QueryTokenID(userID)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		post := Post{
			PostID:         postID,
			ThreadID:       threadID,
			Content:        content,
			UserToken:      userToken,
			Username:       username,
			LikeCounter:    likes,
			DislikeCounter: dislikes,
			CreatedAt:      createdAt.Format("2006-01-02 15:04:05"),
		}
		posts = append(posts, post)
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

	for i, post := range posts {
		comments, err := getCommentsByPostID(post.PostID)
		if err != nil {
			log.Println("Error getting comments:", err)
			return nil, err
		}
		if comments == nil {
			comments = []Comment{}
		} else {
			posts[i].Comment = comments
		}
	}
	return posts, nil
}

func checkToken(token string) bool {
	var user string
	row := user_db.QueryRow("SELECT Token FROM Users WHERE Token = ?", token)
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	if user == "" {
		return false
	}
	return true
}

func deletePostFromDB(PostID int) error {
	statement, err := user_db.Prepare("DELETE FROM Posts WHERE PostID = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(PostID)
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}
	return nil
}
func insertComment(userID, postID int, content string) error {
	statement, err := user_db.Prepare("INSERT INTO Comments (UserID, PostID, Content) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userID, postID, content)
	if err != nil {
		log.Println("Error executing statement,", err)
		return err
	}

	fmt.Println("Data inserted successfully.")
	return nil

}

func deleteCommentFromDB(CommentID int) error {
	statement, err := user_db.Prepare("DELETE FROM Comments WHERE CommentID = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(CommentID)
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}
	return nil
}

type LikeDislikeActions struct {
	UserID int 
	PostID int  
	IsLike bool 
}

func HandleLikeDislike(action LikeDislikeActions) error {

	var currentID int
	var currentIsLike bool

	err := user_db.QueryRow(`SELECT ID, IsLike FROM LikesDislikes WHERE UserID = ? AND PostID = ?`, action.UserID, action.PostID).Scan(&currentID, &currentIsLike)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("query error: %v", err)
	}

	if err == sql.ErrNoRows {
		_, err = user_db.Exec(`INSERT INTO LikesDislikes (UserID, PostID, IsLike) VALUES (?,?,?) `, action.UserID, action.PostID, action.IsLike)
		if err != nil {
			return fmt.Errorf("insert error: %v", err)
		}

		if action.IsLike {
			_, err = user_db.Exec(`UPDATE Posts SET Likes = Likes + 1 WHERE PostID = ?`, action.PostID)
		} else {
			_, err = user_db.Exec(`UPDATE Posts SET Dislikes = Dislikes + 1 WHERE PostID = ?`, action.PostID)
		}
		if err != nil {
			return fmt.Errorf("update count error: %v", err)
		}
	} else {
		if currentIsLike == action.IsLike {
			_, err = user_db.Exec(`DELETE FROM LikesDislikes WHERE ID = ?`, currentID)
			if err != nil {
				return fmt.Errorf("delete error: %v", err)
			}

			if action.IsLike {
				_, err = user_db.Exec(`UPDATE Posts SET Likes = Likes - 1 WHERE PostID = ?`, action.PostID)
			} else {
				_, err = user_db.Exec(`UPDATE Posts SET Dislikes = Dislikes - 1 WHERE PostID = ?`, action.PostID)
			}

			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		} else {
			_, err = user_db.Exec(`UPDATE LikesDislikes SET IsLike = ? WHERE ID = ?`, action.IsLike, currentID)
			if err != nil {
				return fmt.Errorf("update error: %v", err)
			}

			if action.IsLike {
				_, err = user_db.Exec(`UPDATE Posts SET Likes = Likes + 1, Dislikes = Dislikes - 1 WHERE PostID = ?`, action.PostID)
			} else {
				_, err = user_db.Exec(`UPDATE Posts SET Likes = Likes - 1, Dislikes = Dislikes + 1 WHERE PostID = ?`, action.PostID)
			}

			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		}
	}
	return nil
}
func getFilteredPosts(categories []string, title string) ([]Post, error) {
	var posts []Post

	if len(categories) == 0 && title == "" {
		return getAllPosts()
	}

	categoryPosts := make(map[int]Post)
	if len(categories) > 0 {
		ids, err := category2ID(categories)
		if err != nil {
			return nil, err
		}
		categoryIDs := ""
		for _, id := range ids {
			categoryIDs += id + ","
		}
		categoryIDs = categoryIDs[:len(categoryIDs)-1]

		rows, err := user_db.Query("SELECT ThreadID FROM Threads WHERE CategoryIDs = ?", categoryIDs)
		if err != nil {
			log.Println("Error querying data:", err)
			return nil, err
		}
		defer rows.Close()
		var threadID int
		for rows.Next() {
			err := rows.Scan(&threadID)
			if err != nil {
				log.Println("Error scanning row:", err)
				return nil, err
			}
			tempPosts, err := getPostByThreadID(threadID)
			if err != nil {
				log.Println("Error getting posts:", err)
				return nil, err
			}
			for _, post := range tempPosts {
				categoryPosts[post.PostID] = post
			}
		}
	}
	titlePosts := make(map[int]Post)
	if title != "" {
		rows, err := user_db.Query("SELECT ThreadID FROM Threads WHERE Title = ?", title)
		if err != nil {
			log.Println("Error querying data:", err)
			return nil, err
		}
		defer rows.Close()
		var threadID int
		for rows.Next() {
			err := rows.Scan(&threadID)
			if err != nil {
				log.Println("Error scanning row:", err)
				return nil, err
			}
			tempPosts, err := getPostByThreadID(threadID)
			if err != nil {
				log.Println("Error getting posts:", err)
				return nil, err
			}
			for _, post := range tempPosts {
				titlePosts[post.PostID] = post
			}
		}
	}
	if len(categories) > 0 && title != "" {
		for id, post := range categoryPosts {
			if _, exists := titlePosts[id]; exists {
				posts = append(posts, post)
			}
		}
	} else if len(categories) > 0 {
		for _, post := range categoryPosts {
			posts = append(posts, post)
		}
	} else if title != "" {
		for _, post := range titlePosts {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func getPostByThreadID(threadID int) ([]Post, error) {
	posts := []Post{}
	rows, err := user_db.Query("SELECT PostID, UserID, Content, CreatedAt, Likes, Dislikes FROM Posts WHERE ThreadID = ?", threadID)
	if err != nil {
		log.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID, userID int
		var content string
		var createdAt time.Time
		var likes, dislikes int
		err := rows.Scan(&postID, &userID, &content, &createdAt, &likes, &dislikes)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		userToken, username, err := QueryTokenID(userID)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		post := Post{
			PostID:         postID,
			ThreadID:       threadID,
			Content:        content,
			UserToken:      userToken,
			Username:       username,
			LikeCounter:    likes,
			DislikeCounter: dislikes,
			CreatedAt:      createdAt.Format("2006-01-02 15:04:05"),
		}
		posts = append(posts, post)
	}
	fmt.Println("Posts:", posts)
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

	for i, post := range posts {
		comments, err := getCommentsByPostID(post.PostID)
		if err != nil {
			log.Println("Error getting comments:", err)
			return nil, err
		}
		if comments == nil {
			comments = []Comment{}
		} else {
			posts[i].Comment = comments
		}
	}
	return posts, nil
}

func sortByDateAsc(posts []Post) []Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].CreatedAt > posts[j].CreatedAt {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func sortByDateDesc(posts []Post) []Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].CreatedAt < posts[j].CreatedAt {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func sortByLikeAsc(posts []Post) []Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].LikeCounter > posts[j].LikeCounter {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func sortByLikeDesc(posts []Post) []Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].LikeCounter < posts[j].LikeCounter {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}
