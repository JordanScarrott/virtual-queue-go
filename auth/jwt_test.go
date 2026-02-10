package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	user := User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Role:  "admin",
	}

	tokenString, err := GenerateToken(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse back to verify
	token, err := jwt.ParseWithClaims(tokenString, &RedDuckClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("red-duck-secret-key-2026"), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*RedDuckClaims)
	assert.True(t, ok)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Role, claims.Role)
}

func TestWithAuth(t *testing.T) {
	// Mock handler
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserKey).(string)
		role := r.Context().Value(RoleKey).(string)

		// Assertions inside handler to verify context injection
		if userID == "" || role == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	t.Run("No Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		rr := httptest.NewRecorder()

		handler := WithAuth(mockHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Bad Header Format", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Basic xyz")
		rr := httptest.NewRecorder()

		handler := WithAuth(mockHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Generate a real token
		user := User{ID: "uid-123", Role: "user"}
		token, _ := GenerateToken(context.Background(), user)

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		handler := WithAuth(mockHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
