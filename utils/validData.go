package main

import (
	"fmt"
	"strings"
	"errors"
	"slices"

	"tg-app/model"
)

var (
	ERRORNOTNULLINTERVAL = errors.New("Len nterval can't be 0")
	ERRORINCORECTTEMPLATE = errors.New("incorrect template interval")
)

func ParseIntervalData(interval string) (*model.Reminder, error) {
	intervalArr := strings.Split(interval, " ")
	if len(interval) == 0 {
		return nil, ERRORNOTNULLINTERVAL
	}

	if len(intervalArr) < 4 {
		return nil, ERRORINCORECTTEMPLATE
	}

	validFirstValues := []string{"каждый", "каждую", "каждое"}
	 
	if !slices.Contains(validFirstValues, intervalArr[0]) {
		return nil, ERRORINCORECTTEMPLATE
	}



	validWeekOfDay := []string{"понедельник", "вторник", "среда", "четверг", "пятница", "суббота", "воскресенье"}
	
	var typeInterval string
	switch intervalArr[1]{
	case "день":
		typeInterval = "day"
	case
	}
}

func main(){

}