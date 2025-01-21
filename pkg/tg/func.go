package tg

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *Bot) SendMessageUnkownChatIdD(msg tgbotapi.MessageConfig) (int64, tgbotapi.Message, error) {
	originalChatId := msg.ChatID
	var sendedMsg tgbotapi.Message
	var err error

	app.logger.Info().
		Int64("originalChatId", originalChatId).
		Str("messageText", msg.Text[:min(len(msg.Text), 50)]).
		Msg("Начинаем попытку отправки сообщения")

	// Попытка 1: отправка с оригинальным ID
	app.logger.Info().
		Int64("chatId", msg.ChatID).
		Msg("Попытка 1: Отправка с оригинальным ID")

	sendedMsg, err = app.SendMessage(msg)
	if err == nil {
		app.logger.Info().
			Int64("chatId", msg.ChatID).
			Int("messageId", sendedMsg.MessageID).
			Msg("Успешная отправка с оригинальным ID")
		return msg.ChatID, sendedMsg, nil
	}
	app.logger.Error().
		Int64("chatId", msg.ChatID).
		Err(err).
		Msg("Ошибка при отправке с оригинальным ID")

	// Попытка 2: отправка с префиксом -100
	newChatID := addNegative100Prefix(msg.ChatID)
	app.logger.Info().
		Int64("originalChatId", msg.ChatID).
		Int64("newChatId", newChatID).
		Msg("Попытка 2: Отправка с префиксом -100")

	msg.ChatID = newChatID
	sendedMsg, err = app.SendMessage(msg)
	if err == nil {
		app.logger.Info().
			Int64("chatId", msg.ChatID).
			Int("messageId", sendedMsg.MessageID).
			Msg("Успешная отправка с префиксом -100")
		return msg.ChatID, sendedMsg, nil
	}
	app.logger.Error().
		Int64("chatId", msg.ChatID).
		Err(err).
		Msg("Ошибка при отправке с префиксом -100")

	// Попытка 3: отправка с отрицательным ID
	negativeChatID := originalChatId * -1
	app.logger.Info().
		Int64("originalChatId", originalChatId).
		Int64("negativeChatId", negativeChatID).
		Msg("Попытка 3: Отправка с отрицательным ID")

	msg.ChatID = negativeChatID
	sendedMsg, err = app.SendMessage(msg)
	if err == nil {
		app.logger.Info().
			Int64("chatId", msg.ChatID).
			Int("messageId", sendedMsg.MessageID).
			Msg("Успешная отправка с отрицательным ID")
		return msg.ChatID, sendedMsg, nil
	}
	app.logger.Error().
		Int64("chatId", msg.ChatID).
		Err(err).
		Msg("Ошибка при отправке с отрицательным ID")

	// Все попытки неудачны
	app.logger.Error().
		Int64("originalChatId", originalChatId).
		Int64("lastTriedChatId", msg.ChatID).
		Err(err).
		Msg("Все попытки отправки сообщения неудачны")

	return msg.ChatID, sendedMsg, fmt.Errorf("failed all attempts to send message: %w", err)
}

func addNegative100Prefix(num int64) int64 {
	str := fmt.Sprintf("-100%d", num)
	result, _ := strconv.ParseInt(str, 10, 64)
	return result
}

// SendMessage синхронная функция для отправки сообщения
func (app *Bot) SendDeleteMessage(msg tgbotapi.DeleteMessageConfig) (*tgbotapi.APIResponse, error) {
	sendedMsg, err := app.sendDeleteMessage(msg)
	if err != nil {
		return sendedMsg, err
	}
	return sendedMsg, nil
}

// sendMessage асинхронная функция, которая с помощью waitgroup дожидается результатов от отправки сообщения
func (app *Bot) sendDeleteMessage(msg tgbotapi.DeleteMessageConfig) (*tgbotapi.APIResponse, error) {
	app.CheckAPI()

	// Устанавливаем глобальные параметры
	sendedMsg, err := app.BotAPI.Request(msg)
	if err != nil {
		return nil, err
	}

	return sendedMsg, nil
}

