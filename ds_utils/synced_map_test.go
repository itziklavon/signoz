package ds_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSyncedMap(t *testing.T) {
	syncedMap := NewSyncedMap()
	syncedMap.Put("a", "a")
	_, ok := syncedMap.Get("a")
	assert.True(t, ok)
	assert.False(t, syncedMap.Empty())
	assert.Equal(t, 1, syncedMap.Size())
	assert.Equal(t, 1, len(syncedMap.Values()))
	assert.Equal(t, 1, len(syncedMap.Keys()))

	jsonMap, _ := syncedMap.ToJson()
	_ = syncedMap.FromJson(jsonMap)

	syncedMap.Remove("a")

	assert.Equal(t, 0, syncedMap.Size())
	syncedMap.Put("a", "a")
	assert.Equal(t, 1, syncedMap.Size())
	syncedMap.Clear()
	assert.Equal(t, 0, syncedMap.Size())
}
