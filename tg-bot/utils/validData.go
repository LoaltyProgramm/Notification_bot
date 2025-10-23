package utils

import (
	"errors"
	"slices"
	"strings"

	"tg-app/model"
)

var (
	ERRORNOTNULLINTERVAL  = errors.New("len nterval can't be 0")
	ERRORINCORECTTEMPLATE = errors.New("incorrect template interval")
	ERRORTYPEINTERVAL     = errors.New("incorrect data interval")
)

func ParseIntervalData(session *model.UserSession) (*model.Reminder, error) {
	intervalArr := strings.Split(session.Interval, " ")
	if len(intervalArr) == 0 {
		return nil, ERRORNOTNULLINTERVAL
	}

	if len(intervalArr) < 4 {
		return nil, ERRORINCORECTTEMPLATE
	}

	firstValue := strings.ToLower(intervalArr[0])

	validFirstValues := []string{"каждый", "каждую", "каждое"}

	if !slices.Contains(validFirstValues, firstValue) {
		return nil, ERRORINCORECTTEMPLATE
	}

	//правка данных для правильной валидации
	weekDay := strings.ToLower(intervalArr[1])
	switch weekDay {
	case "среду":
		weekDay = "среда"
	case "пятницу":
		weekDay = "пятница"
	case "субботу":
		weekDay = "суббота"
	}

	validWeekOfDay := []string{"понедельник", "вторник", "среда", "четверг", "пятница", "суббота", "воскресенье"}

	var typeInterval string
	if weekDay == "день" {
		typeInterval = "day"
	} else if slices.Contains(validWeekOfDay, weekDay) {
		typeInterval = "week"
	} else {
		return nil, ERRORTYPEINTERVAL
	}

	_, _, err := timeParse(intervalArr[3])
	if err != nil {
		return nil, err
	}

	fullTime := strings.Join(intervalArr, " ")

	reminder := model.Reminder{
		ChatID:       session.Chat_ID,
		Text:         session.UserText,
		TypeInterval: typeInterval,
		WeekDay:      weekDay,
		Time:         intervalArr[3],
		FullTime:     fullTime,
	}

	return &reminder, nil
}
