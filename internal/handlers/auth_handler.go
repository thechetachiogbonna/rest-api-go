package handlers

import (
	"fmt"
	"net/http"
	"rest-api-go/internal/config"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/types"
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
	var payload types.RegisterPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := repository.CreateUser(auth.Pool, payload.FirstName, payload.LastName, payload.Email, payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	session, err := repository.CreateSession(auth.Pool, user.ID, r.RemoteAddr, r.UserAgent())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := utils.GenerateToken(utils.Payload{
		UserID:    user.ID,
		SessionID: session.ID,
	}, utils.AccessTokenRegisterClaims, config.Env.JWT_SECRET)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	refreshToken, err := utils.GenerateToken(utils.Payload{
		UserID:    user.ID,
		SessionID: session.ID,
	}, utils.RefreshTokenRegisterClaims, config.Env.JWT_REFRESH_SECRET)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.SetAuthCookies(w, accessToken, refreshToken)

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": fmt.Sprintf("User created successfully with id %q.", user.ID)})
}

func (auth *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
