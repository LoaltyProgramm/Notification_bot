package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandHandler(command string, chatID int64, userSession *UserSession, bot tgbotapi.BotAPI) {
	switch command {
	case "start":
		msg := tgbotapi.NewMessage(chatID, "Привет👋\nДанный бот позволяет добавить напоминания к группе")

		userSession.State = "main_menu"
		if _, err := bot.Send(msg); err != nil {
			return
		}
	}
}
