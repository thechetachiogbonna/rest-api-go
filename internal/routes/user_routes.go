package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(pool *pgxpool.Pool) *http.ServeMux {
	user := &handlers.User{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /me", user.GetCurrentUser)
	mux.HandleFunc("PUT /{id}", user.UpdateUserByID)

	return mux
}
