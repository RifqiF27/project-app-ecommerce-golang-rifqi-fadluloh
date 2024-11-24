package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty" validate:"required,min=3,regex=^[A-Za-z ]+$"`
	Email     string    `json:"email,omitempty" validate:"required_without=Phone,omitempty,email"`
	Phone     string    `json:"phone,omitempty" validate:"required_without=Email,omitempty,numeric,min=10,max=13"`
	Address   []string  `json:"address,omitempty"`
	Password  string    `json:"password,omitempty" validate:"required,min=8"`
	Token     string    `json:"token,omitempty"`
	UpdatedAt time.Time `json:"-"`
}

type Session struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
