package ds_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSyncedList(t *testing.T) {
	lst := NewSyncedList()
	lst.Add("1")
	assert.Equal(t, 1, lst.Size())

	assert.Equal(t, 1, len(lst.Values()))

	assert.True(t, lst.Contains("1"))

	lst.AddAll("2", "3")
	assert.Equal(t, 3, lst.Size())
}
