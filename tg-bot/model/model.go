package model

type Reminder struct {
	ChatID       int64
	Text         string
	TypeInterval string
	WeekDay      string
	Time         string
	FullTime     string
}

type UserSession struct {
	Chat_ID int64
	State         State
	UserText      string
	Interval      string
	IntervalRetry bool
	Reminder      *Reminder
}
