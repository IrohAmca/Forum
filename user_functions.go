package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password")
	}
	return string(hashedPassword)
}

func login(c *gin.Context) {
	var loginInput struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emailReal, passwordReal, err := Query_email(loginInput.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = authenticate(passwordReal, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": emailReal})
}

func authenticate(storedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		fmt.Println("Password does not match")
		return errors.New("authentication failed")
	}
	return nil
}
func SignUp(c *gin.Context) {
	var user struct {
		Email           string `form:"email" binding:"required"`
		Password        string `form:"password" binding:"required"`
		ConfirmPassword string `form:"confirm_password" binding:"required"`
		Age             int    `form:"age" binding:"required"`
	}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Age, _ = strconv.Atoi(c.PostForm("age"))
	c.JSON(http.StatusOK, gin.H{"message": user.ConfirmPassword})
	c.JSON(http.StatusOK, gin.H{"message": user.Password})
	if user.Password != user.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	if user.Age < 18 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must be at least 18 years old to sign up"})
		return
	}
	updateUserID()
	insertData(userID, user.Email, user.Password, user.Age)
	userID++

	c.JSON(http.StatusOK, gin.H{"message": "User successfully created"})
}