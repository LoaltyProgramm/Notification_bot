package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"tg-app/model"
)

var (
	ERRORNOTNULLINTERVAL = errors.New("len nterval can't be 0")
	ERRORINCORECTTEMPLATE = errors.New("incorrect template interval")
	ERRORTYPEINTERVAL = errors.New("incorrect data interval")
)

func ParseIntervalData(userID int64, text, interval string) (*model.Reminder, error) {
	intervalArr := strings.Split(interval, " ")
	if len(interval) == 0 {
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
	switch weekDay{
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

	h, m, err := timeParse(intervalArr[3])
	if err != nil {
		return nil, err
	}

	reminder := model.Reminder{
		UserID: userID,
		Text: text,
		TypeInterval: typeInterval,
		WeekDay: weekDay,
		Hours: h,
		Minute: m,
	}

	return &reminder, nil
}

func main(){
	rem, err := ParseIntervalData(124512, "asdasdadas", "Каждый понедельник в 12:50")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rem)
}