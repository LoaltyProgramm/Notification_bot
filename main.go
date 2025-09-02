package main

import (
	"fmt"
	"log"
	"os"
	"tg-app/model"
	"tg-app/utils"

	telebotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserSession struct {
	State    string
	UserText string
	Interval string
}

var dbReminder = make(map[int64]*model.Reminder)

var sessions = make(map[int64]*UserSession)

func getSession(chatID int64) *UserSession {
	if _, ok := sessions[chatID]; !ok {
		sessions[chatID] = &UserSession{State: "main_menu"}
	}
	return sessions[chatID]
}

func main() {
	token := os.Getenv("TOKEN_BOT")
	if token == "" {
		log.Fatal("No token")
	}

	bot, err := telebotapi.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return
	}

	bot.Debug = true

	u := telebotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var chatID int64
		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		}

		session := getSession(chatID)

		// команды
		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := telebotapi.NewMessage(chatID, "Привет👋\nДанный бот позволяет добавить напоминания к группе")

				session.State = "main_menu"
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
			}
		}

		// обработка коллбеков
		if update.CallbackQuery != nil {
			callback := telebotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println(err)
				continue
			}
			bot.Send(callback)

			switch update.CallbackQuery.Data {
			case "create_reminder":
				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				session.State = "registred_text"
			case "back":
				delete(dbReminder, update.CallbackQuery.From.ID)

				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				session.State = "main_menu"
			case "success_data":
				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				session.State = "main_menu"
				msg := telebotapi.NewMessage(chatID, dbReminder[update.CallbackQuery.From.ID].Text) //"Напоминание добавлено✅"
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
			case "redirect_main_menu":
				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				session.State = "main_menu"
			case "redirect_registred_text":
				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				session.State = "registred_text"
			}
		}

		// логика состояний
		switch session.State {
		case "main_menu":
			msg := telebotapi.NewMessage(chatID, "<b>Выберите функцию👇</b>")
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Создать напоминание📋", "create_reminder"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}

		case "registred_text":
			msg := telebotapi.NewMessage(chatID, "<b>Введите текст напоминания✍️</b>")
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Назад", "back"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Println("ERROR - ", err)
				continue
			}
			session.State = "registred_interval"

		case "registred_interval":

			session.UserText = update.Message.Text
			msg := telebotapi.NewMessage(chatID, "<b>Введите интервал напоминания⏰</b>")
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Назад к тексту", "redirect_registred_text"),
					telebotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
				),
			)

			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}

			session.State = "registred_pre_final"

		case "registred_pre_final":
			session.Interval = update.Message.Text

			rem, err := utils.ParseIntervalData(update.Message.Chat.ID, session.UserText, session.Interval)
			if err != nil {
				log.Println(err)
				session.State = "registred_interval"

				msg := telebotapi.NewMessage(update.Message.Chat.ID, "Не правильный формат ввода интервала\nВведите интервал заново:")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
				continue
			} 

			log.Println(session.State)
			dbReminder[update.Message.From.ID] = rem
			session.State = "registred_final"
			log.Println(session.State)
		case "registred_final":
			msg := telebotapi.NewMessage(chatID,
				fmt.Sprintf("<b>Подтверждаете напоминание?</b>\nТекст:\n%s\nИнтервал:\n%s",
					session.UserText,
					session.Interval))
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Подтвержаю", "success_data"),
					telebotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}
		}
	}

}
