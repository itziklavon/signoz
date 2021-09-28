package ds_utils

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"sync"
)

// ConcurrentList thread safe list, more like java list and not go slice
type ConcurrentList struct {
	sync.Mutex
	List *arraylist.List
}

// NewSyncedList creates new array list
func NewSyncedList() *ConcurrentList {
	return &ConcurrentList{
		List: arraylist.New(),
	}
}

// Add adds new item to list
func (list *ConcurrentList) Add(item interface{}) {
	list.Lock()
	defer list.Unlock()
	list.List.Add(item)
}

// AddAll adds multiple items to list
func (list *ConcurrentList) AddAll(item ...interface{}) {
	list.Lock()
	defer list.Unlock()
	list.List.Add(item...)
}

// Contains check if item received is already exists in list
func (list *ConcurrentList) Contains(item interface{}) bool {
	list.Lock()
	defer list.Unlock()
	return list.List.Contains(item)
}

// Values get list as slice
func (list *ConcurrentList) Values() []interface{} {
	list.Lock()
	defer list.Unlock()
	return list.List.Values()
}

// Size get list size
func (list *ConcurrentList) Size() int {
	list.Lock()
	defer list.Unlock()
	return list.List.Size()
}
