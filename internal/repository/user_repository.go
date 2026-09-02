package repository

import (
	"context"
	"errors"
	"fmt"
	"rest-api-go/internal/models"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, firstName, lastName, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := GetUserByEmail(pool, email)

	if err == nil {
		return "", fmt.Errorf("User with email %s already exists", email)
	}

	query := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			password
		) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return "", err
	}

	var userID string

	err = pool.QueryRow(ctx, query, firstName, lastName, email, hashedPassword).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
