package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func find_catagory(c *gin.Context) {
	var post struct {
		Content string `form:"email" binding:"required"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content := strings.Split(post.Content, " ")
	for _, word := range content {
		if strings.HasPrefix(word, "#") {
			if word[1:] == "sports" {
				//
			}
			if word[1:] == "politics" {
				//
			}
			if word[1:] == "technology" {
				//
			}
			if word[1:] == "entertainment" {
				//
			}
			if word[1:] == "food" {
				//
			}
			if word[1:] == "travel" {
				//
			}
			if word[1:] == "fashion" {
				//
			}
		}
	}
}

func createPost(c *gin.Context) {
	var post struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	userIDstr, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userID, _ := strconv.Atoi(userIDstr)
	threadID, err := insertThread(userID, 0, post.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	insertPost(threadID, userID, post.Content)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post created successfully"})
}

func getPosts(c *gin.Context){
	posts,err :=getAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "posts": posts})
}