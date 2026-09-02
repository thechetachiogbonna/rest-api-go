package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/types"
	"rest-api-go/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	Pool *pgxpool.Pool
}

func (auth *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := repository.GetUserByEmail(auth.Pool, payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid email or password."))
		return
	}

	err = utils.ComparePasswords(user.Password, payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("Invalid email or password."))
		return
	}

	sessionID, err := repository.CreateSession(auth.Pool, user.ID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, refreshToken, err := utils.GetAccessAndRefreshTokens(user.ID, sessionID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.SetAuthCookies(w, accessToken, refreshToken)

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "User logged in successfully."})
}

func (auth *Auth) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := repository.CreateUser(auth.Pool, payload.FirstName, payload.LastName, payload.Email, payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	sessionID, err := repository.CreateSession(auth.Pool, userID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, refreshToken, err := utils.GetAccessAndRefreshTokens(userID, sessionID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.SetAuthCookies(w, accessToken, refreshToken)

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": fmt.Sprintf("User created successfully with id %q.", userID)})
}

func (auth *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
