package handlers

import (
	"net/http"
	"rest-api-go/internal/middleware"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Pool *pgxpool.Pool
}

func (user *User) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	userData, err := repository.GetUserByID(user.Pool, userID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]any{"user": userData})
}
