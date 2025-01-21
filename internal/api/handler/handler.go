package handler

import (
	"net/http"
	"nstu/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service.Servicer
}

func NewHandler(srv service.Servicer) *Handler {
	return &Handler{
		service: srv,
	}
}

// RegisterRoutes регистрирует все маршруты приложения
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Регистрируем маршруты
	router.HandleFunc("/form", h.HandleNewForm).Methods(http.MethodPost)
}

// HandleNewForm создает новую заявку от пользователя
func (h *Handler) HandleNewForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("form processed"))
}
