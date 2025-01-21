// internal/repository/postgres/postgres.go
package postgres

import (
	"fmt"
	"nstu/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.URL())
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return db, nil
}

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
