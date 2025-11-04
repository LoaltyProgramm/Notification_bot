package bot

import (
	"tg-app/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandHandler(command string, chatID int64, userSession *model.UserSession, bot tgbotapi.BotAPI) {
	switch command {
	case "start":
		if userSession.ValidUser {
			userSession.State = model.StateMainMenu
			return
		}

		msg := tgbotapi.NewMessage(chatID, "Введите пароль от бота:")

		userSession.State = "login_user"
		if _, err := bot.Send(msg); err != nil {
			return
		}
	}
}
