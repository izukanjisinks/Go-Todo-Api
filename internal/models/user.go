package models

import "github.com/google/uuid"

type Login struct {
	Password string `json:"-"`
}

type User struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	IsAdmin  bool      `json:"is_admin"`
	Login
}
