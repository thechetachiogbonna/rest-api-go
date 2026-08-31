package utils

import (
	"net/http"
)

func GetAccessTokenCookieAttributes(accessToken string) *http.Cookie {
	return &http.Cookie{
		Name:     "accessToken",
		Path:     "/",
		Value:    accessToken,
		MaxAge:   int(FIFTEEN_MINUTES), // 15 minutes
		HttpOnly: true,
		Secure:   true,
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
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, GetAccessTokenCookieAttributes(accessToken))

	http.SetCookie(w, GetRefreshTokenCookieAttributes(refreshToken))
}
