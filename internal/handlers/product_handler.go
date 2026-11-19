package handlers

import (
	"net/http"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/types"
	"rest-api-go/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	Pool *pgxpool.Pool
}

func (p *Product) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := repository.GetAllProducts(p.Pool)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, map[string]any{"products": products}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (p *Product) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")

	product, err := repository.GetProductByID(p.Pool, productID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, map[string]any{"product": product}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (p *Product) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	product, err := repository.CreateProduct(
		p.Pool,
		payload.CategoryID,
		payload.Name,
		payload.Slug,
		payload.Description,
		payload.Price,
		payload.StockQuantity,
	)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]any{"product": product})
}
