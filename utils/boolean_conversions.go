package utils

import (
	"goapm/logger"
	"fmt"
	"strconv"
)

//parse boolean val, if parsing is not valid, return false by default
func ParseBoolean(valToParse interface{}) bool {
	if valToParse == nil {
		return false
	}

	valToPrseStr := fmt.Sprint(valToParse)
	parsedBool, err := strconv.ParseBool(valToPrseStr)
	if err != nil {
		logger.LOGGER.Info("parsed string is not boolean returning false", valToPrseStr)
		return false
	}
	return parsedBool
}
