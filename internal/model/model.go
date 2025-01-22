// Package model содержит модели данных
package model

import (
	"time"
)

// BaseModel - базовая модель для всех моделей
type BaseModel struct {
	ID
	UpdatedAt
}

type ID struct {
	ID int64 `json:"-" db:"id" sql:"primary key,autoincrement"` // ID пользователя из телеграмма и БД
}

type UpdatedAt struct {
	UpdatedAt time.Time `json:"-" db:"updated_at" sql:"not null,default:current_timestamp"` // Дата последнего обновления
}

// User пользователь
type User struct {
	ID int64 `json:"-" db:"id" sql:"primary key"` // ID пользователя из телеграмма
	UpdatedAt
	FirstName string `json:"firstName" db:"first_name" sql:"not null,type:varchar(64)" validate:"required,max=64"` // Имя из Telegram, 1-64 символа
	LastName  string `json:"-" db:"last_name" sql:"type:varchar(64)" validate:"max=64"`                            // Фамилия из Telegram, 0-64 символа
	UserName  string `json:"-" db:"username" sql:"type:varchar(32),unique" validate:"max=32"`                      // Username из Telegram, 5-32 символа
}

// Form заявка оставленная пользователем
type Form struct {
	BaseModel
	UserID   int64  `json:"-" db:"user_id" sql:"not null,references:users(id),index"`                    // id пользователя
	Name     string `json:"name" db:"name" sql:"not null,type:varchar(128)" validate:"required,max=128"` // Имя пользователя
	Feedback string `json:"feedback" db:"feedback" sql:"type:varchar(256)" validate:"max=256"`           // Предпочтительный способ обратной связи
	Comment  string `json:"comment" db:"comment" sql:"type:varchar(512)" validate:"max=512"`             // Комментарий к заявке
}

type Request struct {
	Form Form
	User User
}
