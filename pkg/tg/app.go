package tg

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

// Config структура для конфигурации бота
type Config struct {
	Token           string           // Токен бота
	Expiration      time.Duration    // Время хранения состояний пользователя
	CleanupInterval time.Duration    // Интервал очистки кеша
	States          map[string]State // Состояния пользователя
	Logger          *zerolog.Logger  // Логгер для записи событий
	UpdateHandler   HandlerFunc      // Обработчик, который будет вызываться при получении любого обновления
}

// Bot структура для бота
type Bot struct {
	BotAPI        *tgbotapi.BotAPI // API бота. Экспортируется для доступа к нему из вне
	expiration    time.Duration    // Время хранения состояний пользователя
	limiter       *Limiter         // Лимитер для ограничения количества запросов к API
	cache         *gocache.Cache   // Кеш для хранения состояний пользователей
	logger        *zerolog.Logger  // Логгер для записи событий
	states        map[string]State // Состояния пользователя
	globalStates  []*State         // Состояния, в которые может перейти пользователь из любого другоо
	updateHandler HandlerFunc      // Обработчик, который будет вызываться при получении любого обновления
}

// Конструктор нового бота
func NewBot(config Config) (*Bot, error) {
	if config.States == nil {
		return nil, ErrStatesNil
	}
	if config.Expiration < 0 {
		return nil, NewValidationError(ErrNegativeExpiration, config.Expiration)
	}
	if config.CleanupInterval < 0 {
		return nil, NewValidationError(ErrNegativeCleanup, config.CleanupInterval)
	}
	if config.Token == "" {
		return nil, ErrInvalidToken
	}

	botAPI, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("не удается инициализировать бота telegram: %v", err)
	}

	globalStates := make([]*State, 0)
	for _, state := range config.States {
		if state.Global {
			globalStates = append(globalStates, &state)
		}
	}

	app := Bot{
		BotAPI:        botAPI,
		limiter:       NewLimiter(),
		cache:         gocache.New(config.Expiration, config.CleanupInterval),
		states:        config.States,
		globalStates:  globalStates,
		expiration:    config.Expiration,
		logger:        config.Logger,
		updateHandler: config.UpdateHandler,
	}

	go app.HandleUpdates()

	return &app, nil
}

// HandleUpdates запускает обработку всех обновлений поступающих боту из телеграмма
func (app *Bot) HandleUpdates() {
	// Настройка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := app.BotAPI.GetUpdatesChan(u)
	app.logger.Info().Msg("Запуск обработки обновлений")
	for update := range updates {

		go func() {
			if app.updateHandler != nil {
				app.updateHandler(app, update)
			}
		}()

		go func(update tgbotapi.Update) {

			// Обработка локальных стейтов
			if update.SentFrom() == nil {
				return
			}

			// Обработка глобальных стейтов
			globalStateFound, err := app.HandleGlobalStates(update)
			if err != nil {
				app.logger.Error().Err(err).Msg("failed to handle global state")
			}
			// Если глобальное состояние найдено, то выходим из функции
			if globalStateFound {
				return
			}
			// Получение названия состояния пользователя
			userStateName, err := app.GetUserState(update.SentFrom().ID)
			if err != nil {
				app.logger.Error().Err(err).Msg("failed to get user state")
			}
			// Получени состояния
			userState, ok := app.states[userStateName]
			if !ok {
				app.logger.Debug().Str("state", userStateName).Msg("state not found in states map")
			}

			// Выбор обработчика состояния
			_, err = app.SelectHandler(update, &userState)
			if err != nil {
				app.logger.Error().Err(err).Msg("failed to handle user state")
			}

		}(update)

	}
}

func (app *Bot) GetUserState(userId int64) (string, error) {
	userStateInterface, ok := app.cache.Get(strconv.FormatInt(userId, 10))
	if !ok {
		return "", fmt.Errorf("пользовательский статус не найден") // Обработка ошибки
	}

	userState, ok := userStateInterface.(string)
	if !ok {
		return "", fmt.Errorf("пользовательский статус не найден") // Обработка ошибки
	}

	return userState, nil // Возврат состояния пользователя
}

