package service

import (
	"fmt"
	"forum/models"
	"forum/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var post struct {
		Title        string   `json:"title" binding:"required"`
		Content      string   `json:"content" binding:"required"`
		Categories   []string `json:"categories" binding:"required"`
		Image        []byte   `json:"image"`
		Ext          string   `json:"ext"`
		Device_Token string   `json:"device_token" binding:"required"`
		Cookie       string   `json:"cookie" binding:"required"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Please fill in all the fields"})
		return
	}
	if repository.CheckDeviceToken(post.Device_Token[:len(post.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	if len(post.Categories) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Please select at least one category"})
		return
	}
	if post.Cookie == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Please log in to the website first!!!"})
		return
	}
	token, err := repository.GetTokenByCookie(post.Cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Please log in to the website first!!!"})
		return
	}
	userID, err := repository.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	threadID, err := repository.InsertThread(userID, post.Title, post.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	repository.InsertPost(threadID, userID, post.Content, post.Image, post.Ext)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post created successfully"})
}

func ShortPost(short_type string, posts []models.Post) []models.Post {
	if short_type != "" {
		if short_type == "date-asc" {
			fmt.Println("date-asc")
			posts = repository.SortByDateAsc(posts)
		} else if short_type == "date-desc" {
			fmt.Println("date-desc")
			posts = repository.SortByDateDesc(posts)
		} else if short_type == "likes-asc" {
			fmt.Println("likes-asc")
			posts = repository.SortByLikeAsc(posts)
		} else if short_type == "likes-desc" {
			fmt.Println("likes-desc")
			posts = repository.SortByLikeDesc(posts)
		}
	}
	return posts
}
func GetPosts(c *gin.Context) {
	var categories struct {
		Categories   []string `json:"categories"`
		Title        string   `json:"title"`
		Short_type   string   `json:"short_type"`
		Device_Token string   `json:"device_token"`
	}
	if err := c.ShouldBind(&categories); err != nil {
		print("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if repository.CheckDeviceToken(categories.Device_Token[:len(categories.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}

	post, err := repository.GetFilteredPosts(categories.Categories, categories.Title)
	if err != nil {
		print("Error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	post = ShortPost(categories.Short_type, post)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Posts retrieved successfully", "posts": post})
}

func DeletePost(c *gin.Context) {
	var post struct {
		PostID       string `json:"PostID" binding:"required"`
		Device_Token string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	if repository.CheckDeviceToken(post.Device_Token[:len(post.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}

	postID, err := strconv.Atoi(post.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	err = repository.DeletePostFromDB(postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post deleted successfully"})
}
func CreateComment(c *gin.Context) {
	var comment struct {
		PostID       string `json:"postId"`
		Content      string `json:"content"`
		Cookie       string `json:"cookie"`
		Device_Token string `json:"device_token"`
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
	token, err := repository.GetTokenByCookie(comment.Cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := repository.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	err = repository.InsertComment(userID, postID, comment.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment created successfully"})
}

func DeleteComment(c *gin.Context) {
	var comment struct {
		CommentID    string `json:"CommentID" binding:"required"`
		Device_Token string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	if repository.CheckDeviceToken(comment.Device_Token[:len(comment.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	commentID, err := strconv.Atoi(comment.CommentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid CommentID"})
		return
	}
	err = repository.DeleteCommentFromDB(commentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment deleted successfully"})
}

func LikeDislikePost(c *gin.Context) {
	var like struct {
		PostID       string `json:"PostID"`
		Cookie       string `json:"cookie"`
		Device_Token string `json:"device_token"`
		IsLike       bool   `json:"isLike"`
	}
	if err := c.ShouldBindJSON(&like); err != nil {
		fmt.Println("Reading Error" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error: " + err.Error()})
		return
	}
	if repository.CheckDeviceToken(like.Device_Token[:len(like.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}

	postID, err := strconv.Atoi(like.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	token, err := repository.GetTokenByCookie(like.Cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Cookie Error: " + err.Error()})
		return
	}
	userID, err := repository.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	action := repository.LikeDislikePostActions{
		UserID: userID,
		PostID: postID,
		IsLike: like.IsLike,
	}
	err = repository.HandleLikeDislike(action)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Like/Dislike action successful"})
}

func LikeDislikeComment(c *gin.Context) {
	var like struct {
		CommentID    string `json:"CommentID" binding:"required"`
		IsLike       bool   `json:"isLike"`
		Cookie       string `json:"cookie"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error: " + err.Error()})
		return
	}
	if repository.CheckDeviceToken(like.Device_Token[:len(like.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	commentID, err := strconv.Atoi(like.CommentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	token, err := repository.GetTokenByCookie(like.Cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := repository.Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	action := repository.LikeDislikeCommentActions{
		UserID:    userID,
		CommentID: commentID,
		IsLike:    like.IsLike,
	}
	err = repository.HandleLikeDislikeComment(action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Like/Dislike action successful"})
}

func GetImage(c *gin.Context) {
	filename := c.Param("filename")
	postID, err := strconv.Atoi(filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}

	image, ext, err := repository.GetImage(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	var contentType string
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpeg", ".jpg":
		contentType = "image/jpeg"
	case ".svg":
		contentType = "image/svg+xml"
	case ".gif", "webp":
		contentType = "image/gif"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Unsupported image format"})
		return
	}

	c.Data(http.StatusOK, contentType, image)
}
