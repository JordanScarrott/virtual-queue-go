package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key int

const (
	UserKey key = iota
	RoleKey
)

// WithAuth is a middleware that validates JWT tokens
func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract Header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Validate Format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Parse Token
		token, err := jwt.ParseWithClaims(tokenString, &RedDuckClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what we expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("red-duck-secret-key-2026"), nil
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// 4. Verify
		if claims, ok := token.Claims.(*RedDuckClaims); ok && token.Valid {
			// 5. Context Injection
			ctx := context.WithValue(r.Context(), UserKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	}
}
