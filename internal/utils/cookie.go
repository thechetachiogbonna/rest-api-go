package utils

import (
	"net/http"
	"rest-api-go/internal/config"
)

var secureCookie = config.Env.GO_ENV == "production"

func GetAccessTokenCookieAttributes(accessToken string) *http.Cookie {
	return &http.Cookie{
		Name:     "accessToken",
		Path:     "/",
		Value:    accessToken,
		MaxAge:   int(FIFTEEN_MINUTES), // 15 minutes
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
	}
}

func GetRefreshTokenCookieAttributes(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     "refreshToken",
		Path:     REFRESH_PATH,
		Value:    refreshToken,
		MaxAge:   int(THIRTY_DAYS), // 30 days
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
	}
}

func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, GetAccessTokenCookieAttributes(accessToken))

	http.SetCookie(w, GetRefreshTokenCookieAttributes(refreshToken))
}

func ClearAuthCookies(w http.ResponseWriter) {
	accessToken := &http.Cookie{
		Name:     "accessToken",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
	}

	refreshToken := &http.Cookie{
		Name:     "refreshToken",
		Path:     REFRESH_PATH,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, accessToken)
	http.SetCookie(w, refreshToken)
}
