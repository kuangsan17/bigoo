package monitor

import (
	"github.com/atwat/bigoo/other"
	"sync"
	"time"
)

type safeid struct {
	data map[int]int
	*sync.RWMutex
}

func newSafeid(data map[int]int) *safeid {
	return &safeid{data, &sync.RWMutex{}}
}
func (d *safeid) get(key int) (int, bool) {
	d.RLock()
	defer d.RUnlock()
	old, ok := d.data[key]
	return old, ok
}
func (d *safeid) put(key int, value int) (int, bool) {
	d.Lock()
	defer d.Unlock()
	old, ok := d.data[key]
	d.data[key] = value
	return old, ok
}
func (d *safeid) del(key int) (int, bool) {
	d.Lock()
	defer d.Unlock()
	old, ok := d.data[key]
	if ok {
		delete(d.data, key)
	}
	return old, ok
}

//Monitor is all monitor
func Monitor(mainch chan int, btoml *other.Btoml) {
	ch1 := make(chan int)
	//ch1 := make(chan int)
	roomids := newSafeid(map[int]int{})
	go bilive(ch1, mainch, btoml)
	go local(ch1, btoml)
	go yj(ch1, mainch, btoml)
	go lunxun(ch1, btoml)
	for {
		roomid := <-ch1
		go func(ch chan int, roomid int) {
			k, _ := roomids.get(roomid)
			if k == 0 {
				roomids.put(roomid, 1)
			}
			roomids.put(roomid, k+1)
			k, _ = roomids.get(roomid)
			time.Sleep(time.Second * 5)
			k2, _ := roomids.get(roomid)
			if k == k2 {
				ch <- roomid
				roomids.del(roomid)
			}

		}(mainch, roomid)
	}
}
