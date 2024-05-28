package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	EmployeeID string `json:"user_id"`
}

type Authenticator struct {
	SignKey  string
	TokenTTL time.Duration
}

func NewAuthenticator(signKey string, tokenTTL time.Duration) *Authenticator {
	return &Authenticator{SignKey: signKey, TokenTTL: tokenTTL}
}

func (a *Authenticator) ParseToken(token string) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(a.SignKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token with claims: %w", err)
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse token claims")
	}

	return claims.EmployeeID, nil
}

func (a *Authenticator) GenerateToken(employeeID uuid.UUID) (string, error) {
	tokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		employeeID.String(),
	})

	token, err := tokenRaw.SignedString([]byte(a.SignKey))

	if err != nil {
		return "", fmt.Errorf("failed to sign token with claims: %w", err)
	}

	return token, nil
}
