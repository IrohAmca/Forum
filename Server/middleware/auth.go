package middleware

import (
	"fmt"
	"forum/models"
	"forum/repository"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func AuthMiddleware(token string) (bool, error) {
	if repository.CheckDeviceToken(token) {
		return true, nil
	}
	return false, nil
}

func GenerateDeviceToken(device_type string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": device_type,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		panic("Failed to generate token")
	}
	return tokenString
}

func DeviceRegister(c *gin.Context) {
	var device struct {
		Password    string `json:"password" binding:"required"`
		Device_type string `json:"device_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(400, gin.H{"success": false, "message": "Reading Error: " + err.Error()})
		fmt.Println("Reading Error: ", err.Error())
		return
	}
	godotenv.Load("config/.env")
	password := os.Getenv("DEVICE_PASSWORD")
	if device.Password != password {
		c.JSON(400, gin.H{"success": false, "message": "Invalid Password"})
		fmt.Print("Invalid Password")
		return
	}
	token := GenerateDeviceToken(device.Device_type)
	err := repository.InsertDevice(device.Device_type, token)
	if err != nil {
		c.JSON(500, gin.H{"success": false, "message": err.Error()})
		fmt.Print(err.Error())
		return
	}
	register := models.Register{Success: true, Token: token}
	fmt.Print(register)
	c.JSON(200, register)
}
