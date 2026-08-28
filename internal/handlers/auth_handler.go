package handlers

import (
	"net/http"
	"rest-api-go/internal/models"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/utils"

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
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := repository.CreateUser(auth.Pool, payload.FirstName, payload.LastName, payload.Email, payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "User Created Successfully."})
}

func (auth *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
