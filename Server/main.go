package main

import (
	"forum/middleware"
	"forum/repository"
	"forum/service"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.CreateDatabase()
	// db_manager.CloseDatabase() <-- Can be add
	// writeAllData()
	r := gin.Default()
	repository.SetAdmin("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjEwNDUzNzMsInN1YiI6IkFsaSBBeWfDvG4ifQ.PECOF4yo6cxKSgm6DgBmk3e-ueBC1fn9aW7q0UoiSDk")
	r.POST("/device-register", middleware.DeviceRegister)
	
	r.POST("/check-email", service.CheckEmail)
	r.POST("/auth/login", service.AuthLogin)
	r.POST("/auth/signup", service.AuthSignup)

	r.POST("/profile", service.ProfilePage)
	r.POST("/sign-up", service.SignUp)
	r.POST("/login", service.Login)
	r.POST("/check-token", service.UserChecker)
	r.POST("/sign-out", service.SignOut)

	// Post routes
	r.POST("/get-posts", service.GetPosts)
	r.POST("ld_comment", service.LikeDislikeComment)
	r.POST("/delete-comment", service.DeleteComment)
	r.POST("/ld_post", service.LikeDislikePost)
	r.POST("/create-post", service.CreatePost)
	r.POST("/delete-post", service.DeletePost)
	r.POST("/create-comment", service.CreateComment)

	r.GET("/images/:filename", service.GetImage)

	// Manager routes
	r.POST("/setModarator", service.SetModarator)
	// Thread routes
	r.Run(":8081")
}
