package bot

import (
	"sync"
	"tg-app/model"
)

// type State string

// const (
// 	StateMainMenu          State = "main_menu"
// 	StateRegistred         State = "registred_text"
// 	StateRegistredInterval State = "registred_interval"
// 	StateRegistredFinal    State = "registred_final"
// 	StateRegistredError    State = "registred_error"
// )

type UserSession struct {
	State         State
	UserText      string
	Interval      string
	IntervalRetry bool
	Reminder      *model.Reminder
}

type Manager struct {
	session map[int64]*UserSession
	mu      sync.Mutex
}

func (m *Manager) Get(chatID int64) *UserSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.session[chatID]; !ok {
		m.session[chatID] = &UserSession{State: StateMainMenu}
	}

	return m.session[chatID]
}
