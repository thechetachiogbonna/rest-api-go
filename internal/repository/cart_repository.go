package repository

import (
	"context"
	"errors"
	"fmt"
	"rest-api-go/internal/types"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCurrentUserCartItems(pool *pgxpool.Pool, userID string) ([]types.CartItem, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
			p.name, p.slug, p.price
		FROM cart_items ci
		JOIN carts c ON ci.cart_id = c.id
		JOIN products p ON ci.product_id = p.id
		WHERE c.user_id = $1
		ORDER BY ci.created_at, ci.id
	`

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]types.CartItem, 0)
	totalCost := 0
	for rows.Next() {
		var item types.CartItem
		if err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Product.Name,
			&item.Product.Slug,
			&item.Product.Price,
		); err != nil {
			return nil, 0, err
		}

		item.Product.ID = item.ProductID
		item.ItemTotal = item.Product.Price * item.Quantity
		totalCost += item.ItemTotal
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, totalCost, nil
}

func AddItemToCart(pool *pgxpool.Pool, userID, productID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		query := `
			UPDATE products
			SET stock_quantity = stock_quantity - $1, updated_at = NOW()
			WHERE id = $2 AND stock_quantity >= $1
		`
		result, err := tx.Exec(ctx, query, quantity, productID)
		if err != nil {
			return fmt.Errorf("failed to update stock: %w", err)
		}
		if result.RowsAffected() == 0 {
			product, err := GetProductByID(ctx, tx, productID)
			if err != nil {
				return err
			}

			return fmt.Errorf("insufficient stock: requested %d, available %d", quantity, product.StockQuantity)
		}

		var cartID string
		query = `
			INSERT INTO carts (user_id)
			VALUES ($1)
			ON CONFLICT (user_id) DO UPDATE
			SET user_id = EXCLUDED.user_id
			RETURNING id
		`
		err = tx.QueryRow(ctx, query, userID).Scan(&cartID)
		if err != nil {
			return fmt.Errorf("failed to insert cart: %w", err)
		}

		query = `
			INSERT INTO cart_items (cart_id, product_id, quantity)
			VALUES ($1, $2, $3)
			ON CONFLICT (cart_id, product_id) DO UPDATE
			SET quantity = cart_items.quantity + EXCLUDED.quantity,
				updated_at = NOW()
		`
		_, err = tx.Exec(ctx, query, cartID, productID, quantity)
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %w", err)
		}

		return nil
	})
}
