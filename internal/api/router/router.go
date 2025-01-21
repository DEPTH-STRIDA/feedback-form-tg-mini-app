package router

import (
	"nstu/internal/api/handler"
	"nstu/internal/api/middleware"
	"nstu/internal/service"

	"time"

	"github.com/gorilla/mux"
)

func NewRouter(h *handler.Handler, service *service.Service) *mux.Router {
	r := mux.NewRouter()

	m := middleware.NewMiddleware(service)

	// Создаем rate limiter: 1 запрос в секунду
	rateLimiter := middleware.NewRateLimiter(
		5,             // 1 запрос
		time.Second,   // за 1 секунду
		5*time.Minute, // очистка старых записей каждые 5 минут
	)

	// Применяем middleware в порядке:
	// 1. CORS - для всех запросов
	r.Use(m.CORSMiddleware)
	// 2. Логгер - для всех запросов
	r.Use(middleware.LoggerMiddleware)
	// 3. Rate Limiter - ограничиваем количество запросов
	r.Use(rateLimiter.RateLimitMiddleware)
	// 4. JSON заголовки
	r.Use(middleware.JSONMiddleware)
	// 5. Recover - для обработки паник
	r.Use(middleware.RecoverMiddleware)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	// 6. Аутентификация для API маршрутов
	api.Use(m.AuthMiddleware)

	// Пользователи
	api.HandleFunc("/user", h.GetUser).Methods("GET")

	return r
}
