package handlers

import (
	"errors"
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

	utils.WriteJson(w, http.StatusCreated, map[string]any{"user": user})
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

	utils.WriteJson(w, http.StatusOK, map[string]any{"user": user})
}

func (auth *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("accessToken")

	if err != nil && cookie.Value != "" {
		claims, err := utils.VerifyToken(cookie.Value, config.Env.JWT_SECRET)
		if err == nil {
			err = repository.DeleteSession(auth.Pool, claims.UserID, claims.SessionID)
		}
	}

	utils.ClearAuthCookies(w)

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (auth *Auth) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")

	if err != nil || cookie.Value == "" {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	claims, err := utils.VerifyToken(cookie.Value, config.Env.JWT_REFRESH_SECRET)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Invalid refresh token"))
		return
	}

	_, err = repository.GetSessionByUserIDAndSessionID(auth.Pool, claims.UserID, claims.SessionID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("Invalid refresh token"))
		return
	}

	accessToken, refreshToken, err := utils.GetAccessAndRefreshTokens(claims.UserID, claims.SessionID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SetAuthCookies(w, accessToken, refreshToken)
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Access token refreshed successfully"})
}
