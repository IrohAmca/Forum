package models

type Register_CallBack struct {
    Success bool   `json:"success"`
    Token   string `json:"token"`
}

type Register struct {
    Password string `json:"password"`
    Device_Type string `json:"device_type"`
}