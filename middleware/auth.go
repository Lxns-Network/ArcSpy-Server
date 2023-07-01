package middleware

import (
	"github.com/gookit/config/v2"
	"net/http"
)

func AuthMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			if authToken != config.SubDataMap("api").Str("token") {
				RespondWithJSON(w, 401, "unauthorized", nil)
				return
			}

			next(w, r)
		}
	}
}
