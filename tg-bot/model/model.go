package model

type Cfg struct {
	BotAPI    string
	DBConnect string
}

type Reminder struct {
	ID           int
	ChatID       int64
	Text         string
	TypeInterval string
	WeekDay      string
	GroupID      int64
	Time         string
	FullTime     string
}

type Group struct {
	ID         int
	UserID     int64
	GroupID    int64
	TitleGroup string
}

type UserSession struct {
	Chat_ID         int64
	State           State
	UserText        string
	Interval        string
	SendGroupTitle  string
	SendGroupId     int64
	SendGroupIdint  int64
	IntervalRetry   bool
	Reminder        *Reminder
	Group           *Group
	RemoveGroup     int64
	RemoveMSG       int
	RemoveMSGChatID int64
}
