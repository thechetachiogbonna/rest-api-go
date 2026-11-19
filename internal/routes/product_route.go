package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"
	"rest-api-go/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(pool *pgxpool.Pool) *http.ServeMux {
	productHandler := &handlers.Product{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", productHandler.GetAllProducts)
	mux.HandleFunc("GET /{id}", productHandler.GetProductByID)
	mux.Handle("POST /", middleware.CreateStack(
		middleware.Authenticate(productHandler.Pool),
		middleware.IsAdmin(productHandler.Pool),
	)(http.HandlerFunc(productHandler.CreateProduct)))

	return mux
}
