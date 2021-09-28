package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOrDefault(t *testing.T) {
	val := GetOrDefault("actual", "default")
	assert.Equal(t, "actual", val)

	val = GetOrDefault("", "default")
	assert.Equal(t, "default", val)

	val = GetOrDefault(nil, "default")
	assert.Equal(t, "default", val)
}

func TestGetOrDefaultInt(t *testing.T) {
	val := GetOrDefaultInt("1", 2)
	assert.Equal(t, 1, val)

	val = GetOrDefaultInt("aaa", 2)
	assert.Equal(t, 2, val)

	val = GetOrDefaultInt("", 2)
	assert.Equal(t, 2, val)

	val = GetOrDefaultInt(nil, 2)
	assert.Equal(t, 2, val)
}
