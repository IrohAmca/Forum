package models

type Register struct {
    Success bool   `json:"success"`
    Token   string `json:"token"`
}