package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(pool *pgxpool.Pool) *http.ServeMux {
	productHandler := &handlers.Product{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", productHandler.GetAllProducts)

	return mux
}
