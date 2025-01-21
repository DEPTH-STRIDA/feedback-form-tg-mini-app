// // internal/repository/postgres/group.go
package postgres

import (
	"fmt"
	"nstu/internal/model"

	"github.com/jmoiron/sqlx"
)

// FormRepo структура для работы с заявками
type FormRepo struct {
	db *sqlx.DB
}

// NewGroupRepo - создает новый репозиторий для работы с группами
func NewFormRepo(db *sqlx.DB) *FormRepo {
	return &FormRepo{db: db}
}

// CreateForm создает заявку
func (r *FormRepo) CreateForm(form *model.Form) error {
	query := `
		INSERT INTO forms (user_id, name, feedback, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, updated_at`

	return r.db.QueryRow(
		query,
		form.UserID,
		form.Name,
		form.Feedback,
		form.Comment,
	).Scan(&form.ID, &form.UpdatedAt)
}

// GetFormByID получает заявку по id
func (r *FormRepo) GetFormByID(id int64) (*model.Form, error) {
	form := &model.Form{}
	query := `
		SELECT id, user_id, name, feedback, comment, updated_at
		FROM forms
		WHERE id = $1`

	err := r.db.Get(form, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get form: %w", err)
	}

	return form, nil
}

// UpdateForm обновляет заявку
func (r *FormRepo) UpdateForm(form *model.Form) error {
	query := `
		UPDATE forms
		SET name = $2, feedback = $3, comment = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRow(
		query,
		form.ID,
		form.Name,
		form.Feedback,
		form.Comment,
	).Scan(&form.UpdatedAt)
}

// DeleteForm удаляет заявку
func (r *FormRepo) DeleteForm(id int64) error {
	query := `DELETE FROM forms WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete form: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("form not found")
	}

	return nil
}

// ListForms получает список заявок
func (r *FormRepo) ListForms(offset, limit int, userID int64) ([]model.Form, error) {
	forms := []model.Form{}
	query := `
		SELECT id, user_id, name, feedback, comment, updated_at
		FROM forms
		WHERE ($1 = 0 OR user_id = $1)
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3`

	err := r.db.Select(&forms, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list forms: %w", err)
	}

	return forms, nil
}
