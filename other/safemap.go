package other

import (
	"sync"
)

//SafeDict ...
type SafeDict struct {
	data map[string]int
	*sync.RWMutex
}

//NewSafeDict ...
func NewSafeDict(data map[string]int) *SafeDict {
	return &SafeDict{data, &sync.RWMutex{}}
}

//Len ...
func (d *SafeDict) Len() int {
	d.RLock()
	defer d.RUnlock()
	return len(d.data)
}

//Put ...
func (d *SafeDict) Put(key string, value int) (int, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.data[key]
	d.data[key] = value
	return oldValue, ok
}

//Add ...
func (d *SafeDict) Add(key string, n int) {
	old, _ := d.Get(key)
	d.Put(key, old+n)
}

//Get ...
func (d *SafeDict) Get(key string) (int, bool) {
	d.RLock()
	defer d.RUnlock()
	oldValue, ok := d.data[key]
	return oldValue, ok
}

//Delete ...
func (d *SafeDict) Delete(key string) (int, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.data[key]
	if ok {
		delete(d.data, key)
	}
	return oldValue, ok
}

//TheMap ...
func (d *SafeDict) TheMap() map[string]int {
	d.RLock()
	defer d.RUnlock()
	return d.data
}
