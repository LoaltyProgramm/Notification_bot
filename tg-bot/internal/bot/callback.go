package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	//"fmt"
	"log"
	//"strings"
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
	case "all_lists": //отрефакторить код в файл state_handler.go
		deleteMsg := tgbotapi.NewDeleteMessage(callback.CallbackQuery.Message.Chat.ID, callback.CallbackQuery.Message.MessageID)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		reminders, err := service.ListRemindersForChatID(context.Background(), userSession)
		if err != nil {
			log.Println(err)
			return
		}
		
		if len(reminders) <= 0 {
			userSession.State = model.StateErrorInterval
			return
		}

		var rows [][]tgbotapi.InlineKeyboardButton
		for _, v := range reminders {
			btn := tgbotapi.NewInlineKeyboardButtonData(v.Text, fmt.Sprintf("reminder_%d", v.ID))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
		))

		msg := tgbotapi.NewMessage(chatID, "Ваши напоминания:")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
		userSession.State = "idle"
		return
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

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Текст:\n%s\n\nИнтервал:\n%s", reminder.Text, reminder.FullTime))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Удалить", fmt.Sprintf("delete_reminder_%s", idStr)),
				tgbotapi.NewInlineKeyboardButtonData("Назад", "all_lists"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
			),
		)
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
	}

	switch {
	case strings.HasPrefix(callbackData, "delete_reminder_"):
		deleteMsg := tgbotapi.NewDeleteMessage(callback.CallbackQuery.Message.Chat.ID, callback.CallbackQuery.Message.MessageID)
		if _, err := bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		}

		idStr := strings.TrimPrefix(callbackData, "delete_reminder_")
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
