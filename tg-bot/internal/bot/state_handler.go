package bot

import (
	"fmt"
	"log"
	"tg-app/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var StateHandler = map[State]func(h *Handler, update tgbotapi.Update, session *UserSession, chatID int64){
	StateMainMenu:          handlerMainMenu,
	StateRegistredText:     handlerRegistredText,
	StateRegistredInterval: handlerRegistredInterval,
	StateRegistredFinal:    handlerRegistredFinal,
	StateRegistredError:    handlerRegistredError,
}

func handlerMainMenu(h *Handler, update tgbotapi.Update, session *UserSession, chatID int64) {
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

func handlerRegistredText(h *Handler, update tgbotapi.Update, session *UserSession, chatID int64) {
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

func handlerRegistredInterval(h *Handler, update tgbotapi.Update, session *UserSession, chatID int64) {
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

func handlerRegistredFinal(h *Handler, update tgbotapi.Update, session *UserSession, chatID int64) {
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

func handlerRegistredError(h *Handler, update tgbotapi.Update, session *UserSession, chatID int64) {
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
