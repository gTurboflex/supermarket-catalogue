package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
