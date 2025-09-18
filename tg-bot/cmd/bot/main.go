package main

import (
	"log"
	"os"

	tg_internal "tg-app/internal/bot"
	"tg-app/internal/db"
	"tg-app/internal/reminder"
	"tg-app/model"

	telebotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	pool, err := db.InitDB("postgres://postgres:1234@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	repo := reminder.NewRepository(pool)
	reminderService := reminder.NewreminderService(repo)
	session := make(map[int64]*model.UserSession)
	manager := tg_internal.NewManager(session)
	handler := tg_internal.NewHandler(bot, manager, reminderService)

	u := telebotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handler.UpdateHandler(update)
	}
}
