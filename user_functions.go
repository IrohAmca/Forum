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
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	passwordReal, err := Query_email(loginInput.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	err = authenticate(passwordReal, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid email or password"})
		return
	}
	username, _ := Query_username(loginInput.Email)
	userID, _ := getUserID(loginInput.Email)
	c.SetCookie("user_id", strconv.Itoa(userID), 3600, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome " + username})
}

func SignUp(c *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := insertUser(user.Username, user.Email, user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
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

func getinformation(c *gin.Context) {
	useridstr, err := c.Cookie("user_id")
	fmt.Println(useridstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	userid, err := strconv.Atoi(useridstr)
	data, err := postDataByUserID(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "success": true, "data":data})
}
