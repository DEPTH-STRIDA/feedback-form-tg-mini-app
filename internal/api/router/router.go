package router

import (
	"nstu/internal/api/handler"
	"nstu/internal/api/middleware"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

// NewRouter создает новый маршрутизатор с использованием заданных параметров
func NewRouter(h *handler.Handler, rateLimit, burstLimit int) *mux.Router {
	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.LoggerMiddleware())                                     // 1 - логируем все
	r.Use(middleware.RateLimitMiddleware(rate.Limit(rateLimit), burstLimit)) // 2 - проверяем лимиты
	r.Use(middleware.CORSMiddleware())                                       // 3 - настраиваем CORS
	r.Use(middleware.RecoverMiddleware())                                    // 4 - перехватываем панику

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.AuthMiddleware()) // 5 - проверяем авторизацию

	// Регистрация маршрутов
	h.RegisterRoutes(api)

	return r
}
