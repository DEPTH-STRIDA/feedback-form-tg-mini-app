package tg

import (
	"nstu/internal/logger"
	"nstu/pkg/tg"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Bot *tg.Bot
)

type Config interface {
	GetToken() string
	GetExpiration() time.Duration
	GetCleanupInterval() time.Duration
}

// updateHandler обработчик, который вызывается для каждого обновления
func updateHandler() tg.HandlerFunc {
	return func(b *tg.Bot, u tgbotapi.Update) error {
		return nil
	}
}

func InitBot(config Config) {

	bot, err := tg.NewBot(tg.Config{
		Token:           config.GetToken(),
		Expiration:      config.GetExpiration(),
		CleanupInterval: config.GetCleanupInterval(),
		States:          states,
		Logger:          &logger.Log,
		UpdateHandler:   updateHandler(),
	})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка инициализации бота")
	}
	Bot = bot
}
