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
	r.POST("/sign-up", SignUp)
	r.POST("/login", login)
	r.Run("localhost:8080")
}
