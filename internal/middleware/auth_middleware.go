package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"rest-api-go/internal/config"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(pool *pgxpool.Pool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accessToken")
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, errors.New("Missing access token"))
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

		fmt.Println(session.ID)

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
