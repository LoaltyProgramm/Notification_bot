package db

import (
	"context"
	"log"
	"tg-app/model"
	"time"

	telebotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckReminderSend(ctx context.Context, db *pgxpool.Pool, bot *telebotapi.BotAPI) {
	now := time.Now()
	currentTime := now.Format("15:04")
	currentWeekday := now.Weekday()

	query := `
		SELECT text, group_send_id 
		FROM reminder 
		WHERE time = $1 
		AND (type_reminder = 'day' OR (type_reminder = 'week' AND week_day = $2));
	`
	rows, err := db.Query(ctx, query, currentTime, currentWeekday)
	if err != nil {
		log.Println(err)
	}

	var response []*model.ResponseReminderSendGroup

	for rows.Next() {
		var r model.ResponseReminderSendGroup
		if err := rows.Scan(&r.Text, &r.GroupSendId); err != nil {
			log.Println(err)
		}

		response = append(response, &r)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	if len(response) != 0 {
		for _, r := range response {
			msg := telebotapi.NewMessage(r.GroupSendId, r.Text)
			bot.Send(msg)
			log.Printf("Отправлено сообщение %s в группу %d", r.Text, r.GroupSendId)
		}
	}
}
