package tg

import (
	"nstu/pkg/tg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var states = map[string]tg.State{
	"start": Start,
	"menu":  Menu,
}

var Start = tg.State{
	Global:         true,
	Context:        true,
	AtEntranceFunc: nil,
	CatchAllFunc:   nil,
	MessageHandlers: map[string]tg.Handler{
		"/start": {
			Handle: func(b *tg.Bot, u tgbotapi.Update) error {
				b.SendMessage(tgbotapi.NewMessage(u.Message.Chat.ID, "Привет, я бот для студентов НГТУ"))
				return nil
			},
			Description: "Начало работы",
		},
		"/menu": {
			Handle: func(b *tg.Bot, u tgbotapi.Update) error {
				b.SetUserState(u.Message.Chat.ID, "menu", false, &u)
				return nil
			},
			Description: "Выполянет переход в меню",
		},
	},
	CallbackHandlers: nil,
}

var Menu = tg.State{
	Global:  false,
	Context: true,
	AtEntranceFunc: &tg.Handler{
		Handle: func(b *tg.Bot, u tgbotapi.Update) error {
			b.SendMessage(tgbotapi.NewMessage(u.Message.Chat.ID, "Что хотите сделать?\n1)Арзуб2)Выйти"))
			return nil
		},
	},
	CatchAllFunc: nil,
	MessageHandlers: map[string]tg.Handler{
		"арбуз": tg.Handler{
			Handle: func(b *tg.Bot, u tgbotapi.Update) error {
				b.SendMessage(tgbotapi.NewMessage(u.Message.Chat.ID, "АрбузАрбузАрбуз"))
				return nil
			},
		},
	},
}
