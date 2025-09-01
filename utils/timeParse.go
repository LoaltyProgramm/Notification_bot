package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

var (
	ERRORTIMEHOURS  = errors.New("an hour cannot be zero or greater than 23")
	ERRORTIMEMINUTE = errors.New("an minute cannot be zero or greater than 59")
	ERRORLENTIME    = errors.New("an time cannot be zero or greater than 2")
)

func timeParse(time string) (int ,int, error) {
	timeArr := strings.Split(time, ":")
	log.Println(timeArr)
	log.Println(len(timeArr))

	if len(timeArr) <= 0 {
		return 1, 1, ERRORLENTIME
	}

	if len(timeArr) > 2 {
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

// func main() {
// 	h, m, err := timeParse("12:30")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(h, "часов", m, "минут")

// 	fmt.Println()

// 	h1, m1, err1 := timeParse("24:30")
// 	if err1 != nil {
// 		fmt.Println(err1)
// 	}
// 	fmt.Println(h1, m1)

// 	fmt.Println()

// 	h2, m2, err2 := timeParse("23:60")
// 	if err2 != nil {
// 		fmt.Println(err2)
// 	}
// 	fmt.Println(h2, m2)
// }
