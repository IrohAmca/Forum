package services

import (
	"database/sql"

	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func UploadImage(c *gin.Context) {
	// Veritabanı bağlantısını aç
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// Form verilerini al
	title := c.PostForm("title")
	content := c.PostForm("content")
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	defer file.Close()

	// Dosya boyutunu kontrol et
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	if len(fileBytes) > 20*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 20 MB"})
		return
	}

	// Veritabanına kaydet
	_, err = db.Exec("INSERT INTO posts (title, content, image) VALUES (?, ?, ?)", title, content, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database insert failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
