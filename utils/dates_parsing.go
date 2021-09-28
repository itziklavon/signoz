package utils

import (
	"goapm/logger"
	"fmt"
	"runtime/debug"
	"time"
)

//go date format -
//https://programming.guide/go/format-parse-string-time-date-example.html

// layout equals to - YYYY-MM-DD HH:mm:ss
var DefaultLayout = "2006-01-02 15:04:05"

//GetTimestamp get time by given format
func GetTimestamp(date interface{}, layout string) (time.Time, error) {
	t, err := time.Parse(layout, fmt.Sprintf("%v", date))
	if err != nil {
		logger.LOGGER.Error("cannot parse string to time, error = ", err, ", stack trace = ", string(debug.Stack()))
	}
	return t, err
}

//GetTimeInMillis get time in milliseconds(seconds * 1000)
func GetTimeInMillis(timeToParse time.Time) int64 {
	return timeToParse.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
