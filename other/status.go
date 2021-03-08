package other

import (
	"sync"
)

//Status is users and gifts message
type Status struct {
	*sync.RWMutex
	User  []userInfo
	Gifts map[string][]giftsInfo
}

type userInfo struct {
	start      int
	end        int
	err        int
	name       string
	awardsInfo map[string][]userAwards
}

type userAwards struct {
	name string
	num  int
}

type giftsInfo struct {
	name string
	ids  []int
}

//StatusStart ...
func StatusStart(t *Btoml) *Status {
	a := &Status{&sync.RWMutex{}, []userInfo{}, map[string][]giftsInfo{}}
	a.Refresh(t)
	return a
}

//Refresh the Status
func (d *Status) Refresh(t *Btoml) {
	d.Lock()
	defer d.Unlock()
Loop:
	for k, v := range t.BiliUser {
		for _, vv := range d.User {
			if v.UserName == vv.name {
				//fmt.Println("已经存在", v.UserName)
				continue Loop
			}
		}
		d.User = append(d.User, userInfo{0, 0, 0, t.BiliUser[k].UserName, map[string][]userAwards{}})

	}
}

//NewGift :if the gift is new return true, else add the giftId and return false
func (d *Status) NewGift(giftName string, giftID int) bool {
	d.Lock()
	defer d.Unlock()
	for ka, va := range d.Gifts {
		for k, v := range va {
			if v.name == giftName {
				for _, vv := range v.ids {
					if giftID == vv {
						return false
					}
				}
				if ka == TodayDay() {
					nameids := append(d.Gifts[ka][k].ids, giftID)
					giftsinfos := d.Gifts[ka]
					giftsinfos[k] = giftsInfo{giftName, nameids}
					d.Gifts[ka] = giftsinfos
					return true
				}
			}
		}
		//	giftsinfos := append(d.Gifts[ka], giftsInfo{giftName, []int{giftID}})
		//	d.Gifts[ka] = giftsinfos
		//	return true
	}
	d.Gifts[TodayDay()] = append(d.Gifts[TodayDay()], giftsInfo{giftName, []int{giftID}})
	return true
}

//IfUserOK 判断是否领取
func (d *Status) IfUserOK(k int) bool {
	d.RLock()
	defer d.RUnlock()
	if IntTime()-d.User[k].err > 3601 {
		return true
	}
	return false
}

//OkUser 领取成功
func (d *Status) OkUser(userid int, giftname string, giftnum int) {
	d.Lock()
	defer d.Unlock()
	d.User[userid].err = 0
	for ka, va := range d.User[userid].awardsInfo {
		if ka == TodayDay() {
			for k, v := range va {
				if v.name == giftname {
					newuseraward := d.User[userid].awardsInfo[ka]
					newuseraward[k].num = d.User[userid].awardsInfo[ka][k].num + giftnum
					d.User[userid].awardsInfo[ka] = newuseraward
					return
				}
			}
			newuseraward := d.User[userid].awardsInfo[ka]
			newuseraward = append(newuseraward, userAwards{giftname, giftnum})
			d.User[userid].awardsInfo[ka] = newuseraward
			return
		}
	}
	d.User[userid].awardsInfo[TodayDay()] = []userAwards{userAwards{giftname, giftnum}}
}

//ErrUser 领取失败 if return true = 刚黑屋
func (d *Status) ErrUser(k int) bool {
	d.Lock()
	defer d.Unlock()
	b := false
	if d.User[k].err < 9 /*2*/ {
		d.User[k].err++
	} else {
		b = true
		d.User[k].err = IntTime()
	}
	return b
}

//TodayGifts eg. 舰长 X 100 小电视 X 50
func (d *Status) TodayGifts() map[string]int {
	a := map[string]int{}
	d.RLock()
	defer d.RUnlock()
	for ka, va := range d.Gifts {
		if ka == TodayDay() {
			for _, v := range va {
				//fmt.Println(v.name, "X", len(v.ids))
				a[v.name] = len(v.ids)
			}
		}
	}
	return a
}

//TodayUserAward eg. 辣条 X 22 亲密度 X 33
func (d *Status) TodayUserAward() []map[string]int {
	a := []map[string]int{}
	d.RLock()
	defer d.RUnlock()
	for userid := range d.User {
		for ka, va := range d.User[userid].awardsInfo {
			if ka == TodayDay() {
				g := map[string]int{}
				for _, v := range va {
					g[v.name] = v.num
				}
				a = append(a, g)
			}
		}
	}
	return a
}

//YesterTodayGifts eg. 舰长 X 100 小电视 X 50
func (d *Status) YesterTodayGifts() map[string]int {
	a := map[string]int{}
	d.RLock()
	defer d.RUnlock()
	for ka, va := range d.Gifts {
		if ka == BeforeDay(1) {
			for _, v := range va {
				//fmt.Println(v.name, "X", len(v.ids))
				a[v.name] = len(v.ids)
			}
		}
	}
	return a
}

//YesterTodayUserAward eg. 辣条 X 22 亲密度 X 33
func (d *Status) YesterTodayUserAward() []map[string]int {
	a := []map[string]int{}
	d.RLock()
	defer d.RUnlock()
	for userid := range d.User {
		for ka, va := range d.User[userid].awardsInfo {
			if ka == BeforeDay(1) {
				g := map[string]int{}
				for _, v := range va {
					g[v.name] = v.num
				}
				a = append(a, g)
			}
		}
	}
	return a
}

//Clear2Day ...
func (d *Status) Clear2Day() {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.Gifts[BeforeDay(2)]; ok {
		delete(d.Gifts, BeforeDay(2))
		//fmt.Println("清理了2前gift")
	}
	for userid := range d.User {
		if _, ok := d.User[userid].awardsInfo[BeforeDay(2)]; ok {
			delete(d.User[userid].awardsInfo, BeforeDay(2))
			//fmt.Println("清理了2前award userid:", userid)
		}
	}
}
