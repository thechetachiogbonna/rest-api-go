package repository

import (
	"context"
	"fmt"
	"rest-api-go/internal/models"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, firstName, lastName, email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := GetUserByEmail(pool, email)

	if err == nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	query := `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)`

	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, query, firstName, lastName, email, hashedPassword)

	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `
		SELECT * FROM users WHERE email = $1
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
