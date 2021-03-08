package monitor

import (
	"fmt"
	"math/rand"
	"github.com/atwat/bigoo/other"
	"time"

	"github.com/bitly/go-simplejson"
)

const (
	pageSize = 99
	maxPage  = 200
)

func lunxun(ch chan int, btoml *other.Btoml) {
	if !btoml.Monitor.Lun {
		return
	}
	for {
		roomids := lunGetRoomids()
		fmt.Println(other.TI()+"轮询发现抽奖房间数:", len(roomids))
		for roomid := range roomids {
			go func(roomid int) {
				randSleep(0, 180, roomid)
				ch <- roomid
			}(roomid)
		}
		time.Sleep(time.Second * 300)
	}
}

func lunGetRoomids() map[int]string {
	roomids := map[int]string{}
	runn := "-"
	//https://api.live.bilibili.com/room/v3/Area/getRoomList?actionKey=appkey&appkey=1d8b6e7d45233436&area_id=0&build=5540500&cate_id=0&channel=huawei&device=android&device_name=vmos&https_url_req=0&mobi_app=android&page=1&page_size=30&parent_area_id=0&platform=android&qn=0&sort_type=online&statistics={"appId":1,"platform":3,"version":"5.54.0","abtest":""}&tag_version=1&ts=1587177401&sign=74923d7a5d415c9c7504fd7e707ae104
	for i := 1; i <= maxPage; i++ {
		switch runn {
		case "-":
			runn = "\\"
		case "\\":
			runn = "|"
		case "|":
			runn = "/"
		case "/":
			runn = "-"
		}
		fmt.Print(other.TI(), "轮询直播间第", i, "页"+runn)
		url := "https://api.live.bilibili.com/room/v3/Area/getRoomList?" + other.Sign(fmt.Sprintf("actionKey=appkey&appkey=1d8b6e7d45233436&area_id=0&build=5540500&cate_id=0&channel=huawei&device=android&device_name=vmos&https_url_req=0&mobi_app=android&page=%v&page_size=%v&parent_area_id=0&platform=android&qn=0&sort_type=online&statistics={\"appId\":1,\"platform\":3,\"version\":\"5.54.0\",\"abtest\":\"\"}&tag_version=1&ts=%v", i, pageSize, other.StrTime()))
		b, c := other.AppReqGet(url)
		if !c {
			break
		}
		//fmt.Println(i, string(b))
		js, err := simplejson.NewJson(b)
		if err != nil {
			break
		}
		code, err := js.Get("code").Int()
		if err != nil {
			break
		}
		if code == 1024 {
			i = i - 1
			continue
		}
		list, err := js.Get("data").Get("list").Array()
		if err != nil {
			break
		}
		//fmt.Println(len(list))
		l := len(list)
		if l == 0 {
			break
		}
		for ii := 0; ii < l; ii++ {
			roomid, err := js.Get("data").Get("list").GetIndex(ii).Get("roomid").Int()
			if err != nil {
				fmt.Println("1", err)
				break
			}
			pendent, err := js.Get("data").Get("list").GetIndex(ii).Get("pendent_ru").String()
			if err != nil {
				fmt.Println("2", err)
				break
			}
			//fmt.Println(ii, roomid, pendent)
			if pendent != "" {
				roomids[roomid] = pendent
			}
		}
	}
	return roomids
}
func randSleep(timea, timeb int, aa int) {
	if timeb <= 1 {
		return
	}
	rand.Seed(time.Now().UnixNano() - int64(aa))
	if timea < 0 {
		a := rand.Intn(timeb*1000 - 800)
		time.Sleep(time.Duration(int64(a) * 1000000))
	}
	b := timeb - timea
	c := rand.Intn(b*1000-1000) + timea*1000 + 500
	time.Sleep(time.Duration(int64(c) * 1000000))
}
