package models

import (
	"time"
)

type UserPayload struct {
	UserName          string `json:"user_name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserName  string    `gorm:"unique" json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
