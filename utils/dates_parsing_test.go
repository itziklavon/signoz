package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTimeInMillis(t *testing.T) {
	_ = GetTimeInMillis(time.Now().UTC())
}

func TestGetTimestamp(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	parsedTime, err := GetTimestamp("2021-08-08 00:00:00", layout)
	assert.Nil(t, err)
	assert.True(t, parsedTime.After(time.Now().UTC().AddDate(-40, 0, 0)))

	parsedTime, err = GetTimestamp("20210808 00:00:00", layout)
	assert.NotNil(t, err)
}
