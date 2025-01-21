package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(b *Bot, u tgbotapi.Update) error

type Handler struct {
	// Handle обрабатывает входящее обновление от Telegram.
	Handle HandlerFunc

	// Description возвращает описание обработчика.
	Description string
}

// State представляет состояние бота и определяет правила обработки сообщений.
type State struct {
	// Если true, триггеры обработчиков проверяются независимо от текущего состояния пользователя
	// Если в глобальном состоянии найден подходящий обработчик, то он выполняется, а тригеры другого
	// состояния не выполняются.
	// После первого подходящего глобального состояния, другие глобальные состояния не выполняются.
	Global           bool
	Context          bool               // Если true, переход в другое состояние не происходит. Можно использовать для вызова функций при входе.
	AtEntranceFunc   *Handler           // Выполняется при входе в состояние. Не стоит использовать для глобальных состояний.
	CatchAllFunc     *Handler           // Выполняется для всех событий, которые не попали в маршруты. В глобальных состояниях следует использовать аккуратнее.
	MessageHandlers  map[string]Handler // Сопоставляет текст сообщения с обработчиком
	CallbackHandlers map[string]Handler // Сопоставляет данные callback с обработчиком
}

// NewState создает новый экземпляр State с заданными параметрами.
func NewState(global bool, context bool, atEntranceFunc *Handler, catchAllFunc *Handler) *State {
	return &State{
		Global:           global,
		Context:          context,
		AtEntranceFunc:   atEntranceFunc,
		CatchAllFunc:     catchAllFunc,
		MessageHandlers:  make(map[string]Handler),
		CallbackHandlers: make(map[string]Handler),
	}
}
