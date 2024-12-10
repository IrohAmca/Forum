package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frontend/manager"
	"frontend/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func setDeviceToken(token string) error {
	envPath := "config/.env"
	err := godotenv.Load(envPath)
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	envMap, err := godotenv.Read(envPath)
	if err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}
	envMap["DEVICE_TOKEN"] = token

	envFile, err := os.Create(envPath)
	if err != nil {
		return fmt.Errorf("error creating .env file: %v", err)
	}
	defer envFile.Close()

	for key, value := range envMap {
		_, err := envFile.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return fmt.Errorf("error writing to .env file: %v", err)
		}
	}

	fmt.Println("DEVICE_TOKEN set successfully")
	return nil
}
func DeviceRegister() bool {
	register := models.Register{}
	godotenv.Load("config/.env")
	register.Password = os.Getenv("DEVICE_PASSWORD")
	register.Device_Type = os.Getenv("DEVICE_TYPE")
	api := manager.API{}
	url := api.GetURL("DeviceRegister")
	register_marshalled, err := json.Marshal(register)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(register_marshalled))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	response := models.Register_CallBack{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	if response.Success {
		if err := setDeviceToken(response.Token); err != nil {
			log.Fatalf("Failed to set environment variable: %v\n", err)
		}
		print("Device Registered Successfully!!!\n")
		return true
	} else {
		print("Device Registration Failed!!!\n")
		return false
	}
}
