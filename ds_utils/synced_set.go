package ds_utils

import (
	"github.com/emirpasic/gods/sets/hashset"
	"sync"
)

type ConcurrentHashSet struct {
	sync.Mutex
	Set *hashset.Set
}

// NewSyncedHashSet create new HashSet
func NewSyncedHashSet() *ConcurrentHashSet {
	return &ConcurrentHashSet{
		Set: hashset.New(),
	}
}

// Remove remove key from set
func (set *ConcurrentHashSet) Remove(item interface{}) {
	set.Lock()
	defer set.Unlock()
	set.Set.Remove(item)
}

// Add add key to set
func (set *ConcurrentHashSet) Add(item interface{}) {
	set.Lock()
	defer set.Unlock()
	set.Set.Add(item)
}

// AddALl add all keys to set
func (set *ConcurrentHashSet) AddAll(item ...interface{}) {
	set.Lock()
	defer set.Unlock()
	set.Set.Add(item...)
}

// Contains check if keys exists in set
func (set *ConcurrentHashSet) Contains(item interface{}) bool {
	set.Lock()
	defer set.Unlock()
	return set.Set.Contains(item)
}

// Values get keys as slice
func (set *ConcurrentHashSet) Values() []interface{} {
	set.Lock()
	defer set.Unlock()
	return set.Set.Values()
}

// Size get size of set
func (set *ConcurrentHashSet) Size() int {
	set.Lock()
	defer set.Unlock()
	return set.Set.Size()
}
