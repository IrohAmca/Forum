package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/manager"
	"frontend/models"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func load_env() {
	err := godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
}
func GetPosts(c *gin.Context) {
	load_env()
	var Post struct {
		Categories   []string `json:"categories"`
		Title        string   `json:"title"`
		Short_type   string   `json:"short_type"`
		Device_Token string   `json:"device_token"`
	}
	if err := c.ShouldBind(&Post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	Post.Device_Token = env["DEVICE_TOKEN"]
	postData, err := json.Marshal(Post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("GetPosts")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	var response struct {
		Success bool          `json:"success"`
		Post    []models.Post `json:"posts"`
		Message string        `json:"message"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	success := response.Success
	post := response.Post
	message := response.Message
	c.JSON(http.StatusOK, gin.H{"success": success, "posts": post, "message": message})
}

func LikeDislikePost(c *gin.Context) {
	var like struct {
		PostID       string `json:"PostID" binding:"required"`
		IsLike       bool   `json:"isLike"`
		Cookie       string `json:"cookie"`
		Device_Token string `json:"device_token"`
	}

	if err := c.ShouldBind(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	like.Device_Token = env["DEVICE_TOKEN"]
	like.Cookie, err = c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	postData, err := json.Marshal(like)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	api := manager.API{}
	url := api.GetURL("LikeDislikePost")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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

func LikeDislikeComment(c *gin.Context) {
	var like struct {
		CommentID    string `json:"CommentID" binding:"required"`
		IsLike       bool   `json:"isLike"`
		Cookie       string `json:"cookie"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBind(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	like.Device_Token = env["DEVICE_TOKEN"]
	like.Cookie, err = c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	commentData, err := json.Marshal(like)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("LikeDislikeComment")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(commentData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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

func DeleteComment(c *gin.Context) {
	var comment struct {
		CommentID    string `json:"CommentID" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	comment.Device_Token = env["DEVICE_TOKEN"]
	commentData, err := json.Marshal(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("DeleteComment")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(commentData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
	if !response.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"success": response.Success, "message": response.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": response.Success, "message": response.Message})
}

func file2Blob(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	return ioutil.ReadAll(src)
}

func CreatePost(c *gin.Context) {
	var post struct {
		Title        string
		Content      string
		Categories   []string
		Image        []byte
		Ext          string
		Device_Token string
		Cookie       string
	}
	post.Title = c.PostForm("title")
	post.Content = c.PostForm("content")
	post.Categories = c.PostFormArray("categories")
	var err error
	if post.Title == "" || post.Content == "" || len(post.Categories) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Please fill all the fields"})
		return
	}
	file, _ := c.FormFile("image")
	if file != nil {
		if file.Size > 1024*1024*20 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Image size should be less than 20MB"})
			return
		}
		post.Ext = filepath.Ext(file.Filename)
		post.Image, err = file2Blob(file)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Please upload an image"})
			return
		}
	}
	post.Cookie, err = c.Cookie("cookie")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Please log in to the website first!!!"})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	post.Device_Token = env["DEVICE_TOKEN"]
	postData, err := json.Marshal(post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	api := manager.API{}
	url := api.GetURL("CreatePost")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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

func DeletePost(c *gin.Context) {
	var post struct {
		PostID       string `json:"PostID" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	post.Device_Token = env["DEVICE_TOKEN"]
	postData, err := json.Marshal(post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("DeletePost")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(postData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
	if !response.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"success": response.Success, "message": response.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": response.Success, "message": response.Message})
}

func CreateComment(c *gin.Context) {
	var comment struct {
		PostID       string `json:"postId" binding:"required"`
		Content      string `json:"content" binding:"required"`
		Cookie       string `json:"cookie" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reading Error" + err.Error()})
		return
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	comment.Device_Token = env["DEVICE_TOKEN"]
	commentData, err := json.Marshal(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	api := manager.API{}
	url := api.GetURL("CreateComment")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(commentData))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
	if !response.Success {
		c.JSON(http.StatusInternalServerError, gin.H{"success": response.Success, "message": response.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": response.Success, "message": response.Message})
}
