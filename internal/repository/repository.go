// Package repository содержит интерфейсы для работы с базой данных
package repository

import (
	"pet1/internal/model"
)

type Repository interface {
	User
	Form
}

type User interface {
	CreateUser(user *model.User) error
	CreateUserIfNotExists(user *model.User) error
	UpdateUser(user *model.User) error
	GetUserByID(id int64) (*model.User, error)
}

type Form interface {
	CreateForm(form *model.Form) error
	GetFormByID(id int64) (*model.Form, error)
	UpdateForm(form *model.Form) error
	DeleteForm(id int64) error
	ListForms(offset, limit int, userID int64) ([]model.Form, error)
}
