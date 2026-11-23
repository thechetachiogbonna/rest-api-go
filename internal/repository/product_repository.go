package repository

import (
	"context"
	"errors"
	"rest-api-go/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllProducts(pool *pgxpool.Pool) ([]*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, category_id, name, slug, description, price, stock_quantity, created_at, updated_at
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
			&product.Description,
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
		SELECT id, category_id, name, slug, description, price, stock_quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &models.Product{}

	if err := pool.QueryRow(ctx, query, productID).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return product, nil
}

func CreateProduct(pool *pgxpool.Pool, categoryID, name, slug, description string, price, stockQuantity int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (category_id, name, slug, description, price, stock_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, category_id, name, slug, description, price, stock_quantity, created_at, updated_at
	`
	product := &models.Product{}
	err := pool.QueryRow(ctx, query, categoryID, name, slug, description, price, stockQuantity).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func UpdateProductByID(pool *pgxpool.Pool, productID, categoryID, name, slug, description string, price, stockQuantity int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE products
		SET category_id = $2, name = $3, slug = $4, description = $5, price = $6, stock_quantity = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, category_id, name, slug, description, price, stock_quantity, created_at, updated_at
	`
	product := &models.Product{}
	err := pool.QueryRow(ctx, query, productID, categoryID, name, slug, description, price, stockQuantity).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func DeleteProduct(pool *pgxpool.Pool, productID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM products
		WHERE id = $1
	`
	_, err := pool.Query(ctx, query, productID)

	if err != nil {
		return err
	}

	return nil
}
