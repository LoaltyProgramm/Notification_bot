package reminder

import (
	"context"
	"errors"
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

func (s *ReminderService) ListReminderForID(ctx context.Context, id int) (*model.Reminder, error) {
	return s.repo.GetReminderForID(ctx, id)
}

func (s *ReminderService) RemoveReminderForID(ctx context.Context, id int) error {
	return s.repo.DeleteReminderForID(ctx, id)
}

func (s *ReminderService) CreateGroup(ctx context.Context, userSession *model.UserSession) error {
	if userSession.Group.GroupID == 0 {
		return errors.New("нету ID группы!")
	}

	if userSession.Group.UserID == 0 {
		return errors.New("нету ID пользователя!")
	}

	if userSession.Group.TitleGroup == "" {
		return errors.New("нету названия группы!")
	}

	err := s.repo.AddGroup(ctx, userSession)
	if err != nil {
		return err
	}

	return nil
}

