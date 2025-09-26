package reminder

import (
	"context"
	"tg-app/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Addreminder(ctx context.Context, rem *model.Reminder) error
	GetReminders(ctx context.Context, userSession *model.UserSession) ([]*model.Reminder, error)
	GetReminderForID(ctx context.Context, id int) (*model.Reminder, string)
}

type PGXRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *PGXRepository {
	return &PGXRepository{
		pool: pool,
	}
}

func (r *PGXRepository) Addreminder(ctx context.Context, rem *model.Reminder) error {
	query := `
		INSERT INTO reminder (chat_id, text, week_day, type_reminder, time, full_time) VALUES ($1, $2, $3, $4, $5, $6); 
	`

	if _, err := r.pool.Exec(ctx, query, rem.ChatID, rem.Text, rem.WeekDay, rem.TypeInterval, rem.Time, rem.FullTime); err != nil {
		return err
	}

	return nil
}

func (r *PGXRepository) GetReminders(ctx context.Context, userSession *model.UserSession) ([]*model.Reminder, error) {
	query := `
		SELECT id, text, full_time FROM reminder WHERE chat_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, userSession.Chat_ID)
	if err != nil {
		return nil, err
	}

	var reminders []*model.Reminder
	for rows.Next() {
		var r model.Reminder
		if err := rows.Scan(&r.ID, &r.Text, &r.FullTime); err != nil {
			return nil, err
		}

		reminders = append(reminders, &r)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

//ошибка в данной функции? надо разораться 
func (r *PGXRepository) GetReminderForID(ctx context.Context, id int) (*model.Reminder, string) {
	query := `
		SELECT text, full_time FROM reminder WHERE id = $1;
	`

	var text string
	var fullTime string
	err := r.pool.QueryRow(context.Background(), query, id).Scan(&text, &fullTime).Error()
	if err != "" {
		return nil, err
	}

	reminder := &model.Reminder{
		Text: text,
		FullTime: fullTime,
	}

	return reminder, ""
}