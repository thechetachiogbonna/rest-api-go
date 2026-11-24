package handlers

import (
	"context"
	"errors"
	"net/http"
	"rest-api-go/internal/types"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cart struct {
	Pool *pgxpool.Pool
}

func (c *Cart) GetCurrentUserCartItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
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

	rows, err := c.Pool.Query(ctx, query, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
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
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		item.Product.ID = item.ProductID
		item.ItemTotal = item.Product.Price * item.Quantity
		totalCost += item.ItemTotal
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, map[string]any{
		"items":      items,
		"total_cost": totalCost,
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}
