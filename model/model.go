package model

type Reminder struct {
	UserID       int64
	Text         string
	TypeInterval string
	WeekDay      string
	Hours        int
	Minute       int
}
