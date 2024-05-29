package main

import (
	"errors"
	"fmt"
	"net/http"

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
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": false, "message": err.Error()})
		return
	}
	passwordReal, err := Query_email(loginInput.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": false, "message": err.Error()})
		return
	}

	err = authenticate(passwordReal, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": false, "message": "Invalid email or password"})
		return
	}
	username, _ := Query_username(loginInput.Email)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome " + username})
}

func SignUp(c *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": false, "message": err.Error()})
		return
	}
	updateUserID()
	if err := insertData(userID, user.Username, user.Email, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": false, "message": err.Error()})
		return
	}
	userID++
	c.JSON(http.StatusOK, gin.H{"success": true, "message": user.Username + " has been registered successfully"})
}

func authenticate(storedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		fmt.Println("Password does not match")
		return errors.New("authentication failed")
	}
	return nil
}
