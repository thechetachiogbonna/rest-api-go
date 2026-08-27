package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoutes(pool *pgxpool.Pool) *http.ServeMux {
	auth := &handlers.Auth{Pool: pool}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", auth.LoginHandler)
	mux.HandleFunc("POST /register", auth.RegisterHandler)
	mux.HandleFunc("POST /logout", auth.LogoutHandler)

	return mux
}
