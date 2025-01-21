package validator

import (
	validator "github.com/go-playground/validator/v10"
)

// Validator интерфейс для валидации
type Validator interface {
	ValidateStruct(s interface{}) error
}

// validate реализация валидатора
type validate struct {
	validator *validator.Validate
}

// New создает новый экземпляр валидатора
func New() Validator {
	return &validate{
		validator: validator.New(),
	}
}

// ValidateStruct валидирует структуру по тегам validate
func (v *validate) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}
