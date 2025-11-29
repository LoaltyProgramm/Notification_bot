package main

import (
	"context"
	"log"
	"time"

	"tg-app/config"
	tg_internal "tg-app/internal/bot"
	"tg-app/internal/db"
	"tg-app/internal/reminder"
	"tg-app/model"

	telebotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"github.com/joho/godotenv"
)

func main() {
	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telebotapi.NewBotAPI(cfg.BotAPI)
	if err != nil {
		log.Println(err)
		return
	}

	bot.Debug = true

	pool, err := db.InitDB(cfg.DBConnect)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for range ticker.C {
			db.CheckReminderSend(context.Background(), pool, bot)
		}
	}()

	repo := reminder.NewRepository(pool)
	reminderService := reminder.NewreminderService(repo)
	session := make(map[int64]*model.UserSession)
	manager := tg_internal.NewManager(session)
	handler := tg_internal.NewHandler(bot, manager, reminderService)

	u := telebotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handler.UpdateHandler(update, cfg)
	}
}
