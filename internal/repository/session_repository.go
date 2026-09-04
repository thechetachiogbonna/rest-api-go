package repository

import (
	"context"
	"errors"
	"rest-api-go/internal/models"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
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

func GetSessionByUserIDAndSessionID(pool *pgxpool.Pool, userID, sessionID string) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, ip_address, user_agent, expires_at
		FROM sessions
		WHERE id = $1 AND user_id = $2
	`

	var session models.Session
	err := pool.QueryRow(ctx, query, sessionID, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.IPAddress,
		&session.UserAgent,
		&session.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}

	return &session, nil
}
