package job

import (
	"fmt"
	"strconv"
	"time"

	"github.com/atwat/bigoo/other"

	"github.com/bitly/go-simplejson"
)

//EnterRoom ...
func EnterRoom(roomid int, btoml *other.Btoml, Status *other.Status, roomidbyprint bool) {
	if !roomidbyprint {
		if !ifrun(btoml) {
			return
		}
	}
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v1/lottery/getLotteryInfo"
	var access string
	if len(btoml.BiliUser) != 0 {
		access = btoml.BiliUser[0].AccessToken
	}
	data := other.Sign(`access_key=` + access + `&actionKey=appkey&appkey=1d8b6e7d45233436&build=` + build + `&channel=` + channel + `&device=android&mobi_app=android&platform=android&roomid=` + strconv.Itoa(roomid) + `&statistics=` + statistics + `&ts=` + other.StrTime())
	body, c := other.AppReqGet(url + "?" + data)
	if !c {
		return
	}
	//fmt.Println(string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Print(other.TI(), "房间", roomid, " 查看礼物时返回了非json字符串\n")
		fmt.Println(string(body))
		fmt.Println(url + "?" + data)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("enterroom %v get message err :%v\n", roomid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("enterroom %v get code err :%v\n", roomid, err)
		return
	}
	if code != 0 {
		fmt.Print(other.TI(), "房间", roomid, " 查看礼物时:", message, "\n")
		return
	}
	giftmap := map[string]int{}
	//storm
	stormID, err := js.Get("data").Get("storm").Get("id").Int64()
	if err == nil {
		stormtime, err := js.Get("data").Get("storm").Get("time").Int()
		if err == nil {
			stormnum, err := js.Get("data").Get("storm").Get("num").Int()
			if err == nil {
				if Status.NewGift("节奏风暴", int(stormID/1000000)) {
					giftmap["节奏风暴"]++
					//fmt.Println("节奏风暴", stormID, stormnum, stormtime)
					for userid := range btoml.BiliUser {
						go JoinStorm(roomid, stormID, btoml, Status, userid, stormnum, stormtime)
					}
				}
			}
		}
	}
	//redpocket
	redID, err := js.Get("data").Get("red_pocket").Get("id").Int()
	if err == nil {
		redtime, err := js.Get("data").Get("red_pocket").Get("remain_time").Int()
		if err == nil {
			if Status.NewGift("红包", redID) {
				giftmap["红包"]++
				//fmt.Println("红包", stormID, stormnum, stormtime)
				for userid := range btoml.BiliUser {
					go JoinRed(roomid, redID, btoml, Status, userid, redtime)
				}
			}
		}
	}
	//anchor
	_, err = js.Get("data").Get("anchor").Get("id").Int()
	if err == nil {
		if Status.NewGift("天选时刻", redID) {
			giftmap["天选时刻"]++
		}
	}
	//guard
	allguard, err := js.Get("data").Get("guard").Array()
	if err != nil {
		fmt.Printf("enterroom %v get data.guard err :%v\n", roomid, err)
		return
	}
	for i := range allguard {
		gid, err := js.Get("data").Get("guard").GetIndex(i).Get("id").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.guard.%v.id err :%v\n", roomid, i, err)
			continue
		}
		gtype, err := js.Get("data").Get("guard").GetIndex(i).Get("privilege_type").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.guard.%v.privilege_type err :%v\n", roomid, i, err)
			continue
		}
		gstatus, err := js.Get("data").Get("guard").GetIndex(i).Get("status").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.guard.%v.status err :%v\n", roomid, i, err)
			continue
		}
		gkeyword, err := js.Get("data").Get("guard").GetIndex(i).Get("keyword").String()
		if err != nil {
			fmt.Printf("enterroom %v get data.guard.%v.keyword err :%v\n", roomid, i, err)
			continue
		}
		gtime, err := js.Get("data").Get("guard").GetIndex(i).Get("time").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.guard.%v.keyword err :%v\n", roomid, i, err)
			continue
		}
		giftname := "上船"
		switch gtype {
		case 1:
			giftname = "总督"
		case 2:
			giftname = "提督"
		case 3:
			giftname = "舰长"
		}
		if gstatus == 1 && Status.NewGift(giftname, gid) {
			giftmap[giftname]++
			for userid := range btoml.BiliUser {
				go JoinGuard(roomid, gid, gtime, gkeyword, giftname, btoml, Status, userid)
			}
		}
	}
	//gifts
	allgifts, err := js.Get("data").Get("gift_list").Array()
	if err != nil {
		fmt.Printf("enterroom %v get data.gift_list err :%v\n", roomid, err)
		return
	}
	for i := range allgifts {
		gid, err := js.Get("data").Get("gift_list").GetIndex(i).Get("raffleId").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.gift_list.%v.raffleId err :%v\n", roomid, i, err)
			continue
		}
		giftname, err := js.Get("data").Get("gift_list").GetIndex(i).Get("title").String()
		if err != nil {
			fmt.Printf("enterroom %v get data.gift_list.%v.title err :%v\n", roomid, i, err)
			continue
		}
		gtype, err := js.Get("data").Get("gift_list").GetIndex(i).Get("type").String()
		if err != nil {
			fmt.Printf("enterroom %v get data.gift_list.%v.type err :%v\n", roomid, i, err)
			continue
		}
		gwait, err := js.Get("data").Get("gift_list").GetIndex(i).Get("time_wait").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.gift_list.%v.time_wait err :%v\n", roomid, i, err)
			continue
		}
		gtime, err := js.Get("data").Get("gift_list").GetIndex(i).Get("time").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.gift_list.%v.time err :%v\n", roomid, i, err)
			continue
		}
		if Status.NewGift(giftname, gid) {
			giftmap[giftname]++
			for userid := range btoml.BiliUser {
				go JoinGifts(roomid, gid, gtime, gwait, gtype, giftname, btoml, Status, userid)
			}
		}
	}
	//pk
	allpk, err := js.Get("data").Get("pk").Array()
	if err != nil {
		fmt.Printf("enterroom %v get data.pk err :%v\n", roomid, err)
		return
	}
	for i := range allpk {
		gid, err := js.Get("data").Get("pk").GetIndex(i).Get("id").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.pk.%v.id err :%v\n", roomid, i, err)
			continue
		}
		giftname, err := js.Get("data").Get("pk").GetIndex(i).Get("title").String()
		if err != nil {
			fmt.Printf("enterroom %v get data.pk.%v.title err :%v\n", roomid, i, err)
			continue
		}
		groomid, err := js.Get("data").Get("pk").GetIndex(i).Get("room_id").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.pk.%v.room_id err :%v\n", roomid, i, err)
			continue
		}
		gtime, err := js.Get("data").Get("pk").GetIndex(i).Get("time").Int()
		if err != nil {
			fmt.Printf("enterroom %v get data.pk.%v.time err :%v\n", roomid, i, err)
			continue
		}
		if Status.NewGift(giftname, gid) {
			giftmap[giftname]++
			for userid := range btoml.BiliUser {
				go JoinPk(groomid, gid, gtime, giftname, btoml, Status, userid)
			}
		}
	}
	//BOX
	boxTitle, err := js.Get("data").Get("activity_box").Get("title").String()
	if err == nil {
		gid, err := js.Get("data").Get("activity_box").Get("activity_id").Int()
		//fmt.Println(gid, boxTitle)
		if err == nil {
			if Status.NewGift(boxTitle, gid) {
				giftmap[boxTitle]++
				go joinBox(gid, btoml, Status)
			}
		}
	}

	if len(giftmap) == 0 {
		return
	}
	txt := ""
	for k, v := range giftmap {
		if txt != "" {
			txt = txt + "、"
		}
		if v != 1 {
			txt = fmt.Sprintf("%v%vX%v", txt, k, v)
		} else {
			txt = fmt.Sprintf("%v%v", txt, k)
		}
	}

	//winGreenPrint(fmt.Sprintf("%v房间%v 发现 %v\n", other.TI(), roomid, txt))

	fmt.Print("\033[32m", other.TI(), "房间", roomid, " 发现 ", txt, "\n\033[0m\r     \r")

}

