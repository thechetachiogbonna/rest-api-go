package repository

import (
	"context"
	"rest-api-go/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllProducts(pool *pgxpool.Pool) ([]*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, category_id, name, slug, price, stock_quantity, created_at, updated_at
		FROM products
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Product{}

	for rows.Next() {
		product := &models.Product{}
		if err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Slug,
			&product.Price,
			&product.StockQuantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductByID(pool *pgxpool.Pool, productID string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, category_id, name, slug, price, stock_quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &models.Product{}

	if err := pool.QueryRow(ctx, query, productID).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Price,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return product, nil
}
