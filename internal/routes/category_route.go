package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"
	"rest-api-go/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CategoryRoutes(pool *pgxpool.Pool) http.Handler {
	categoryHandler := &handlers.Category{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", categoryHandler.GetCategories)
	mux.HandleFunc("GET /{slug}", categoryHandler.GetCategoryBySlug)
	mux.Handle("POST /", middleware.CreateStack(
		middleware.Authenticate(pool),
		middleware.IsAdmin(pool),
	)(http.HandlerFunc(categoryHandler.CreateCategory)))

	return mux
}
