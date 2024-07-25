package service

import (
	"forum/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetModarator(c *gin.Context) {
	var Mod struct {
		Username     string `json:"username" binding:"required"`
		Device_Token string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&Mod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if repository.CheckDeviceToken(Mod.Device_Token[:len(Mod.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	err := repository.SetMod(Mod.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User " + Mod.Username + " is now a modarator"})
}