/*
func winGreenPrint(s string) { //设置终端字体颜色
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(2|8))
	fmt.Print(s)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
}
*/
var run bool = true

func ifrun(btoml *other.Btoml) bool {
	a := btoml.Setting.Lottery.SleepTime
	//now, _ := time.LoadLocation("Asia/Shanghai")
	n := time.Now().In(time.FixedZone("CST", 28800))
	todaysecond := n.Hour()*3600 + n.Minute()*60 + n.Second()
	for _, a2 := range a {
		if a2[0]/10000*3600+(a2[0]-a2[0]/10000*10000)/100*60+a2[0]-a2[0]/100*100 <= todaysecond && todaysecond <= a2[1]/10000*3600+(a2[1]-a2[1]/10000*10000)/100*60+a2[1]-a2[1]/100*100 {
			if run {
				fmt.Printf("%v当前为休眠时间段，暂停抽奖\n", other.TI())
			}
			run = false
			return false
		}
	}
	if !run {
		fmt.Printf("%v休眠时间结束，启动抽奖\n", other.TI())
	}
	run = true
	return true
}

func getRealRoomID(rid int) int {
	url := "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + strconv.Itoa(rid)
	b, c := other.PcReqGet(url, "")
	if !c {
		return rid
	}
	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Printf("getRealRoomID json: %v\n", err)
		return rid
	}
	roomid, err := js.Get("data").Get("room_id").Int()
	if err != nil {
		fmt.Printf("getRealRoomID json data.room_id: %v\n", err)
		return rid
	}
	return roomid
}
