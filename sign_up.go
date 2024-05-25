package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Email           string
	Age             int
	Password        string
	ConfirmPassword string
}

func checkSignUp(c *gin.Context) {
	fmt.Println("Sign up request received")
	var loginInput struct {
		Email           string `form:"email" binding:"required"`
		Password        string `form:"password" binding:"required"`
		ConfirmPassword string `form:"confirm_password" binding:"required"`
		Age             int    `form:"age" binding:"required"`
	}
	if loginInput.Password != loginInput.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}
	if loginInput.Age < 18 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must be at least 18 years old to sign up"})
		return
	}
	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertData(userID, loginInput.Email, loginInput.Password, loginInput.Age)
	userID++
	c.JSON(http.StatusOK, gin.H{"message": "User successfully created"})
}
