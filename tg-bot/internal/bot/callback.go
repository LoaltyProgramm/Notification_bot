package bot

import (
	"context"
	"log"
	"strconv"
	"strings"

	"tg-app/internal/reminder"
	"tg-app/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CallbackHandlers(callbackData string, callback tgbotapi.Update, bot *tgbotapi.BotAPI, userSession *model.UserSession, chatID int64, service *reminder.ReminderService) {
	switch callbackData {
	case "create_reminder":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateRegistredText
	case "back":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateMainMenu
	case "back_interval":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.IntervalRetry = true

		userSession.State = model.StateRegistredInterval
	case "success_data":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateAddREminder
	case "redirect_main_menu":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		if userSession.RemoveMSG != 0 && userSession.RemoveMSGChatID != 0 && userSession.State == model.StateRemoveGroup {
			deleteMsg = tgbotapi.NewDeleteMessage(
				userSession.RemoveMSGChatID, userSession.RemoveMSG,
			)
			if _, err := bot.Request(deleteMsg); err != nil {
				log.Println(err)
				return
			} else {
				userSession.RemoveMSG = 0
				userSession.RemoveMSGChatID = 0
			}
		}

		userSession.State = model.StateMainMenu
	case "redirect_registred_text":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateRegistredText
	case "all_lists":
		deleteMsg := tgbotapi.NewDeleteMessage(callback.CallbackQuery.Message.Chat.ID, callback.CallbackQuery.Message.MessageID)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateAllLists
	case "add_group":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateAddGroup
	case "back_add_group":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateMainMenu
	case "add":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateFinalAddGroup
	case "all_group":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateAllGroup
	}

	switch {
	case strings.HasPrefix(callbackData, "reminder_"):
		deleteMsg := tgbotapi.NewDeleteMessage(callback.CallbackQuery.Message.Chat.ID, callback.CallbackQuery.Message.MessageID)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		idStr := strings.TrimPrefix(callbackData, "reminder_")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println(err)
			return
		}

		reminder, err := service.ListReminderForID(context.Background(), id)
		if err != nil {
			log.Println(err)
			return
		}

		userSession.Reminder = reminder
		userSession.State = model.StateList
	}

	switch {
	case strings.HasPrefix(callbackData, "delete_reminder_"):
		deleteMsg := tgbotapi.NewDeleteMessage(callback.CallbackQuery.Message.Chat.ID, callback.CallbackQuery.Message.MessageID)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		idStr := strings.TrimPrefix(callbackData, "delete_reminder_")
		log.Println(idStr)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println(err)
			return
		}

		err = service.RemoveReminderForID(context.Background(), id)
		if err != nil {
			log.Println(err)
			return
		}

		userSession.State = model.StateMainMenu
		msg := tgbotapi.NewMessage(chatID, "Напоминание удалено")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
	}
}
