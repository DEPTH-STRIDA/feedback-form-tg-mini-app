package service

import (
	"nstu/internal/model"
	"nstu/internal/repository"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// Servicer интерфейс для работы с бизнес логикой
type Servicer interface {
	CreateForm(userID model.User, form *model.Request) error
	GetMessageChan() chan *model.Request
}

// Service содержит бизнес-логику приложения
type Service struct {
	repo repository.Repository // репозиторий для работы с базой данных
	ch   chan *model.Request   // канал для отправки сообщений о новых заявках
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}
func (srv *Service) Auth(initData initdata.InitData) (model.User, error) {
	return model.User{}, nil
}

func (srv *Service) CreateForm(user model.User, form model.Form) error {

	srv.ch <- request
	return nil
}
func (srv *Service) GetMessageChan() chan *model.Request {
	return srv.ch
}
