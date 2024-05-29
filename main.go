package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	createDatabase()
	WriteAllData()
	r := gin.Default()

	r.Static("/static", "templates/index.html")

	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	r.POST("/sign-up", SignUp)
	r.POST("/login", login)

	r.Run("localhost:8080")
}
