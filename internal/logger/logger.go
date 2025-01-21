package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func init() {
	// Установка московского времени
	moscow, _ := time.LoadLocation("Europe/Moscow")
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(moscow)
	}

	// Настройка вывода в JSON формате
	Log = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}

// Get возвращает настроенный логгер
func Get() *zerolog.Logger {
	return &Log
}
