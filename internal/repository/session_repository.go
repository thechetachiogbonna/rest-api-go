package repository

import (
	"context"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateSession(pool *pgxpool.Pool, userId, ipAddress, userAgent string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO sessions (user_id, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	expiresAt := time.Now().Add(utils.THIRTY_DAYS)

	var sessionID string
	err := pool.QueryRow(ctx, query, userId, ipAddress, userAgent, expiresAt).Scan(&sessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
