package router

import (
	"net/http"
	"nstu/internal/api/handler"
	"nstu/internal/api/middleware"
	"nstu/internal/service"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func NewRouter(h *handler.Handler, service *service.Service) *mux.Router {
	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.LoggerMiddleware())                     // 1 - логируем все
	r.Use(middleware.RateLimitMiddleware(rate.Limit(5), 10)) // 2 - проверяем лимиты
	r.Use(middleware.CORSMiddleware())                       // 3 - настраиваем CORS
	r.Use(middleware.RecoverMiddleware())                    // 4 - перехватываем панику

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.AuthMiddleware(service)) // 5 - проверяем авторизацию

	// Пользователи
	api.HandleFunc("/form", h.HandleNewForm).Methods(http.MethodPost)

	return r
}
