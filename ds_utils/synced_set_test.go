package ds_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSyncedHashSet(t *testing.T) {
	syncedHash := NewSyncedHashSet()
	syncedHash.Add("a")
	assert.Equal(t, 1, syncedHash.Size())
	assert.Equal(t, 1, len(syncedHash.Values()))

	syncedHash.AddAll("b", "c")
	assert.Equal(t, 3, syncedHash.Size())

	syncedHash.Remove("a")
	assert.Equal(t, 2, syncedHash.Size())

	assert.True(t, syncedHash.Contains("b"))
}
