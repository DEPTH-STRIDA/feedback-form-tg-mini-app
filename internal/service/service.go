package service

import (
	"context"
	"nstu/internal/model"
	"nstu/internal/repository"
)

// Servicer интерфейс для работы с бизнес логикой
type Servicer interface {
	CreateForm(ctx context.Context, userID int64, form *model.Form) error
	GetForm(ctx context.Context, id int64) (*model.Form, error)
}

// Service содержит бизнес-логику приложения
type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (srv *Service) CreateForm(ctx context.Context, userID int64, form *model.Form) error { return nil }
func (srv *Service) GetForm(ctx context.Context, id int64) (*model.Form, error)           { return nil, nil }
