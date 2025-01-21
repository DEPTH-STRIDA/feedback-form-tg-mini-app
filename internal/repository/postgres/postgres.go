// internal/repository/postgres/postgres.go
package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config интерфейс для конфигурации базы данных
type Config interface {
	URL() string
}

// NewPostgresDB создает новое подключение к базе данных
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return db, nil
}

// PostgresRepository репозиторий для работы с базой данных
type PostgresRepository struct {
	*UserRepo
	*FormRepo
}

// NewRepository создает новый репозиторий
func NewRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		UserRepo: NewUserRepo(db),
		FormRepo: NewFormRepo(db),
	}
}
