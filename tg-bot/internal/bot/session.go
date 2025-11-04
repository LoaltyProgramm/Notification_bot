package bot

import (
	"sync"
	"tg-app/model"
)

type Manager struct {
	session map[int64]*model.UserSession
	mu      *sync.Mutex
}

func NewManager(session map[int64]*model.UserSession) *Manager {
	return &Manager{
		session: session,
		mu:      &sync.Mutex{},
	}
}

func (m *Manager) Get(chatID int64) *model.UserSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.session[chatID]; !ok {
		m.session[chatID] = &model.UserSession{ValidUser: false, State: model.StateMainMenu, Chat_ID: chatID}
		m.session[chatID].Group = &model.Group{}
		m.session[chatID].Reminder = &model.Reminder{}
	}

	return m.session[chatID]
}
