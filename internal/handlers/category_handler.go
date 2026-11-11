package handlers

import (
	"fmt"
	"net/http"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/utils"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	Pool *pgxpool.Pool
}

func (c *Category) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := repository.GetCategories(c.Pool)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJson(w, http.StatusOK, map[string]any{"categories": categories}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (c *Category) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	category, err := repository.GetCategoryBySlug(c.Pool, slug)

	if err != nil {
		if strings.Contains(err.Error(), "category not found") {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("Category not found"))
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]any{"category": category})
}
