package utils

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ERRORTIMEHOURS  = errors.New("an hour cannot be zero or greater than 23")
	ERRORTIMEMINUTE = errors.New("an minute cannot be zero or greater than 59")
	ERRORLENTIME    = errors.New("an time cannot be zero or greater than 2")
)

func timeParse(time string) (int, int, error) {
	timeArr := strings.Split(time, ":")

	if len(timeArr) <= 0 {
		return 1, 1, ERRORLENTIME
	}

	if len(timeArr) > 2 {
		return 1, 1, ERRORLENTIME
	}

	if len(timeArr) < 2 {
		return 1, 1, ERRORLENTIME
	}

	hours, err := strconv.Atoi(timeArr[0])
	if err != nil {
		return 0, 0, err
	}

	minute, err := strconv.Atoi(timeArr[1])
	if err != nil {
		return 0, 0, err
	}

	if hours < 0 || hours >= 24 {
		return 1, 1, ERRORTIMEHOURS
	}

	if minute < 0 || minute > 59 {
		return 1, 1, ERRORTIMEMINUTE
	}

	return hours, minute, nil
}
