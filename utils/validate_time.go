package utils

import (
	"strconv"
	"time"
)

func ValidateYear(startYear, endYear string) (string, string) {
	defaultYear := strconv.Itoa(time.Now().Year())
	return ValidateTime(startYear, endYear, defaultYear)
}

func ValidateMonth(startMonth, endMonth string) (string, string) {
	defaultMonth := strconv.Itoa(int(time.Now().Month()))
	return ValidateTime(startMonth, endMonth, defaultMonth)
}

func ValidateDay(startDay, endDay string) (string, string) {
	defaultDay := strconv.Itoa(time.Now().Day())
	return ValidateTime(startDay, endDay, defaultDay)
}

func ValidateTime(start, end, defaultVal string) (string, string) {
	if start != "" && end != "" {
		return start, end
	}

	if start != "" {
		return start, defaultVal
	}

	return defaultVal, ""
}