// SendMessageRepet делает несколько попыток отправки сообщений.
// Останавливает попытки после первой успешной.
func (app *Bot) SendMessageRepet(msg tgbotapi.MessageConfig, numberRepetion int) (tgbotapi.Message, error) {
	for i := 0; i < numberRepetion; i++ {
		sendedMsg, err := app.SendMessage(msg)
		if err != nil {
			app.logger.Info().
				Int("attempt", i).
				Err(err).
				Msg("Ошибка при отправке сообщения с повтором")
		} else {
			return sendedMsg, nil
		}
	}
	return tgbotapi.Message{}, fmt.Errorf("ни одна попытка не оказалось результативной")
}

// SendMessage синхронная функция для отправки сообщения
func (app *Bot) SendMessage(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	sendedMsg, err := app.sendMessage(msg)
	if err != nil {
		return sendedMsg, err
	}
	return sendedMsg, nil
}

// sendMessage асинхронная функция, которая с помощью waitgroup дожидается результатов от отправки сообщения
func (app *Bot) sendMessage(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	app.CheckMessage(msg.ChatID)

	sendedMsg, err := app.BotAPI.Send(msg)
	if err != nil {
		return sendedMsg, err
	}

	return sendedMsg, nil
}

// SendPinMessageEvent синхронная функция для отправки события на закрепление сообщения
func (app *Bot) SendPinMessageEvent(messageID int, ChatID int64, disableNotification bool) (*tgbotapi.APIResponse, error) {
	APIResponse, err := app.sendPinMessageEvent(messageID, ChatID, disableNotification)
	if err != nil {
		return APIResponse, err
	}
	return APIResponse, nil
}

// sendPinMessageEvent асинхронная функция, которая с помощью waitgroup дожидается результатов закрепления сообщения
// DisableNotification - если true, уведомление о закреплении не будет отправлено
func (app *Bot) sendPinMessageEvent(messageID int, ChatID int64, disableNotification bool) (*tgbotapi.APIResponse, error) {

	// Закрепление отправленного сообщения
	pinConfig := tgbotapi.PinChatMessageConfig{
		ChatID:              ChatID,
		MessageID:           messageID,
		DisableNotification: disableNotification,
	}

	app.CheckAPI()

	// Устанавливаем глобальные параметры
	APIResponse, err := app.BotAPI.Request(pinConfig)
	if err != nil {
		return nil, err
	}

	return APIResponse, nil
}

// SendSticker синхронная функция для отправки стикера
func (app *Bot) SendSticker(stickerID string, chatID int64) (*tgbotapi.Message, error) {
	sendedMsg, err := app.sendSticker(stickerID, chatID)
	if err != nil {
		return sendedMsg, err
	}
	return sendedMsg, nil
}

// sendSticker асинхронная функция, которая с помощью waitgroup дожидается результатов от отправки стикера
func (app *Bot) sendSticker(stickerID string, chatID int64) (*tgbotapi.Message, error) {

	msg := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID))

	app.CheckAPI()
	sendedMsg, err := app.BotAPI.Send(msg)
	if err != nil {
		return nil, err
	}

	return &sendedMsg, nil
}

// SendUnPinAllMessageEvent синхронная функция для отправки события на открепление всех сообщений.
func (app *Bot) SendUnPinAllMessageEvent(ChannelUsername string, chatID int64) (*tgbotapi.APIResponse, error) {
	sendedMsg, err := app.sendUnPinAllMessageEvent(ChannelUsername, chatID)
	if err != nil {
		return sendedMsg, err
	}
	return sendedMsg, nil
}

// sendUnPinAllMessageEvent асинхронная функция, которая с помощью waitgroup дожидается результатов от отправки события открепления всех сообщений
func (app *Bot) sendUnPinAllMessageEvent(ChannelUsername string, chatID int64) (*tgbotapi.APIResponse, error) {
	unpinConfig := tgbotapi.UnpinAllChatMessagesConfig{
		ChatID:          chatID,
		ChannelUsername: ChannelUsername,
	}

	app.CheckAPI()
	APIresponse, err := app.BotAPI.Request(unpinConfig)
	if err != nil {
		return nil, err
	}

	return APIresponse, err
}

