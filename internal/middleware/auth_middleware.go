package middleware

import (
	"context"
	"errors"
	"net/http"
	"rest-api-go/internal/config"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(pool *pgxpool.Pool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("accessToken")
			if err != nil || cookie.Value == "" {
				utils.WriteError(w, http.StatusUnauthorized, errors.New("Missing access token"))
				return
			}

			claims, err := utils.VerifyToken(cookie.Value, config.Env.JWT_SECRET)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			session, err := repository.GetSessionByUserIDAndSessionID(pool, claims.UserID, claims.SessionID)
			if err != nil || session.ExpiresAt.Before(time.Now()) {
				utils.WriteError(w, http.StatusUnauthorized, errors.New("Invalid or expired session"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func IsAdmin(pool *pgxpool.Pool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok || userID == "" {
				utils.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
				return
			}

			user, err := repository.GetUserByID(pool, userID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			if !utils.IsUserRoleAdmin(user.Role) {
				utils.WriteError(w, http.StatusForbidden, errors.New("Forbidden"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
