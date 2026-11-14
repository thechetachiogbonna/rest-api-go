package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer (func() {
				fmt.Print("Running since ", time.Since(start), " ", r.Method, " ", r.URL.Path)
			})()
		})
	}
}
