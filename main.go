package main

import (
	"fmt"
	"log"
	"os"

	telebotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var(
	State string
	UserText string
	Interval string
)

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

		if update.CallbackQuery != nil {
			if update.CallbackQuery.Message != nil {
				callback := telebotapi.NewCallback(update.CallbackQuery.ID, "")
				bot.Send(callback)

				switch update.CallbackQuery.Data {
				case "create_reminder":
					State = "registred_text"
				case "back":
					State = "main_menu"
				case "success_data":
					msg := telebotapi.NewMessage(chatID, "Напоминание добавлено✅")
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
						continue
					}
				}

			}
		}

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := telebotapi.NewMessage(chatID, "Привет👋\nДанный бот позволяет добавить напоминания к группе")
				State = "main_menu"
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
			}
		}

		switch State {
		case "main_menu":
			log.Println("PING")
			msg := telebotapi.NewMessage(chatID, "*Выберите функцию*👇")
			msg.ParseMode = "MarkDownV2"
			msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Создать напоминание📋", "create_reminder"),
				),
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Помощь🆘", "help"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}
		case "registred_text":
			msg := telebotapi.NewMessage(chatID, "*Введите текст напоминания✍️*")
			msg.ParseMode = "MarkDownV2"
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}
			State = "registred_interval"
		case "registred_interval":
			UserText = update.Message.Text
			msg := telebotapi.NewMessage(chatID, "*Введите интервал напоминания⏰*")
			msg.ParseMode = "MarkDownV2"
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}
			State = "registred_final"
		case "registred_final":
			Interval = update.Message.Text
			msg := telebotapi.NewMessage(chatID, fmt.Sprintf("*Вы подтверждаете добавление данного напоминания*?\nТекст:\n%s\nИнтервал:\n%s", UserText, Interval))
			msg.ParseMode = "MarkDownV2"
			msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Подтвержаю", "success_data"),
					telebotapi.NewInlineKeyboardButtonData("Вернуться назад", "back"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
