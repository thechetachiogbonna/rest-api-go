package middleware

import (
	"errors"
	"net/http"
	"rest-api-go/internal/repository"
	"rest-api-go/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func IsAdmin(pool *pgxpool.Pool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(utils.UserIDKey).(string)
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
