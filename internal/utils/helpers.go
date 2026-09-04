package utils

import (
	"rest-api-go/internal/config"
)

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