func (app *Bot) EditMessageRepet(editMsg tgbotapi.EditMessageTextConfig, numberRepetion int) (*tgbotapi.APIResponse, error) {
	var err error
	var response *tgbotapi.APIResponse

	for i := 0; i < numberRepetion; i++ {
		response, err = app.editMessage(editMsg)
		if err != nil {
			app.logger.Info().
				Int("attempt", i).
				Err(err).
				Msg("Ошибка при редактировании сообщения с повтором")
		} else {
			return response, nil
		}
	}
	return nil, fmt.Errorf("ни одна попытка не стала результативной: %w", err)
}

// EditMessage синхронно редактирует сообщение
func (app *Bot) EditMessage(editMsg tgbotapi.EditMessageTextConfig) (*tgbotapi.APIResponse, error) {
	response, err := app.editMessage(editMsg)
	if err != nil {
		return response, err
	}

	return response, nil
}

// editMessage редактирует сообщение в чате, отправив функцию редактирования в запросы
func (app *Bot) editMessage(editMsg tgbotapi.EditMessageTextConfig) (*tgbotapi.APIResponse, error) {
	app.CheckAPI()
	// Устанавливаем глобальные параметры
	response, err := app.BotAPI.Request(editMsg)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (app *Bot) DeleteMessageRepet(msgToDelete tgbotapi.DeleteMessageConfig, numberRepetion int) error {
	var err error

	for i := 0; i < numberRepetion; i++ {
		err = app.deleteMessage(msgToDelete)
		if err != nil {
			app.logger.Info().
				Int("attempt", i).
				Err(err).
				Msg("Не удалось удалить сообщение из чата")
		} else {
			return nil
		}
	}

	return fmt.Errorf("ни одна попытка не стала результативной: %w", err)
}

// DeleteMessage удаляет сообщение
func (app *Bot) DeleteMessage(msgToDelete tgbotapi.DeleteMessageConfig) error {
	err := app.deleteMessage(msgToDelete)
	if err != nil {
		return err
	}

	return nil
}

func (app *Bot) deleteMessage(deleteMsg tgbotapi.DeleteMessageConfig) error {
	app.CheckAPI()
	_, err := app.BotAPI.Request(deleteMsg)
	if err != nil {
		return err
	}

	return err
}

// ShowAlert показывает пользователю предупреждение. alert по типу браузерного.
// Для закрытия такого уведомления потребуется нажать "ок"
func (app *Bot) ShowAlert(CallbackQueryID string, alertText string) {
	callback := tgbotapi.NewCallback(CallbackQueryID, alertText)
	// Это заставит текст появиться во всплывающем окне
	callback.ShowAlert = true
	app.CheckAPI()
	_, err := app.BotAPI.Request(callback)
	if err != nil {
		app.logger.Info().
			Err(err).
			Msg("Не удалось показать alert после CallbackQuery")
	}
}

func CreateKeyboard(input []string, buttonsPerRow int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	for i := 0; i < len(input); i += buttonsPerRow {
		var row []tgbotapi.KeyboardButton
		end := i + buttonsPerRow
		if end > len(input) {
			end = len(input)
		}
		for _, text := range input[i:end] {
			row = append(row, tgbotapi.NewKeyboardButton(text))
		}
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewReplyKeyboard(keyboard...)
}

type ButtonData struct {
	Text string
	Data string
}

//	buttons := [][]telegram.ButtonData{
//		{
//			{Text: "1.com", Data: "http://1.com"},
//			{Text: "2", Data: "2"},
//			{Text: "3", Data: "3"},
//		},
//		{
//			{Text: "4", Data: "4"},
//			{Text: "5", Data: "5"},
//			{Text: "6", Data: "6"},
//		},
//	}
func CreateInlineKeyboard(buttons [][]ButtonData) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, row := range buttons {
		var keyboardRow []tgbotapi.InlineKeyboardButton
		for _, btn := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Data))
		}
		keyboard = append(keyboard, keyboardRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
