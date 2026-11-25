package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CartRoutes(pool *pgxpool.Pool) *http.ServeMux {
	cartHandler := &handlers.Cart{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", cartHandler.GetCurrentUserCartItems)
	mux.HandleFunc("POST /items", cartHandler.AddItemToCart)

	return mux
}
