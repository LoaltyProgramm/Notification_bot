package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Bot     *tgbotapi.BotAPI
	Session *Manager
}

func NewHandler(bot *tgbotapi.BotAPI, session *Manager) *Handler {
	return &Handler{
		Bot:     bot,
		Session: session,
	}
}

func (h *Handler) UpdateHandler(update tgbotapi.Update) {
	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	userSession := h.Session.Get(chatID)

	//обработка команд
	if update.Message != nil && update.Message.IsCommand() {
		CommandHandler(update.Message.Command(), chatID, userSession, *h.Bot)
	}

	//обработка коллбеков
	if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := h.Bot.Request(callback); err != nil {
			log.Println(err)
			return
		}

		CallbackHandlers(update.CallbackQuery.Data, update, h.Bot, userSession, chatID)
	}

	if UserStateFunc, ok := StateHandler[userSession.State]; ok {
		UserStateFunc(h, update, userSession, chatID)
	}
}
