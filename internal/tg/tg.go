package tg

import (
	"fmt"
	"nstu/internal/logger"
	"nstu/internal/model"
	"nstu/pkg/tg"
	"strings"
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
	GetMessageChats() *[]int64
}

// updateHandler обработчик, который вызывается для каждого обновления
func updateHandler() tg.HandlerFunc {
	return func(b *tg.Bot, u tgbotapi.Update) error {
		return nil
	}
}

func InitBot(config Config, newForms chan *model.Request) {

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

	go sendForm(config.GetMessageChats(), newForms)
}

func sendForm(chats *[]int64, newForms chan *model.Request) {
	for request := range newForms {
		message := formatMessage(request)
		for _, chatID := range *chats {
			msg := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: chatID,
				},
				Text: message,
			}
			msg.ParseMode = tgbotapi.ModeMarkdownV2

			_, err := Bot.SendMessage(msg)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Ошибка отправки сообщения")
			}
		}
	}
}

func formatMessage(request *model.Request) string {
	var builder strings.Builder

	// Заголовок
	builder.WriteString("📝 *Новая заявка*\n\n")

	// Информация о пользователе
	builder.WriteString(fmt.Sprintf("👤 *От:* %s", request.User.FirstName))
	if request.User.LastName != "" {
		builder.WriteString(" " + request.User.LastName)
	}
	if request.User.UserName != "" {
		builder.WriteString(fmt.Sprintf(" (@%s)", request.User.UserName))
	}
	builder.WriteString("\n\n")

	// Информация из формы
	builder.WriteString(fmt.Sprintf("📋 *Имя:* %s\n", request.Form.Name))

	// Добавляем способ обратной связи только если он указан
	if request.Form.Feedback != "" {
		builder.WriteString(fmt.Sprintf("📞 *Способ связи:* %s\n", request.Form.Feedback))
	}

	// Добавляем комментарий только если он есть
	if request.Form.Comment != "" {
		builder.WriteString(fmt.Sprintf("\n💬 *Комментарий:*\n%s\n", request.Form.Comment))
	}

	// Добавляем время создания
	builder.WriteString(fmt.Sprintf("\n🕐 *Время:* %s", request.Form.UpdatedAt.UpdatedAt.Format("02.01.2006 15:04")))

	return builder.String()
}
