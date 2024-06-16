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
	r.POST("/get-posts", getPosts)

	r.POST("/sign-out", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "localhost", false, false)
		c.JSON(200, gin.H{"success": true, "message": "You have been signed out"})
		c.Redirect(302, "/")
	})
	r.POST("/delete-comment", deleteComment)
	r.POST("/sign-up", SignUp)
	r.POST("/login", login)
	r.POST("/ld_post",likeDislikePost)
	r.POST("/create-post", createPost)
	r.POST("/check-token", UserChecker)
	r.POST("/delete-post", deletePost)
	r.POST("/create-comment", createComment)
	r.Run("localhost:8080")
}
