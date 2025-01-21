package middleware

import (
	"net/http"
	"nstu/internal/logger"
)

// RecoverMiddleware обрабатывает паники в HTTP обработчиках
func RecoverMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Log.Error().Interface("error", err).Msg("Паника в HTTP обработчике")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
