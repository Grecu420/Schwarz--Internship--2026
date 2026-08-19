package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT signed with a secret key that expires after duration.
func GenerateToken(userID int64, duration time.Duration, issuer string, secretKey []byte) (string, error) {
	// Define claims
	claims := LoginClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func ParseToken(tokenStr string, secret []byte) (*LoginClaims, error) {

	claims := LoginClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return &claims, nil
}
