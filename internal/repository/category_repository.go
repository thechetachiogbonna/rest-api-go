package repository

import (
	"context"
	"errors"
	"rest-api-go/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCategories(pool *pgxpool.Pool) ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, slug, created_at, updated_at
		FROM categories
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		category := &models.Category{}
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategoryBySlug(pool *pgxpool.Pool, slug string) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, slug, created_at, updated_at
		FROM categories
		WHERE slug = $1
	`
	category := &models.Category{}

	err := pool.QueryRow(ctx, query, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return category, nil
}

func CreateCategory(pool *pgxpool.Pool, name, slug string) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO categories (name, slug, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, name, slug, created_at, updated_at
	`
	category := &models.Category{}
	err := pool.QueryRow(ctx, query, name, slug).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return category, nil
}
