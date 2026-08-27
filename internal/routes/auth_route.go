package routes

import (
	"net/http"
	"rest-api-go/internal/handlers"
)

func AuthRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /register", handlers.RegisterHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)

	return mux
}
