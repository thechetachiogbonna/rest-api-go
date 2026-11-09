package handlers

import (
	"fmt"
	"net/http"
	"rest-api-go/internal/middleware"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/types"
	"rest-api-go/internal/utils"
	"strings"

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

func (user *User) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("Missing %q params", "id"))
		return
	}

	userData, err := repository.GetUserByID(user.Pool, userID)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("User not found"))
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var payload types.UpdateUserByIDPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	firstName := payload.FirstName
	if firstName == "" {
		firstName = userData.FirstName
	}
	lastName := payload.LastName
	if lastName == "" {
		lastName = userData.LastName
	}
	email := payload.Email
	if email == "" {
		email = userData.Email
	}

	updatedUserData, err := repository.UpdateUserByID(user.Pool, userID, firstName, lastName, email)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]any{"user": updatedUserData})
}
