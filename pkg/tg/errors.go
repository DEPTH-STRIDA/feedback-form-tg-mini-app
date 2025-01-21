package tg

import "fmt"

var (
	// ErrStatesNil возникает когда карта состояний nil
	ErrStatesNil = fmt.Errorf("states map is nil")

	// ErrInvalidToken возникает при пустом или невалидном токене
	ErrInvalidToken = fmt.Errorf("invalid bot token")

	// ErrNegativeExpiration возникает при отрицательном времени хранения
	ErrNegativeExpiration = fmt.Errorf("expiration time cannot be negative")

	// ErrNegativeCleanup возникает при отрицательном интервале очистки
	ErrNegativeCleanup = fmt.Errorf("cleanup interval cannot be negative")
)

// ValidationError представляет ошибку валидации с дополнительной информацией
type ValidationError struct {
	Err   error
	Value interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%v: %v", e.Err, e.Value)
}

// NewValidationError создает новую ошибку валидации
func NewValidationError(err error, value interface{}) error {
	return &ValidationError{
		Err:   err,
		Value: value,
	}
}
