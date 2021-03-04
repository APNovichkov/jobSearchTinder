package main

import (
	"strconv"
	"strings"
)

func ConvertStringDateToInt(postingDate string) int{
	// Converts string like '8 days ago' to number of hours in 8 days, returns integer

	out := 0
	splitDate := strings.Split(postingDate, " ")

	if splitDate[1] == "hours" || splitDate[1] == "hour"{
		intDate, err := strconv.Atoi(splitDate[0])
		if err != nil{
			panic(err)
		}
		out = intDate
	}else if splitDate[1] == "day" || splitDate[1] == "days" {
		intDate, err := strconv.Atoi(splitDate[0])
		if err != nil{
			panic(err)
		}
		out = intDate * 24
	}

	return out
}