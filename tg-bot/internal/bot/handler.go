package bot

import (
	"fmt"
	"log"

	"tg-app/internal/reminder"
	"tg-app/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Bot             *tgbotapi.BotAPI
	Session         *Manager
	ServiceReminder *reminder.ReminderService
}

func NewHandler(bot *tgbotapi.BotAPI, session *Manager, serviceReminder *reminder.ReminderService) *Handler {
	return &Handler{
		Bot:             bot,
		Session:         session,
		ServiceReminder: serviceReminder,
	}
}

func (h *Handler) UpdateHandler(update tgbotapi.Update, cfg *model.Cfg) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v", r)

			var chatID int64
			if update.Message != nil {
				chatID = update.Message.Chat.ID
			} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
				chatID = update.CallbackQuery.Message.Chat.ID
			}
			userSession := h.Session.Get(chatID)

			h.handleError(update, userSession, chatID, h.ServiceReminder, err, cfg)
		}
	}()

	if update.MyChatMember != nil {
		chat := update.MyChatMember.Chat
		userID := update.MyChatMember.From.ID
		userSession := h.Session.Get(userID)
		newStatus := update.MyChatMember.NewChatMember.Status
		typeMember := chat.Type
		log.Println(userSession.State)
		if userSession.State != model.StateWaitAddGroup && typeMember != "private" {
			if newStatus != "kicked" && newStatus != "left" {
				userSession.State = model.StateErrorAddGroup
				if userHandler, ok := StateHandler[userSession.State]; ok {
					userHandler(h, update, userSession, userID, h.ServiceReminder, cfg)
				} else {
					log.Println("ERORR")
					return
				}
				return
			}
		}

		// если это группа или супергруппа
		if chat.Type == "group" || chat.Type == "supergroup" {
			switch newStatus {
			case "member":
				if userHandler, ok := StateHandler[userSession.State]; ok {
					userHandler(h, update, userSession, userID, h.ServiceReminder, cfg)
				} else {
					log.Println("ERORR")
					return
				}
				return
			case "kicked", "left":
				userSession.State = model.StateRemoveGroup
				idGroup := update.MyChatMember.Chat.ID
				userSession.RemoveGroup = idGroup
				if userHandler, ok := StateHandler[userSession.State]; ok {
					userHandler(h, update, userSession, userID, h.ServiceReminder, cfg)
				} else {
					log.Println("ERORR")
					return
				}
			}

		}
		return
	}

	if update.Message != nil && (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") {
		return
	}

	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	userSession := h.Session.Get(chatID)

	if !userSession.ValidUser {
		if userSession.State != model.StateLoginUser {
			userSession.State = model.StateReturnLogin
			if userHandler, ok := StateHandler[userSession.State]; ok {
				userHandler(h, update, userSession, chatID, h.ServiceReminder, cfg)
			} else {
				log.Println("ERORR")
				return
			}
			return
		}
	}

	// var chatID int64
	// if update.Message != nil {
	// 	chatID = update.Message.Chat.ID
	// } else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
	// 	chatID = update.CallbackQuery.Message.Chat.ID
	// }

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

	if !userSession.ValidUser {
		if userSession.State != model.StateLoginUser {
			return
		}

		if UserStateFunc, ok := StateHandler[userSession.State]; ok {
			UserStateFunc(h, update, userSession, chatID, h.ServiceReminder, cfg)
		} else { //переделать
			msg := tgbotapi.NewMessage(chatID, "Нету данного колбека")
			if _, err := h.Bot.Send(msg); err != nil {
				log.Println("Error callback")
			}
		}
	} else if userSession.ValidUser {
		if UserStateFunc, ok := StateHandler[userSession.State]; ok {
			UserStateFunc(h, update, userSession, chatID, h.ServiceReminder, cfg)
		} else { //переделать
			msg := tgbotapi.NewMessage(chatID, "Нету данного колбека")
			if _, err := h.Bot.Send(msg); err != nil {
				log.Println("Error callback")
			}
		}
	}

}

func (h *Handler) handleError(update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService, err error, cfg *model.Cfg) {
	log.Printf("Ошибка: %v", err)

	// бизнес-логика обработки ошибок
	session.State = "main_menu"

	msg := tgbotapi.NewMessage(chatID, "⚠️ Ошибка. Вас вернули в главное меню.")
	log.Println("Перед отправкой сообщения об ошибке")
	h.Bot.Send(msg)

	handlerMainMenu(h, update, session, chatID, service, cfg)
}
