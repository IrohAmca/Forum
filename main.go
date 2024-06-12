package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	createDatabase()
	defer user_db.Close()
	WriteAllData()
	r := gin.Default()

	r.Static("/static", "./static")
	r.Static("/png", "./png")
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})
	r.GET("/get-posts", getPosts)
	r.POST("/sign-out", func(c *gin.Context) {
		c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
		c.JSON(200, gin.H{"success": true, "message": "You have been signed out"})
		c.Redirect(302, "/")
	})

	r.POST("/sign-up", SignUp)
	r.POST("/login", login)
	r.POST("/create-post", createPost)

	r.Run("localhost:8000")
}
