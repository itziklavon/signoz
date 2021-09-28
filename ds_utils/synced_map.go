package ds_utils

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"sync"
)

type ConcurrentHashMap struct {
	sync.Mutex
	Map *hashmap.Map
}

// NewSyncedMap create new map
func NewSyncedMap() *ConcurrentHashMap {
	return &ConcurrentHashMap{
		Map: hashmap.New(),
	}
}

// Put add key-value to map
func (syncedMap *ConcurrentHashMap) Put(key interface{}, value interface{}) {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	syncedMap.Map.Put(key, value)
}

// Get get value from map, if not exists - bool = false
func (syncedMap *ConcurrentHashMap) Get(key interface{}) (interface{}, bool) {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	return syncedMap.Map.Get(key)
}

// Clear remove all keys from map
func (syncedMap *ConcurrentHashMap) Clear() {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	syncedMap.Map.Clear()
}

// Empty check if map is empty
func (syncedMap *ConcurrentHashMap) Empty() bool {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	return syncedMap.Map.Empty()
}

// Size get size of map(keys in map)
func (syncedMap *ConcurrentHashMap) Size() int {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	return syncedMap.Map.Size()
}

// Keys get keys as slice
func (syncedMap *ConcurrentHashMap) Keys() []interface{} {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	return syncedMap.Map.Keys()
}

// Values get values as slice
func (syncedMap *ConcurrentHashMap) Values() []interface{} {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	return syncedMap.Map.Values()
}

// FromJson init map from json
func (syncedMap *ConcurrentHashMap) FromJson(json []byte) error {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	return syncedMap.Map.FromJSON(json)
}

// ToJson create json from map
func (syncedMap *ConcurrentHashMap) ToJson() ([]byte, error) {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	return syncedMap.Map.ToJSON()
}

// Remove remove key from map
func (syncedMap *ConcurrentHashMap) Remove(key interface{}) {
	syncedMap.Lock()
	defer syncedMap.Unlock()
	syncedMap.Map.Remove(key)
}
