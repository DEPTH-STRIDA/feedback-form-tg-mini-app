package postgres

import (
	"fmt"
	"nstu/internal/model"

	"github.com/jmoiron/sqlx"
)

// UserRepo структура для работы с пользователями
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo - создает новый репозиторий для работы с пользователями
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CreateUser создает пользователя
func (r *UserRepo) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, username)
		VALUES ($1, $2, $3, $4)
		RETURNING updated_at`

	return r.db.QueryRow(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.UserName,
	).Scan(&user.UpdatedAt)
}

// CreateUserIfNotExists создает пользователя если не существует
func (r *UserRepo) CreateUserIfNotExists(user *model.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, username)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET first_name = $2, last_name = $3, username = $4, updated_at = CURRENT_TIMESTAMP
		RETURNING updated_at`

	return r.db.QueryRow(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.UserName,
	).Scan(&user.UpdatedAt)
}

func (r *UserRepo) UpdateUser(user *model.User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, username = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	result := r.db.QueryRow(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.UserName,
	)

	return result.Scan(&user.UpdatedAt)
}

func (r *UserRepo) GetUserByID(id int64) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, first_name, last_name, username, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.Get(user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