// SetUserState меняет состояние пользователя
// immediate - если true, то новое состояние применится сразу к текущему сообщению
func (app *Bot) SetUserState(userId int64, state string, immediate bool, update *tgbotapi.Update) {
	key := strconv.FormatInt(userId, 10)

	st, ok := app.states[state]
	if !ok {
		app.logger.Error().Str("state", state).Msg("state not found")
		return
	}
	if !st.Context {
		return
	}

	app.cache.Set(key, state, app.expiration)

	if newState, ok := app.states[state]; ok {
		// Вызываем действие при входе, если оно есть и это не глобальное состояние
		if newState.AtEntranceFunc != nil && !newState.Global && update != nil {
			if err := newState.AtEntranceFunc.Handle(app, *update); err != nil {
				app.logger.Error().
					Err(err).
					Str("state", state).
					Msg("failed to handle entrance function")
			}
		}

		// Если нужна немедленная реакция
		if immediate && update != nil {
			_, err := app.SelectHandler(*update, &newState)
			if err != nil {
				app.logger.Error().
					Err(err).
					Str("state", state).
					Msg("failed to handle immediate reaction")
			}
		}
	}
}

// HandleGlobalStates проверяет подходит ли действие пользователя под
// глобальные состояния и если подходит, то выполняет его.
// Возвращает true, если глобальное состояние найдено и false, если не найдено.
func (app *Bot) HandleGlobalStates(update tgbotapi.Update) (bool, error) {
	// Обработка всех глобальных состояний
	for _, state := range app.globalStates {
		// Обработка состояния
		handlerIsFound, err := app.SelectHandler(update, state)
		// Если ошибка, то пропускаем состояние
		if err != nil {
			app.logger.Error().Err(err).Msg("failed to handle global state")
			continue
		}
		// Если обработчик найден, то возвращаем true
		if handlerIsFound {
			return true, nil
		}
	}
	return false, nil
}

func (app *Bot) SelectHandler(update tgbotapi.Update, userState *State) (bool, error) {
	switch {
	case update.Message != nil:
		if userState.MessageHandlers != nil {
			return app.handleMessage(userState, update)
		} else {
			app.logger.Info().
				Str("command", update.Message.Text).
				Int64("chat_id", update.Message.Chat.ID).
				Str("username", update.Message.Chat.UserName).
				Msg("command not found")
			return false, nil
		}
	case update.CallbackQuery != nil:
		if userState.CallbackHandlers != nil {
			return app.handleCallback(userState, update)
		} else {
			app.logger.Info().
				Str("callback", update.CallbackQuery.Data).
				Int64("user_id", update.CallbackQuery.From.ID).
				Str("username", update.CallbackQuery.From.UserName).
				Msg("callback not found")
			return false, nil
		}
	}
	return false, nil
}

// handleMessage ищет команду в map'е и выполняет ее
func (app *Bot) handleMessage(userState *State, update tgbotapi.Update) (bool, error) {
	messageFound := false

	if currentAction, ok := userState.MessageHandlers[strings.ToLower(strings.TrimSpace(update.Message.Text))]; ok {
		messageFound = true
		if err := currentAction.Handle(app, update); err != nil {
			app.logger.Error().
				Err(err).
				Str("command", update.Message.Text).
				Int64("chat_id", update.Message.Chat.ID).
				Str("username", update.Message.Chat.UserName).
				Msg("failed to handle command")
		} else {
			app.logger.Info().
				Str("command", update.Message.Text).
				Int64("chat_id", update.Message.Chat.ID).
				Str("username", update.Message.Chat.UserName).
				Msg("command handled successfully")
		}
	} else {
		if userState.CatchAllFunc != nil {
			err := userState.CatchAllFunc.Handle(app, update)
			if err != nil {
				app.logger.Error().
					Err(err).
					Int64("chat_id", update.Message.Chat.ID).
					Str("username", update.Message.Chat.UserName).
					Str("command", update.Message.Text).
					Msg("failed to handle command")
			}
		} else {
			app.logger.Info().
				Int64("chat_id", update.Message.Chat.ID).
				Str("username", update.Message.Chat.UserName).
				Str("command", update.Message.Text).
				Msg("command not found")
		}

	}
	return messageFound, nil
}

