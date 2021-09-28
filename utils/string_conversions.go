package utils

import (
	"fmt"
	"strconv"
)

// GetOrDefault check given value(actual) is not empty, if it is empty return default value
func GetOrDefault(actual interface{}, defaultVal string) string {
	if actual == nil || len(fmt.Sprint(actual)) == 0 {
		return defaultVal
	}
	return fmt.Sprint(actual)
}

// GetOrDefaultInt check given value(actual) is not empty, if it is empty or cannot parse to int return default value
func GetOrDefaultInt(actual interface{}, defaultVal int) int {
	if actual == nil || len(fmt.Sprint(actual)) == 0 {
		return defaultVal
	}
	returnVal, err := strconv.Atoi(fmt.Sprint(actual))
	if err != nil {
		return defaultVal
	}

	return returnVal
}
