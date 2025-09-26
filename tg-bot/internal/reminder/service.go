package reminder

import (
	"context"
	"tg-app/model"
)

type ReminderService struct {
	repo Repository
}

func NewreminderService(repo Repository) *ReminderService {
	return &ReminderService{
		repo: repo,
	}
}

func (s *ReminderService) Createreminder(ctx context.Context, rem *model.Reminder) error {
	return s.repo.Addreminder(ctx, rem)
}

func (s *ReminderService) ListRemindersForChatID(ctx context.Context, userSession *model.UserSession) ([]*model.Reminder, error) {
	return s.repo.GetReminders(ctx, userSession)
}

func (s *ReminderService) ListReminderForID(ctx context.Context, id int) (*model.Reminder, string) {
	return s.repo.GetReminderForID(ctx, id)
}

