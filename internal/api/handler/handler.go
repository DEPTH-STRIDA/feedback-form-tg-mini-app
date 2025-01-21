package handler

import (
	"encoding/json"
	"net/http"
	"nstu/internal/api/middleware"
	u "nstu/internal/api/utils"
	"nstu/internal/repository"
	"nstu/internal/service"
)

type Handler struct {
	service *service.Service
	repo    repository.Repository
}

func NewHandler(repo repository.Repository, service *service.Service) *Handler {
	return &Handler{
		repo:    repo,
		service: service,
	}
}

// HandleNewForm создает новую заявку от пользователя
func (h *Handler) HandleNewForm(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(u.NewResponse("error", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
