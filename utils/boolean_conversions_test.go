package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBoolean(t *testing.T) {
	assert.True(t, ParseBoolean("True"))
	assert.False(t, ParseBoolean("False"))
	assert.False(t, ParseBoolean(""))
	assert.False(t, ParseBoolean(nil))
}
