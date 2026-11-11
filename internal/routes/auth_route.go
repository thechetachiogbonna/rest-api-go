package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoutes(pool *pgxpool.Pool) *http.ServeMux {
	authHandler := &handlers.Auth{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", authHandler.LoginHandler)
	mux.HandleFunc("POST /register", authHandler.RegisterHandler)
	mux.HandleFunc("POST /logout", authHandler.LogoutHandler)
	mux.HandleFunc("GET /refresh", authHandler.RefreshHandler)

	return mux
}
