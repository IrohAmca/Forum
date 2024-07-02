package services

import (
	"fmt"
	"forum/db_manager"
	"forum/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var post struct {
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		Categories []string `json:"categories" binding:"required"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Please fill in all the fields"})
		return
	}
	if len(post.Categories) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Please select at least one category"})
		return
	}
	cookie, err := c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Please log in to the website first!!!"})
		return
	}
	token, err := db_manager.GetTokenByCookie(cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Please log in to the website first!!!"})
		return
	}
	userID, err := db_manager.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	threadID, err := db_manager.InsertThread(userID, post.Title, post.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	db_manager.InsertPost(threadID, userID, post.Content)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post created successfully"})
}

func ShortPost(short_type string, posts []models.Post) []models.Post {
	if short_type != "" {
		if short_type == "date-asc" {
			fmt.Println("date-asc")
			posts = db_manager.SortByDateAsc(posts)
		} else if short_type == "date-desc" {
			fmt.Println("date-desc")
			posts = db_manager.SortByDateDesc(posts)
		} else if short_type == "likes-asc" {
			fmt.Println("likes-asc")
			posts = db_manager.SortByLikeAsc(posts)
		} else if short_type == "likes-desc" {
			fmt.Println("likes-desc")
			posts = db_manager.SortByLikeDesc(posts)
		}
	}
	return posts
}
func GetPosts(c *gin.Context) {
	var categories struct {
		Categories []string `json:"categories"`
		Title      string   `json:"title"`
		Short_type string   `json:"short_type"`
	}
	if err := c.ShouldBind(&categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	post, err := db_manager.GetFilteredPosts(categories.Categories, categories.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	post = ShortPost(categories.Short_type, post)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Posts retrieved successfully", "posts": post})
}

func DeletePost(c *gin.Context) {
	var post struct {
		PostID string `json:"PostID" binding:"required"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	postID, err := strconv.Atoi(post.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	err = db_manager.DeletePostFromDB(postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post deleted successfully"})
}
func CreateComment(c *gin.Context) {
	var comment struct {
		PostID  string `json:"postId" binding:"required"`
		Content string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	postID, err := strconv.Atoi(comment.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	cookie, err := c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	token, err := db_manager.GetTokenByCookie(cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := db_manager.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	err = db_manager.InsertComment(userID, postID, comment.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment created successfully"})
}

func DeleteComment(c *gin.Context) {
	var comment struct {
		CommentID string `json:"CommentID" binding:"required"`
	}
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	commentID, err := strconv.Atoi(comment.CommentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid CommentID"})
		return
	}
	err = db_manager.DeleteCommentFromDB(commentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment deleted successfully"})
}

func LikeDislikePost(c *gin.Context) {
	var like struct {
		PostID string `json:"PostID" binding:"required"`
		IsLike bool   `json:"isLike"`
	}
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error: " + err.Error()})
		return
	}
	postID, err := strconv.Atoi(like.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	cookie, err := c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	token, err := db_manager.GetTokenByCookie(cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := db_manager.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	action := db_manager.LikeDislikePostActions{
		UserID: userID,
		PostID: postID,
		IsLike: like.IsLike,
	}
	err = db_manager.HandleLikeDislike(action)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Like/Dislike action successful"})
}

func LikeDislikeComment(c *gin.Context) {
	var like struct {
		CommentID string `json:"CommentID" binding:"required"`
		IsLike    bool   `json:"isLike"`
	}
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error: " + err.Error()})
		return
	}
	commentID, err := strconv.Atoi(like.CommentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	cookie, err := c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	token, err := db_manager.GetTokenByCookie(cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := db_manager.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	action := db_manager.LikeDislikeCommentActions{
		UserID:    userID,
		CommentID: commentID,
		IsLike:    like.IsLike,
	}
	err = db_manager.HandleLikeDislikeComment(action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Like/Dislike action successful"})
}
