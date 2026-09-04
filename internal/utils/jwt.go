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

func GetAccessTokenRegisterClaims() *jwt.RegisteredClaims {
	return &jwt.RegisteredClaims{
		Issuer:    ISSUER,
		Audience:  jwt.ClaimStrings{AUDIENCE},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(FIFTEEN_MINUTES)), // 15 minutes
	}
}

func GetRefreshTokenRegisterClaims() *jwt.RegisteredClaims {
	return &jwt.RegisteredClaims{
		Issuer:    ISSUER,
		Audience:  jwt.ClaimStrings{AUDIENCE},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(THIRTY_DAYS)), // 30 days
	}
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

func VerifyToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}
	return claims, nil
}
