package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}

type Claims struct {
	Payload
	jwt.RegisteredClaims
}

var AccessTokenRegisterClaims = &jwt.RegisteredClaims{
	Issuer:    ISSUER,
	Audience:  jwt.ClaimStrings{AUDIENCE},
	IssuedAt:  jwt.NewNumericDate(time.Now()),
	ExpiresAt: jwt.NewNumericDate(time.Now().Add(FIFTEEN_MINUTES)), // 15 minutes
}

var RefreshTokenRegisterClaims = &jwt.RegisteredClaims{
	Issuer:    ISSUER,
	Audience:  jwt.ClaimStrings{AUDIENCE},
	IssuedAt:  jwt.NewNumericDate(time.Now()),
	ExpiresAt: jwt.NewNumericDate(time.Now().Add(THIRTY_DAYS)), // 30 days
}

func GenerateToken(payload Payload, registerClaims *jwt.RegisteredClaims, secret string) (string, error) {
	claims := Claims{
		Payload:          payload,
		RegisteredClaims: *registerClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
