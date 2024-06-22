package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createPost(c *gin.Context) {
	var post struct {
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		Categories []string `json:"categories" binding:"required"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	threadID, err := insertThread(userID, post.Title, post.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	insertPost(threadID, userID, post.Content)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post created successfully"})
}

func shortPost(short_type string, posts []Post) []Post {
	if short_type != "" {
		if short_type == "date-asc" {
			fmt.Println("date-asc")
			posts = sortByDateAsc(posts)
		} else if short_type == "date-desc" {
			fmt.Println("date-desc")
			posts = sortByDateDesc(posts)
		} else if short_type == "likes-asc" {
			fmt.Println("likes-asc")
			posts = sortByLikeAsc(posts)
		} else if short_type == "likes-desc" {
			fmt.Println("likes-desc")
			posts = sortByLikeDesc(posts)
		}
	}
	return posts
}
func getPosts(c *gin.Context) {
	var categories struct {
		Categories []string `json:"categories"`
		Title      string   `json:"title"`
		Short_type string   `json:"short_type"`
	}
	if err := c.ShouldBind(&categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	post ,err :=getFilteredPosts(categories.Categories, categories.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	post = shortPost(categories.Short_type, post)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Posts retrieved successfully", "posts": post})
}

func deletePost(c *gin.Context) {
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
	err = deletePostFromDB(postID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post deleted successfully"})
}
func createComment(c *gin.Context) {
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
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	err = insertComment(userID, postID, comment.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment created successfully"})
}

func deleteComment(c *gin.Context) {
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
	err = deleteCommentFromDB(commentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment deleted successfully"})
}

func likeDislikePost(c *gin.Context) {
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
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, err := Query_ID(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	action := LikeDislikeActions{
		UserID: userID,
		PostID: postID,
		IsLike: like.IsLike,
	}
	err = HandleLikeDislike(action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Like/Dislike action successful"})
}