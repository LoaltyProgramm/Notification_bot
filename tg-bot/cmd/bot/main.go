package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"tg-app/model"
	"tg-app/utils"

	telebotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserSession struct {
	State         string
	UserText      string
	Interval      string
	IntervalRetry bool
	Remin         *model.Reminder
}

var count int

var dbReminder = make(map[int]*model.Reminder)

var rem *model.Reminder

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
				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				session.State = "main_menu"
			case "back_interval":
				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				session.IntervalRetry = true

				session.State = "registred_interval"
			case "success_data":
				deleteMsg := telebotapi.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}
				// логика добавления в бд записи о напоминаниях
				count += 1
				dbReminder[count] = session.Remin

				session.State = "main_menu"
				msg := telebotapi.NewMessage(chatID, "Напоминание добавлено✅")
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
			case "all_lists":
				deleteMsg := telebotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println(err)
					continue
				}

				if len(dbReminder) <= 0 {
					msg := telebotapi.NewMessage(chatID, "Список пуст\nПополняй скорее его)")
					msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
						telebotapi.NewInlineKeyboardRow(
							telebotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
						),
					)
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
						continue
					}
					continue
				}

				lists := make([]string, 0, 10)
				for _, v := range dbReminder {
					lists = append(lists, fmt.Sprintf("Текст-\n%s\nИнтвервал-\nКаждый %s в %d:%d\n\n", v.Text, v.WeekDay, v.Hours, v.Minute))
				}

				listsStr := strings.Join(lists, "\n")
				msg := telebotapi.NewMessage(chatID, listsStr)
				msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
					telebotapi.NewInlineKeyboardRow(
						telebotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
					),
				)
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
				continue
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
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Список напоминаний", "all_lists"),
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
			if session.IntervalRetry {
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

				session.IntervalRetry = false

				session.State = "registred_final"
				continue
			}

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

			session.State = "registred_final"

		case "registred_final":
			session.Interval = update.Message.Text

			session.Remin, err = utils.ParseIntervalData(update.Message.Chat.ID, session.UserText, session.Interval)
			if err != nil {
				log.Println(err)
				session.State = "registred_error"

				msg := telebotapi.NewMessage(update.Message.Chat.ID, "Не 1 правильный формат ввода интервала\nВведите интервал заново:")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
				continue
			}

			msg := telebotapi.NewMessage(chatID,
				fmt.Sprintf("<b>Подтверждаете напоминание?</b>\nТекст:\n%s\nИнтервал:\n%s",
					session.UserText,
					session.Interval))
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = telebotapi.NewInlineKeyboardMarkup(
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Подтвержаю", "success_data"),
					telebotapi.NewInlineKeyboardButtonData("Назад к интервалу", "back_interval"),
				),
				telebotapi.NewInlineKeyboardRow(
					telebotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
				continue
			}
		case "registred_error":
			session.Interval = update.Message.Text
			session.Remin, err = utils.ParseIntervalData(update.Message.Chat.ID, session.UserText, session.Interval)
			if err != nil {
				session.State = "registred_error"
				msg := telebotapi.NewMessage(update.Message.Chat.ID, "Не 2 правильный формат ввода интервала\nВведите интервал заново:")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
				continue
			}

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
