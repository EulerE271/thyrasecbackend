package models

import "github.com/dgrijalva/jwt-go"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID   string `json:"user_id"` // Change to string for UUID
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}
