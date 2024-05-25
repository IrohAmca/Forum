package main

import "github.com/gin-gonic/gin"

func main() {
	createDatabase()
	listUsers()
	r := gin.Default()

	r.Static("/static", "./index.html")

	r.GET("/sign-up", func(c *gin.Context) {
		c.File("./sign-up.html")
	})
	r.POST("/sign-up", SignUp)

	r.GET("/login", func(c *gin.Context) {
		c.File("./login.html")
	})
	
	r.GET("/", func(c *gin.Context) {
		c.File("./index.html")
	})

	r.POST("/login", login)

	r.Run("localhost:8080")
}
