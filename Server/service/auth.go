package service

import (
	"fmt"
	"forum/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthLogin(c *gin.Context) {
	var user struct {
		Email        string `json:"email" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if repository.CheckDeviceToken(user.Device_Token[:len(user.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	username, _ := repository.Query_username(user.Email)
	token, err := repository.QueryToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	cookie := GenerateCookie(token)

	CheckSession(token, c)
	err = repository.InsertSession(token, cookie)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome " + username, "cookie": cookie})
}

func AuthSignup(c *gin.Context) {
	var user struct {
		Username     string `json:"username" binding:"required"`
		Email        string `json:"email" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Reading Error" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if repository.CheckDeviceToken(user.Device_Token[:len(user.Device_Token)-8]) {
		fmt.Println("Invalid Device Token")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	token := GenerateToken(user.Username)
	err := repository.InsertAuthUser(user.Username, user.Email, token)
	if err != nil {
		fmt.Println("Insert Error" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome " + user.Username})

}
