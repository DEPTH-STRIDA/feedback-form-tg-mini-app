package tg

import (
	"nstu/internal/config"
	"nstu/internal/logger"
	"nstu/internal/service"
	"nstu/pkg/tg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Bot *tg.Bot
)

func updateHandler(s *service.Service) tg.HandlerFunc {
	return func(b *tg.Bot, u tgbotapi.Update) error {

		if u.SentFrom() == nil {
			return nil
		}

		_, err := s.EnsureUser(u.Message.From.ID, u.Message.From.FirstName, u.Message.From.LastName, u.Message.From.UserName)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Ошибка при сохранении пользователя в БД: ")
			return err
		}

		return nil
	}
}

func InitBot(config *config.TGConfig, s *service.Service) {

	bot, err := tg.NewBot(tg.Config{
		Token:           config.Token,
		Expiration:      config.Expiration,
		CleanupInterval: config.CleanupInterval,
		States:          states,
		Logger:          &logger.Log,
		UpdateHandler:   updateHandler(s),
	})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Ошибка инициализации бота")
	}
	Bot = bot
}
