package main

import "github.com/gin-gonic/gin"

func main() {
	// Assume these functions are defined elsewhere
	createDatabase()

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./")

	// Serve the sign-in page
	r.GET("/sign-in.html", func(c *gin.Context) {
		c.File("./sign-in.html")
	})

	// Serve the login page
	r.GET("/login", func(c *gin.Context) {
		c.File("./login.html")
	})

	// Handle the login form submission
	r.POST("/login", login)

	r.Run("localhost:8080")
}
