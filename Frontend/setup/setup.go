package setup

import (
	"os"
	"frontend/auth"
	"github.com/joho/godotenv"
)

func Setup() {
	godotenv.Load("config/.env")
	env ,err:= godotenv.Read("config/.env"); if err != nil {
		print("Error reading .env file\n")
		os.Exit(1)
	}
	token := env["DEVICE_TOKEN"]
	if token == "" {
		if auth.DeviceRegister(){
			print("Device Registration Successful!!!\n")
		}else{
			print("Device Registration Failed!!!\n")
			print("Please check your internet connection and password, try again \n")
			os.Exit(1)
		}
	}
}
