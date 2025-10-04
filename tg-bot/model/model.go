package model

type Reminder struct {
	ID           int
	ChatID       int64
	Text         string
	TypeInterval string
	WeekDay      string
	Time         string
	FullTime     string
}

type Group struct {
	UserID     int64
	GroupID    int64
	TitleGroup string
}

type UserSession struct {
	Chat_ID       int64
	State         State
	UserText      string
	Interval      string
	IntervalRetry bool
	Reminder      *Reminder
	Group         *Group
}
