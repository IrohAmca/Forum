package main

import (
	"frontend/setup"
	"frontend/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	setup.Setup()

	// Default route
	r.Static("/static", "./static")
	r.Static("/png", "./png")
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	// User routes
	r.GET("/profile/:username", services.ProfilePage)
	r.GET("/auth/google", services.GoogleLogin)
	r.GET("/auth/google/callback", services.GoogleCallback)
	r.GET("/auth/github", services.GithubLogin)
	r.GET("/auth/github/callback", services.GithubCallback)
	r.GET("/auth/facebook", services.FacebookLogin)
	r.GET("/auth/facebook/callback", services.FacebookCallback)
	r.GET("/admin", services.AdminPage)
	r.GET("/moderator/:username", services.ModeratorPage)
	r.POST("/setModarator", services.SetModarator)
	r.POST("/sign-up", services.SignUp)
	r.POST("/login", services.Login)
	r.POST("/check-token", services.UserChecker)
	r.POST("/sign-out", services.SignOut)

	// Post routes
	r.POST("/get-posts", services.GetPosts)
	r.POST("ld_comment", services.LikeDislikeComment)
	r.POST("/delete-comment", services.DeleteComment)
	r.POST("/ld_post", services.LikeDislikePost)
	r.POST("/create-post", services.CreatePost)
	r.POST("/delete-post", services.DeletePost)
	r.POST("/create-comment", services.CreateComment)

	// Thread routes
	r.Run(":8080")
}
