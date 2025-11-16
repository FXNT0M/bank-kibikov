package models

import (
	"time"
)

type User struct {
	ID         int       `json:"id"`
	Cipher     string    `json:"cipher"`
	Password   string    `json:"-"`
	Name       string    `json:"name"`
	Group      string    `json:"group"`
	Balance    int       `json:"balance"`
	JoinedDate time.Time `json:"joined_date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRegisterRequest struct {
	Cipher   string `json:"cipher" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Group    string `json:"group" binding:"required"`
}

type UserLoginRequest struct {
	Cipher   string `json:"cipher" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID         int       `json:"id"`
	Cipher     string    `json:"cipher"`
	Name       string    `json:"name"`
	Group      string    `json:"group"`
	Balance    int       `json:"balance"`
	JoinedDate time.Time `json:"joined_date"`
}
