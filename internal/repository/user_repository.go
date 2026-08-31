package repository

import (
	"context"
	"fmt"
	"rest-api-go/internal/models"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, firstName, lastName, email, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := GetUserByEmail(pool, email)

	if err == nil {
		return nil, fmt.Errorf("User with email %s already exists", email)
	}

	query := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			password
		) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, email, first_name, last_name, created_at, updated_at
	`

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return nil, err
	}

	var user models.User

	err = pool.QueryRow(ctx, query, firstName, lastName, email, hashedPassword).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		SELECT id, email, first_name, last_name, created_at, updated_at
		FROM users WHERE email = $1
	`

	user := models.User{}

	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
