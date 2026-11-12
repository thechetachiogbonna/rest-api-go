package utils

import (
	"rest-api-go/internal/config"
	"time"
)

const FIFTEEN_MINUTES = time.Minute * 15
const THIRTY_DAYS = time.Hour * 24 * 30
const ISSUER = "Rest-API-Go"
const AUDIENCE = "user"
const REFRESH_PATH = "/api/auth/refresh"

func GetAccessAndRefreshTokens(userID, sessionID string) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateToken(Payload{
		UserID:    userID,
		SessionID: sessionID,
	}, GetAccessTokenRegisterClaims(), config.Env.JWT_SECRET)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateToken(Payload{
		UserID:    userID,
		SessionID: sessionID,
	}, GetRefreshTokenRegisterClaims(), config.Env.JWT_REFRESH_SECRET)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func IsUserRoleAdmin(userRole string) bool {
	return userRole == "admin"
}
