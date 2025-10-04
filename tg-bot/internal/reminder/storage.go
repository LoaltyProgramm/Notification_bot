package reminder

import (
	"context"
	"tg-app/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Addreminder(ctx context.Context, rem *model.Reminder) error
	GetReminders(ctx context.Context, userSession *model.UserSession) ([]*model.Reminder, error)
	GetReminderForID(ctx context.Context, id int) (*model.Reminder, error)
	DeleteReminderForID(ctx context.Context, id int) error
	AddGroup(ctx context.Context, userSession *model.UserSession) error
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

func (r *PGXRepository) GetReminderForID(ctx context.Context, id int) (*model.Reminder, error) {
	query := `
		SELECT text, full_time FROM reminder WHERE id = $1;
	`

	var text string
	var fullTime string

	err := r.pool.QueryRow(ctx, query, id).Scan(&text, &fullTime)
	if err != nil {
		return nil, err
	}

	reminder := &model.Reminder{
		ID: id,
		Text:     text,
		FullTime: fullTime,
	}

	return reminder, nil
}

func (r *PGXRepository) DeleteReminderForID(ctx context.Context, id int) error {
	query := `
		DELETE FROM reminder WHERE id = $1;
	`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PGXRepository) AddGroup(ctx context.Context, userSession *model.UserSession) error {
	query := `
		INSERT INTO user_group (chat_id_group, user_id, title_group) VALUES ($1, $2, $3);
	`

	_, err := r.pool.Exec(ctx, query, userSession.Group.GroupID, userSession.Group.UserID, userSession.Group.TitleGroup)
	if err != nil {
		return err
	}

	return nil
}
