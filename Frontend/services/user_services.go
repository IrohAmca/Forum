package services

import (
	"bytes"
	"encoding/json"
	"frontend/manager"
	"frontend/models"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func ProfilePage(c *gin.Context) {
	var user struct {
		Username     string `json:"username" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	w := c.Writer
	user.Username = c.Params.ByName("username")
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user.Device_Token = env["DEVICE_TOKEN"]
	postData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("ProfilePage")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))
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
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    models.UserData `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	tpl, err := template.ParseFiles("templates/userprofile.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	tpl.Execute(w, &response.Data)
}

func SignUp(c *gin.Context) {
	var user struct {
		Username     string `json:"username" binding:"required"`
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user.Device_Token = env["DEVICE_TOKEN"]
	postData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("SignUp")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))
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
	success := response.Success
	message := response.Message
	c.JSON(http.StatusOK, gin.H{"success": success, "message": message})

}
func Login(c *gin.Context) {
	var loginInput struct {
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	loginInput.Device_Token = env["DEVICE_TOKEN"]
	loginData, err := json.Marshal(loginInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("Login")
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
		Token   string `json:"token"`
		Cookie  string `json:"cookie"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.SetCookie("cookie", response.Cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", response.Token)
	c.JSON(http.StatusOK, gin.H{"success": response.Success, "message": response.Message})
}
func UserChecker(c *gin.Context) {
	var user struct {
		Cookie       string `json:"cookie" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user.Device_Token = env["DEVICE_TOKEN"]
	postData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("UserChecker")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))
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
		Success  bool   `json:"success"`
		Message  string `json:"message"`
		Username string `json:"username"`
		Token    string `json:"token"`
		Userlevel string `json:"userlevel"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"success": response.Success, "message": response.Message, "username": response.Username, "token": response.Token, "userlevel": response.Userlevel})
}
func SignOut(c *gin.Context) {
	var user struct {
		Cookie       string `json:"cookie"`
		Device_Token string `json:"device_token"`
	}

	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user.Device_Token = env["DEVICE_TOKEN"]
	user.Cookie, err = c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	postData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("SignOut")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))
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
	c.SetCookie("cookie", "", -1, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"success": response.Success, "message": response.Message})
}
