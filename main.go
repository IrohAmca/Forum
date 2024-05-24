package main

import "github.com/gin-gonic/gin"

func main() {
	// Assume these functions are defined elsewhere
	createDatabase()
	deleteData(1)
	insertData(1, "sebo@gmail.com", "got123")
	Query(1)

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./")

	// Serve the login page
	r.GET("/login", func(c *gin.Context) {
		c.File("./login.html")
	})

	// Handle the login form submission
	r.POST("/login", login)

	r.Run("localhost:8080")
}
