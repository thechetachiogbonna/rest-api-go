package repository

import (
	"context"
	"rest-api-go/internal/models"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateSession(pool *pgxpool.Pool, userId, ipAddress, userAgent string) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO sessions (user_id, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, ip_address, user_agent, expires_at, created_at
	`

	expiresAt := time.Now().Add(utils.THIRTY_DAYS)

	var session models.Session
	err := pool.QueryRow(ctx, query, userId, ipAddress, userAgent, expiresAt).Scan(
		&session.ID,
		&session.UserID,
		&session.IPAddress,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
