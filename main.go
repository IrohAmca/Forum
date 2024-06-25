package main

import (
	"forum/db_manager"
	"forum/services"
	"github.com/gin-gonic/gin"
)

func main() {
	db_manager.CreateDatabase()
	// db_manager.CloseDatabase() <-- Can be add
	// writeAllData()
	r := gin.Default()

	// Default route
	r.Static("/static", "./static")
	r.Static("/png", "./png")
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	// User routes
	r.GET("/profile/:username", services.ProfilePage)

	r.POST("/sign-up", services.SignUp)
	r.POST("/login", services.Login)
	r.POST("/check-token", services.UserChecker)
	r.POST("/sign-out", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "localhost", false, false)
		c.JSON(200, gin.H{"success": true, "message": "You have been signed out"})
		c.Redirect(302, "/")
	})

	// Post routes
	r.POST("/get-posts", services.GetPosts)
	r.POST("ld_comment", services.LikeDislikeComment)
	r.POST("/delete-comment", services.DeleteComment)
	r.POST("/ld_post", services.LikeDislikePost)
	r.POST("/create-post", services.CreatePost)
	r.POST("/delete-post", services.DeletePost)
	r.POST("/create-comment", services.CreateComment)

	// Thread routes
	r.Run("localhost:8080")
}
