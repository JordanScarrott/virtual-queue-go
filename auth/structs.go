package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// User Struct
type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	BusinessID string `json:"business_id,omitempty"` // nullable/empty if customer
}

// RedDuckClaims Struct
type RedDuckClaims struct {
	jwt.RegisteredClaims
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	BusinessID string `json:"business_id"`
}

// API Contracts (for HTTP Handlers)
type LoginRequest struct {
	Email string `json:"email"`
}

type VerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
