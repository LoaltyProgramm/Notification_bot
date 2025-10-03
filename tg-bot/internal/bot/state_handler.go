package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"tg-app/internal/reminder"
	"tg-app/model"
	"tg-app/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var StateHandler = map[model.State]func(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService){
	model.StateMainMenu:          handlerMainMenu,
	model.StateRegistredText:     handlerRegistredText,
	model.StateRegistredInterval: handlerRegistredInterval,
	model.StateRegistredFinal:    handlerRegistredFinal,
	model.StateRegistredError:    handlerRegistredError,
	model.StateIdle:              handlerIdle,
	model.StateEmptyLists:     handlerEmptyLists,
	model.StateAddREminder:       handlerAddReminder,
	model.StateAllLists:          handlerAllLists,
	model.StateList:              handlerList,
}

func handlerMainMenu(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, "<b>Выберите функцию👇</b>")
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать напоминание📋", "create_reminder"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список напоминаний", "all_lists"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerRegistredText(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, "<b>Введите текст напоминания✍️</b>")
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println("ERROR - ", err)
		return
	}
	session.State = "registred_interval"
}

func handlerRegistredInterval(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	if session.IntervalRetry {
		msg := tgbotapi.NewMessage(chatID, "<b>Введите интервал напоминания⏰</b>")
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад к тексту", "redirect_registred_text"),
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
			),
		)

		if _, err := h.Bot.Send(msg); err != nil {
			log.Println(err)
			return
		}

		session.IntervalRetry = false

		session.State = "registred_final"
		return
	}

	session.UserText = update.Message.Text
	msg := tgbotapi.NewMessage(chatID, "<b>Введите интервал напоминания⏰</b>")
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад к тексту", "redirect_registred_text"),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
		),
	)

	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	session.State = "registred_final"
}

func handlerRegistredFinal(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	session.Interval = update.Message.Text

	var err error
	session.Reminder, err = utils.ParseIntervalData(chatID, session.UserText, session.Interval)
	if err != nil {
		log.Println(err)
		session.State = "registred_error"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не 1 правильный формат ввода интервала\nВведите интервал заново:")
		if _, err := h.Bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
		return
	}

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("<b>Подтверждаете напоминание?</b>\nТекст:\n%s\nИнтервал:\n%s",
			session.UserText,
			session.Interval))
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвержаю", "success_data"),
			tgbotapi.NewInlineKeyboardButtonData("Назад к интервалу", "back_interval"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerRegistredError(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	session.Interval = update.Message.Text

	var err error
	session.Reminder, err = utils.ParseIntervalData(update.Message.Chat.ID, session.UserText, session.Interval)
	if err != nil {
		session.State = "registred_error"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не 2 правильный формат ввода интервала\nВведите интервал заново:")
		if _, err := h.Bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
		return
	}

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("<b>Подтверждаете напоминание?</b>\nТекст:\n%s\nИнтервал:\n%s",
			session.UserText,
			session.Interval))
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвержаю", "success_data"),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerIdle(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	return
}

func handlerEmptyLists(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, "Список пуст\nПополняй скорее его)")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerAddReminder(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	// логика добавления в бд записи о напоминаниях
	err := service.Createreminder(context.Background(), session.Reminder)
	if err != nil {
		log.Fatal(err)
	}
	//---------------------------------------------
	session.State = "main_menu"
	msg := tgbotapi.NewMessage(chatID, "Напоминание добавлено✅")
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	handlerMainMenu(h, update, session, chatID, service)
}

func handlerAllLists(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	reminders, err := service.ListRemindersForChatID(context.Background(), session)
	if err != nil {
		log.Println(err)
		return
	}

	if len(reminders) <= 0 {
		session.State = model.StateEmptyLists
		handlerEmptyLists(h, update, session, chatID, service)
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

	session.State = "idle"
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	handlerIdle(h, update, session, chatID, service)
}

func handlerList(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Текст:\n%s\n\nИнтервал:\n%s", session.Reminder.Text, session.Reminder.FullTime))
	idStr := strconv.Itoa(session.Reminder.ID)
	log.Println(session.Reminder)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить", fmt.Sprintf("delete_reminder_%s", idStr)),//error
			tgbotapi.NewInlineKeyboardButtonData("Назад", "all_lists"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}
