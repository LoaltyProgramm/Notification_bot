package bot

import (
	"log"

	"tg-app/internal/reminder"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Bot     *tgbotapi.BotAPI
	Session *Manager
	ServiceReminder *reminder.ReminderService //вставить сервис
}

func NewHandler(bot *tgbotapi.BotAPI, session *Manager, serviceReminder *reminder.ReminderService) *Handler {
	return &Handler{
		Bot:     bot,
		Session: session,
		ServiceReminder: serviceReminder,
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

		CallbackHandlers(update.CallbackQuery.Data, update, h.Bot, userSession, chatID, h.ServiceReminder)
	}

	if UserStateFunc, ok := StateHandler[userSession.State]; ok {
		UserStateFunc(h, update, userSession, chatID)
	} else { //переделать
		msg := tgbotapi.NewMessage(chatID, "Нету данного колбека")
		if _, err := h.Bot.Send(msg); err != nil {
			log.Println("Error callback")
		}
	}
}
