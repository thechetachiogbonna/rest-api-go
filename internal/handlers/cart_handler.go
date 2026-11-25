package handlers

import (
	"errors"
	"net/http"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/types"
	"rest-api-go/internal/utils"

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

	items, totalCost, err := repository.GetCurrentUserCartItems(c.Pool, userID)

	if err != nil {
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

func (c *Cart) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(string)
	if !ok || userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	var payload types.AddItemToCartPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := repository.AddItemToCart(c.Pool, userID, payload.ProductID, payload.Quantity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "product added to cart successfully."})
}
