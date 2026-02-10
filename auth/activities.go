package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.temporal.io/sdk/activity"
)

// SendMagicCode sending an email with the magic code
func SendMagicCode(ctx context.Context, email string, code string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("MAGIC LOGIN CODE FOR " + email + ": " + code)
	return nil
}

// GenerateToken creates a signed JWT for the session
func GenerateToken(ctx context.Context, user User) (string, error) {
	claims := RedDuckClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days
		},
		UserID:     user.ID,
		Email:      user.Email,
		Role:       user.Role,
		BusinessID: user.BusinessID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Secret key (hardcoded for now as per requirements)
	secretKey := []byte("red-duck-secret-key-2026")

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