// handleCallback ищет команду в map'е и выполняет ее
func (app *Bot) handleCallback(userState *State, update tgbotapi.Update) (bool, error) {
	callbackFound := false

	if currentAction, ok := userState.CallbackHandlers[update.CallbackQuery.Data]; ok {
		callbackFound = true
		if err := currentAction.Handle(app, update); err != nil {
			app.logger.Error().
				Err(err).
				Str("callback", update.CallbackQuery.Data).
				Int64("user_id", update.CallbackQuery.From.ID).
				Str("username", update.CallbackQuery.From.UserName).
				Msg("failed to handle callback")
			return callbackFound, err
		}

		app.logger.Info().
			Str("callback", update.CallbackQuery.Data).
			Int64("user_id", update.CallbackQuery.From.ID).
			Str("username", update.CallbackQuery.From.UserName).
			Msg("callback handled successfully")
	} else {
		if userState.CatchAllFunc != nil {
			err := userState.CatchAllFunc.Handle(app, update)
			if err != nil {
				app.logger.Error().
					Err(err).
					Int64("user_id", update.CallbackQuery.From.ID).
					Str("username", update.CallbackQuery.From.UserName).
					Str("callback", update.CallbackQuery.Data).
					Msg("failed to handle callback")
				return callbackFound, err
			}
		} else {
			app.logger.Info().
				Int64("user_id", update.CallbackQuery.From.ID).
				Str("username", update.CallbackQuery.From.UserName).
				Str("callback", update.CallbackQuery.Data).
				Msg("callback not found")
		}
	}
	return callbackFound, nil
}

// CheckMessage проверяет возможность отправки сообщения в чат
func (b *Bot) CheckMessage(chatID int64) {
	b.limiter.mu.Lock()
	defer b.limiter.mu.Unlock()

	now := time.Now()

	if lastTime, ok := b.limiter.ChatTimes[chatID]; ok {
		if diff := time.Second - now.Sub(lastTime); diff > 0 {
			time.Sleep(diff)
			now = time.Now()
		}
	}

	b.cleanup(now)
	if len(b.limiter.MessageTimes) >= MultiChatLimit || len(b.limiter.ApiTimes) >= GlobalLimit {
		time.Sleep(WaitTime)
		now = time.Now()
		b.cleanup(now)
	}

	b.limiter.ChatTimes[chatID] = now
	b.limiter.MessageTimes = append(b.limiter.MessageTimes, now)
	b.limiter.ApiTimes = append(b.limiter.ApiTimes, now)
}

// CheckAPI проверяет возможность отправки запроса к API
func (b *Bot) CheckAPI() {
	b.limiter.mu.Lock()
	defer b.limiter.mu.Unlock()

	now := time.Now()
	b.cleanup(now)

	if len(b.limiter.ApiTimes) >= GlobalLimit {
		time.Sleep(WaitTime)
		now = time.Now()
		b.cleanup(now)
	}

	b.limiter.ApiTimes = append(b.limiter.ApiTimes, now)
}

// cleanup удаляет записи старше 1 секунды
func (b *Bot) cleanup(now time.Time) {
	newMessageTimes := b.limiter.MessageTimes[:0]
	for _, t := range b.limiter.MessageTimes {
		if now.Sub(t) < time.Second {
			newMessageTimes = append(newMessageTimes, t)
		}
	}
	b.limiter.MessageTimes = newMessageTimes

	newAPITimes := b.limiter.ApiTimes[:0]
	for _, t := range b.limiter.ApiTimes {
		if now.Sub(t) < time.Second {
			newAPITimes = append(newAPITimes, t)
		}
	}
	b.limiter.ApiTimes = newAPITimes
}
