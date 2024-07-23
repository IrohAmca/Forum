package repository

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"strconv"
	"strings"
	"time"
)

func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	comments := []models.Comment{}

	rows, err := User_db.Query("SELECT CommentID, UserID, PostID, Content, CreatedAt, UpdatedAt, Likes, Dislikes FROM Comments WHERE PostID=?", postID)
	if err != nil {
		log.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		var createdAt, updatedAt time.Time
		err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.PostID, &comment.Content, &createdAt, &updatedAt, &comment.LikeCounter, &comment.DislikeCounter)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		comment.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		comment.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		comments = append(comments, comment)

	}

	for i, comment := range comments {
		rows, err = User_db.Query("SELECT Nickname FROM Users WHERE UserID=?", comment.UserID)
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

func GetCategoriesByID(id string) ([]string, error) {
	var categories []string
	idList := strings.Split(id, ",")
	dict := map[string]string{
		"1": "Gündem",
		"2": "Ev&Yaşam",
		"3": "Para&Ekonomi",
		"4": "Moda&Stil",
		"5": "İnternet&Teknoloji",
		"6": "Eğitim&Kariyer",
	}
	for _, id := range idList {
		category, ok := dict[id]
		if !ok {
			return nil, fmt.Errorf("no category with id %s", id)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetAllPosts() ([]models.Post, error) {
	posts := []models.Post{}

	rows, err := User_db.Query("SELECT PostID, ThreadID, UserID, Image, Content, CreatedAt, Likes, Dislikes FROM Posts")
	if err != nil {
		log.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var image []byte
		var createdAt time.Time
		err := rows.Scan(&post.PostID, &post.ThreadID, &post.UserID, &image, &post.Content, &createdAt, &post.LikeCounter, &post.DislikeCounter)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		post.Image = ""
		if image != nil {
			post.Image = "http://localhost:8081/images/" + strconv.Itoa(post.PostID)
		}

		userToken, username, err := QueryTokenID(post.UserID)
		if err != nil {
			log.Println("Error querying user token:", err)
			return nil, err
		}

		post.UserToken = userToken
		post.Username = username
		post.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		posts = append(posts, post)
	}

	if err := addThreadDetails(posts); err != nil {
		return nil, err
	}

	if err := addCommentsToPosts(posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func addThreadDetails(posts []models.Post) error {
	for i, post := range posts {
		threadRows, err := User_db.Query("SELECT Title, CategoryIDs FROM Threads WHERE ThreadID = ?", post.ThreadID)
		if err != nil {
			log.Println("Error querying thread data:", err)
			return err
		}
		defer threadRows.Close()

		if threadRows.Next() {
			var title, categoryIDs string
			if err := threadRows.Scan(&title, &categoryIDs); err != nil {
				log.Println("Error scanning thread row:", err)
				return err
			}

			categories, err := GetCategoriesByID(categoryIDs)
			if err != nil {
				log.Println("Error querying categories:", err)
				return err
			}

			posts[i].Title = title
			posts[i].Categories = categories
		}
	}
	return nil
}

func addCommentsToPosts(posts []models.Post) error {
	for i, post := range posts {
		comments, err := GetCommentsByPostID(post.PostID)
		if err != nil {
			log.Println("Error getting comments:", err)
			return err
		}
		posts[i].Comment = comments
	}
	return nil
}

func DeletePostFromDB(PostID int) error {
	var threadID int
	row := User_db.QueryRow("SELECT ThreadID FROM Posts WHERE PostID = ?", PostID)
	err := row.Scan(&threadID)
	if err != nil {
		log.Println("Error querying data:", err)
		return err
	}

	statement, err := User_db.Prepare("DELETE FROM Posts WHERE PostID = ?")
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

	statement, err = User_db.Prepare("DELETE FROM Comments WHERE PostID = ?")
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

	statement, err = User_db.Prepare("DELETE FROM Threads WHERE ThreadID = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(threadID)
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}
	return nil
}

func DeleteCommentFromDB(CommentID int) error {
	statement, err := User_db.Prepare("DELETE FROM Comments WHERE CommentID = ?")
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

type LikeDislikePostActions struct {
	UserID int
	PostID int
	IsLike bool
}

func HandleLikeDislike(action LikeDislikePostActions) error {

	var currentID int
	var currentIsLike bool

	err := User_db.QueryRow(`SELECT ID, IsLike FROM PostLikesDislikes WHERE UserID = ? AND PostID = ?`, action.UserID, action.PostID).Scan(&currentID, &currentIsLike)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("query error: %v", err)
	}

	if err == sql.ErrNoRows {
		_, err = User_db.Exec(`INSERT INTO PostLikesDislikes (UserID, PostID, IsLike) VALUES (?,?,?) `, action.UserID, action.PostID, action.IsLike)
		if err != nil {
			return fmt.Errorf("insert error: %v", err)
		}

		if action.IsLike {
			_, err = User_db.Exec(`UPDATE Posts SET Likes = Likes + 1 WHERE PostID = ?`, action.PostID)
		} else {
			_, err = User_db.Exec(`UPDATE Posts SET Dislikes = Dislikes + 1 WHERE PostID = ?`, action.PostID)
		}
		if err != nil {
			return fmt.Errorf("update count error: %v", err)
		}
	} else {
		if currentIsLike == action.IsLike {
			_, err = User_db.Exec(`DELETE FROM PostLikesDislikes WHERE ID = ?`, currentID)
			if err != nil {
				return fmt.Errorf("delete error: %v", err)
			}

			if action.IsLike {
				_, err = User_db.Exec(`UPDATE Posts SET Likes = Likes - 1 WHERE PostID = ?`, action.PostID)
			} else {
				_, err = User_db.Exec(`UPDATE Posts SET Dislikes = Dislikes - 1 WHERE PostID = ?`, action.PostID)
			}

			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		} else {
			_, err = User_db.Exec(`UPDATE PostLikesDislikes SET IsLike = ? WHERE ID = ?`, action.IsLike, currentID)
			if err != nil {
				return fmt.Errorf("update error: %v", err)
			}

			if action.IsLike {
				_, err = User_db.Exec(`UPDATE Posts SET Likes = Likes + 1, Dislikes = Dislikes - 1 WHERE PostID = ?`, action.PostID)
			} else {
				_, err = User_db.Exec(`UPDATE Posts SET Likes = Likes - 1, Dislikes = Dislikes + 1 WHERE PostID = ?`, action.PostID)
			}

			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		}
	}
	return nil
}

type LikeDislikeCommentActions struct {
	UserID    int
	CommentID int
	IsLike    bool
}

func HandleLikeDislikeComment(action LikeDislikeCommentActions) error {
	var currentID int
	var currentIsLike bool

	err := User_db.QueryRow(`SELECT ID, IsLike FROM CommentLikesDislikes WHERE UserID = ? AND CommentID = ?`, action.UserID, action.CommentID).Scan(&currentID, &currentIsLike)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("query error: %v", err)
	}

	if err == sql.ErrNoRows {
		_, err = User_db.Exec(`INSERT INTO CommentLikesDislikes (UserID, CommentID, IsLike) VALUES (?,?,?) `, action.UserID, action.CommentID, action.IsLike)
		if err != nil {
			return fmt.Errorf("insert error: %v", err)
		}

		if action.IsLike {
			_, err = User_db.Exec(`UPDATE Comments SET Likes = Likes + 1 WHERE CommentID = ?`, action.CommentID)
		} else {
			_, err = User_db.Exec(`UPDATE Comments SET Dislikes = Dislikes + 1 WHERE CommentID = ?`, action.CommentID)
		}
		if err != nil {
			return fmt.Errorf("update count error: %v", err)
		}
	} else {
		if currentIsLike == action.IsLike {
			_, err = User_db.Exec(`DELETE FROM CommentLikesDislikes WHERE ID = ?`, currentID)
			if err != nil {
				return fmt.Errorf("delete error: %v", err)
			}

			if action.IsLike {
				_, err = User_db.Exec(`UPDATE Comments SET Likes = Likes - 1 WHERE CommentID = ?`, action.CommentID)
			} else {
				_, err = User_db.Exec(`UPDATE Comments SET Dislikes = Dislikes - 1 WHERE CommentID = ?`, action.CommentID)
			}

			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		} else {
			_, err = User_db.Exec(`UPDATE CommentLikesDislikes SET IsLike = ? WHERE ID = ?`, action.IsLike, currentID)
			if err != nil {
				return fmt.Errorf("update error: %v", err)
			}
			if action.IsLike {
				_, err = User_db.Exec(`UPDATE Comments SET Likes = Likes + 1, Dislikes = Dislikes - 1 WHERE CommentID = ?`, action.CommentID)
			} else {
				_, err = User_db.Exec(`UPDATE Comments SET Likes = Likes - 1, Dislikes = Dislikes + 1 WHERE CommentID = ?`, action.CommentID)
			}

			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		}
	}
	return nil
}

func GetFilteredPosts(categories []string, title string) ([]models.Post, error) {
	var posts []models.Post

	if len(categories) == 0 && title == "" {
		return GetAllPosts()
	}

	categoryPosts := make(map[int]models.Post)
	if len(categories) > 0 {
		ids, err := Category2ID(categories)
		if err != nil {
			return nil, err
		}
		categoryIDs := ""
		for _, id := range ids {
			categoryIDs += id + ","
		}
		categoryIDs = categoryIDs[:len(categoryIDs)-1]

		rows, err := User_db.Query("SELECT ThreadID FROM Threads WHERE CategoryIDs = ?", categoryIDs)
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
			tempPosts, err := GetPostByThreadID(threadID)
			if err != nil {
				log.Println("Error getting posts:", err)
				return nil, err
			}
			for _, post := range tempPosts {
				categoryPosts[post.PostID] = post
				post := categoryPosts[post.PostID]
				post.Categories = categories
				categoryPosts[post.PostID] = post
			}
		}
	}
	titlePosts := make(map[int]models.Post)
	if title != "" {
		rows, err := User_db.Query("SELECT ThreadID FROM Threads WHERE Title = ?", title)
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
			tempPosts, err := GetPostByThreadID(threadID)
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

func GetPostByThreadID(threadID int) ([]models.Post, error) {
	posts := []models.Post{}
	rows, err := User_db.Query("SELECT PostID, UserID, Content, CreatedAt, Likes, Dislikes FROM Posts WHERE ThreadID = ?", threadID)
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
		post := models.Post{
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
		rows, err = User_db.Query("SELECT Title FROM Threads WHERE ThreadID = ?", post.ThreadID)
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
		comments, err := GetCommentsByPostID(post.PostID)
		if err != nil {
			log.Println("Error getting comments:", err)
			return nil, err
		}
		if comments != nil {
			posts[i].Comment = comments
		}
	}
	return posts, nil
}

func SortByDateAsc(posts []models.Post) []models.Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].CreatedAt > posts[j].CreatedAt {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func SortByDateDesc(posts []models.Post) []models.Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].CreatedAt < posts[j].CreatedAt {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func SortByLikeAsc(posts []models.Post) []models.Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].LikeCounter > posts[j].LikeCounter {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func SortByLikeDesc(posts []models.Post) []models.Post {
	for i := 0; i < len(posts); i++ {
		for j := i + 1; j < len(posts); j++ {
			if posts[i].LikeCounter < posts[j].LikeCounter {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func InsertComment(userID, postID int, content string) error {
	statement, err := User_db.Prepare("INSERT INTO Comments (UserID, PostID, Content, Likes, Dislikes) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userID, postID, content, 0, 0)
	if err != nil {
		log.Println("Error executing statement,", err)
		return err
	}

	fmt.Println("Data inserted successfully.")
	return nil

}
func GetPostByID(postID int) ([]models.Post, error) {
	var posts []models.Post
	row := User_db.QueryRow("SELECT ThreadID, UserID, Content, CreatedAt, Likes, Dislikes FROM Posts WHERE PostID = ?", postID)
	var post models.Post
	var threadID, userID, likes, dislikes int
	var content, createdAt string
	err := row.Scan(&threadID, &userID, &content, &createdAt, &likes, &dislikes)
	if err != nil {
		return []models.Post{}, err
	}
	userToken, username, err := QueryTokenID(userID)
	if err != nil {
		return []models.Post{}, err
	}
	post = models.Post{
		PostID:         postID,
		ThreadID:       threadID,
		Content:        content,
		UserToken:      userToken,
		Username:       username,
		LikeCounter:    likes,
		DislikeCounter: dislikes,
		CreatedAt:      createdAt,
	}
	posts = append(posts, post)

	return posts, nil
}

func GetImage(postID int) ([]byte, string, error) {
	row := User_db.QueryRow("SELECT Image ,Ext FROM Posts WHERE PostID = ?", postID)
	var image []byte
	var ext string
	err := row.Scan(&image, &ext)
	if err != nil {
		return nil, "", err
	}
	return image, ext, nil
}


func InsertPost(threadID, userID int, content string, image []byte, ext string) error {
	statement, err := User_db.Prepare("INSERT INTO Posts (ThreadID, UserID, Content, Image, Ext, Likes, Dislikes) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(threadID, userID, content, image, ext, 0, 0)
	if err != nil {
		log.Println("Error executing statement,", err)
		return err
	}
	return nil
}