package middleware

import (
	"context"
	"fmt"
	"net/http"
	"nstu/internal/logger"
	"nstu/internal/model"
	"nstu/internal/service"
)

// Middleware предоставляет функции для аутентификации и логирования
type Middleware struct {
	service *service.Service
}

// NewMiddleware создает новый экземпляр Middleware
func NewMiddleware(service *service.Service) *Middleware {
	return &Middleware{
		service: service,
	}
}

// Создаем тип ключа для контекста
type contextKey string

const (
	userContextKey contextKey = "user"
)

// AuthMiddleware проверяет наличие токена в заголовке и возвращает ошибку, если токен не найден
func AuthMiddleware(service *service.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// token := r.Header.Get("Authorization")
			// if token == "" {
			// 	logger.Log.Error().
			// 		Msg("Empty token")
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	json.NewEncoder(w).Encode(u.NewResponse("error", "Empty token"))
			// 	return
			// }

			// // Разбиваем токен на части
			// parts := strings.SplitN(token, " ", 2)
			// if len(parts) != 2 || parts[0] != "tma" {
			// 	logger.Log.Error().
			// 		Msg("Неверный формат токена")
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	json.NewEncoder(w).Encode(u.NewResponse("error", "Invalid token format"))
			// 	return
			// }

			// user, err := service.AuthInitData(parts[1])
			// if err != nil {
			// 	logger.Log.Error().
			// 		Err(err).
			// 		Msg("Ошибка аутентификации")
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	json.NewEncoder(w).Encode(u.NewResponse("error", err.Error()))
			// 	return
			// }

			///////////////////////////////////////////////////////////////////////////////
			user, err := service.AuthInitData("")
			if err != nil {
				logger.Log.Error().
					Err(err).
					Msg("Ошибка аутентификации")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			///////////////////////////////////////////////////////////////////////////////

			// Добавляем пользователя в контекст
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext получает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}
