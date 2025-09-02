package model

type Reminder struct {
	ChatID       int64
	Text         string
	TypeInterval string
	WeekDay      string
	Hours        int
	Minute       int
}
