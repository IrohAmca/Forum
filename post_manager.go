package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
