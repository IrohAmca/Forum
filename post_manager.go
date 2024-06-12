package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createPost(c *gin.Context) {
	var post struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
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
	threadID, err := insertThread(userID, 0, post.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	insertPost(threadID, userID, post.Content)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post created successfully"})
}

func getPosts(c *gin.Context) {
	posts, err := getAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "posts": posts})
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
func likedislikepost(c *gin.Context) {
	var like struct {
		PostID string `json:"postID" binding:"required"`
		Liked  bool   `json:"checker" binding:"required"`
	}
	if err := c.ShouldBind(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	postID, err := strconv.Atoi(like.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid PostID"})
		return
	}
	token := c.GetHeader("Authorization")
	fmt.Println(token)
	userID, err := Query_ID(token)
	fmt.Println(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	action := LikeDislikeActions{
		UserID: userID,
		PostID: postID,
		IsLike: like.Liked,
	}
	err = HandleLikeDislike(action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
}
