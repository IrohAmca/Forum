package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/manager"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AdminPage(c *gin.Context) {
	c.File("templates/admin.html")
}
func ModeratorPage(c *gin.Context) {
	c.File("templates/moderator.html")
}


func SetModarator(c *gin.Context) {
	var Mod struct {
		Username     string `json:"username" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBindJSON(&Mod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	fmt.Println(Mod)
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	Mod.Device_Token = env["DEVICE_TOKEN"]
	loginData, err := json.Marshal(Mod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("SetModarator")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(loginData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": response.Success, "message": response.Message})
}
