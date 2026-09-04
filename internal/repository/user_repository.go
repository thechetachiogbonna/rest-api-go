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
		RETURNING id, first_name, last_name, email, created_at, updated_at
	`

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return nil, err
	}

	var user models.User

	err = pool.QueryRow(ctx, query, firstName, lastName, email, hashedPassword).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
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
