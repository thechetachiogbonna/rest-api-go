package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CategoryRoutes(pool *pgxpool.Pool) http.Handler {
	categoryHandler := &handlers.Category{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", categoryHandler.GetCategories)
	mux.HandleFunc("GET /{slug}", categoryHandler.GetCategoryBySlug)

	return mux
}
