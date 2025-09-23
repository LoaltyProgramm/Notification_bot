package bot

import (
	"context"
	"fmt"
	"log"
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

		userSession.State = "registred_text"
	case "back":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = "main_menu"
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

		userSession.State = "registred_interval"
	case "success_data":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}
		// логика добавления в бд записи о напоминаниях
		err := service.Createreminder(context.Background(), userSession.Reminder)
		if err != nil {
			log.Fatal(err)
		}
		//---------------------------------------------
		userSession.State = "main_menu"
		msg := tgbotapi.NewMessage(chatID, "Напоминание добавлено✅")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
	case "redirect_main_menu":
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.CallbackQuery.Message.Chat.ID,
			callback.CallbackQuery.Message.MessageID,
		)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		userSession.State = "main_menu"
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

		reminders, err := service.ListReminderForChatID(context.Background(), userSession)
		if err != nil {
			log.Println(err)
			return
		}

		if len(reminders) <= 0 {
			userSession.State = model.StateErrorInterval
		}

		lists := make([]string, 0, 10)
		for _, v := range reminders {
			lists = append(lists, fmt.Sprintf("Текст-\n%s\nИнтвервал-\n%s\n", v.Text, v.FullTime))
		}

		listsStr := strings.Join(lists, "\n")
		msg := tgbotapi.NewMessage(chatID, listsStr)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
			),
		)
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
		userSession.State = "idle"
		return
	}
}
