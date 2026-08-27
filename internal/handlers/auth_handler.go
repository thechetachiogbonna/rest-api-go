package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api-go/internal/models"
	"rest-api-go/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	Pool *pgxpool.Pool
}

func (auth *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (auth *Auth) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := repository.CreateUser(auth.Pool, payload.FirstName, payload.LastName, payload.Email, payload.Password)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
}

func (auth *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
