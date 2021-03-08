package job

import (
	"fmt"
	"math/rand"
	"github.com/atwat/bigoo/other"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

var sendGiftNumList []int = []int{1314, 520, 99, 10, 1}

//ClearBag ...
func ClearBag(btoml *other.Btoml) {
	for userid := range btoml.BiliUser {
		go func(userid int) {
			for {
				if other.NowHour() == 2 || other.NowHour() == 3 {
					randSleep(0, 300, (userid+1)*23)
					go userBag(userid, btoml)
					time.Sleep(time.Hour * 17)
				} else {
					if other.NowHour() == 23 {
						randSleep(0, 60, (userid+1)*23)
						go userBag(userid, btoml)
					}
					time.Sleep(time.Minute * 33)
				}
			}
		}(userid)
	}
}
func userBag(userid int, btoml *other.Btoml) {
	//randSleep(0, 300, (userid+1)*20)
	if !btoml.MoreSetting.UseBag.Send2Wearing && !btoml.MoreSetting.UseBag.CleanExpiring2OtherMedal && !btoml.MoreSetting.KeepMedalColorful {
		return
	}
	accessKey := btoml.BiliUser[userid].AccessToken
	//medals := getUserMedals(accessKey)
	medals := getUserMedalsPc(btoml.BiliUser[userid].Cookie)
	medals = getUserMedalsPcRealRoomID(medals)
	bags := getBagList(accessKey)
	/*for k, v := range bags {
		fmt.Printf("  %v %v\n", k, v)
	}*/
	if btoml.MoreSetting.KeepMedalColorful {
	GrayLoop:
		for medalid, medal := range medals {
			//fmt.Printf("  %v %v\n", medalid, medal)
			if medal.isLighted == 0 {
				//灰色牌子
				for bagid, bag := range bags {
					if bag.id == 30607 && bag.num > 0 {
						code, msg := bagSend(accessKey, medal.targetID, medal.roomID, bag.bagID, bag.id, 1)
						if code != 0 {
							fmt.Printf("%v用户%v 从背包赠送给%v %v个%v: %v\n", other.TI(), userid, medal.targetName, 1, bag.name, msg)
							continue GrayLoop
						}
						fmt.Printf("%v用户%v 从背包赠送给%v %v个%v 成功\n", other.TI(), userid, medal.targetName, 1, bag.name)
						medals[medalid].todayFeed += 50
						bags[bagid].num--
						continue GrayLoop
					}
				}
			}
		}
	}
	if btoml.MoreSetting.UseBag.Send2Wearing {
	WearingLoop:
		for medalid, medal := range medals {
			if medal.iconCode == 1 && medal.level < 20 {
				allNum := medal.dayLimit - medal.todayFeed
				if allNum <= 0 {
					break WearingLoop
				}
				//
				for bagid, bag := range bags {
					if (bag.id == 1 || bag.id == 6 || bag.id == 30607) && bag.expire != 0 {
						price := 9999
						if bag.id == 1 {
							price = 1
						} else if bag.id == 6 {
							price = 10
						} else if bag.id == 30607 {
							price = 50
						}
						for _, sendNum := range sendGiftNumList {
							for {
								if allNum >= sendNum*price && bags[bagid].num >= sendNum {
									code, msg := bagSend(accessKey, medal.targetID, medal.roomID, bag.bagID, bag.id, sendNum)
									if code != 0 {
										fmt.Printf("%v用户%v 从背包赠送给%v %v个%v: %v\n", other.TI(), userid, medal.targetName, sendNum, bag.name, msg)
										break WearingLoop
									}
									fmt.Printf("%v用户%v 从背包赠送给%v %v个%v 成功\n", other.TI(), userid, medal.targetName, sendNum, bag.name)
									medals[medalid].todayFeed += sendNum * price
									allNum -= sendNum * price
									bags[bagid].num -= sendNum
									if allNum == 0 {
										break WearingLoop
									}
									time.Sleep(time.Second / 5)
								} else {
									time.Sleep(time.Second * 1)
									break
								}
							}
						}
					}
				}
			}
		}
	}
	if btoml.MoreSetting.UseBag.CleanExpiring2OtherMedal {
		for bagid, bag := range bags {
			if bag.expire-other.IntTime() < 86400 && (bag.id == 1 || bag.id == 6 || bag.id == 30607) && bag.expire != 0 {
			MLoop:
				for medalid, medal := range medals {
					if medal.level >= 20 {
						continue MLoop
					}
					allNum := medal.dayLimit - medal.todayFeed
					if bags[bagid].num == 0 {
						break
					}
					if allNum <= 0 {
						continue MLoop
					}
					price := 9999
					if bag.id == 1 {
						price = 1
					} else if bag.id == 6 {
						price = 10
					} else if bag.id == 30607 {
						price = 50
					}
					for _, sendNum := range sendGiftNumList {
						for {
							if allNum >= sendNum*price && bags[bagid].num >= sendNum {
								code, msg := bagSend(accessKey, medal.targetID, medal.roomID, bag.bagID, bag.id, sendNum)
								if code != 0 {
									fmt.Printf("%v用户%v 从背包赠送给%v %v个%v: %v\n", other.TI(), userid, medal.targetName, sendNum, bag.name, msg)
									continue MLoop
								}
								fmt.Printf("%v用户%v 从背包赠送给%v %v个%v 成功\n", other.TI(), userid, medal.targetName, sendNum, bag.name)
								medals[medalid].todayFeed += sendNum * price
								allNum -= sendNum * price
								bags[bagid].num -= sendNum
								if allNum == 0 {
									continue MLoop
								}
								time.Sleep(time.Second / 5)
							} else {
								time.Sleep(time.Second * 1)
								break
							}
						}
					}
				}
				for medalid, medal := range medals {
					allNum := bags[bagid].num
					if medal.iconCode == 1 {
						for _, sendNum := range sendGiftNumList {
							for {
								if allNum >= sendNum {
									code, msg := bagSend(accessKey, medal.targetID, medal.roomID, bag.bagID, bag.id, sendNum)
									if code != 0 {
										fmt.Printf("%v用户%v 从背包赠送给%v %v个%v: %v\n", other.TI(), userid, medal.targetName, sendNum, bag.name, msg)
										return
									}
									fmt.Printf("%v用户%v 从背包赠送给%v %v个%v 成功\n", other.TI(), userid, medal.targetName, sendNum, bag.name)
									medals[medalid].todayFeed += sendNum
									allNum -= sendNum
									bags[bagid].num -= sendNum
									time.Sleep(time.Second / 5)
								} else {
									time.Sleep(time.Second * 1)
									break
								}
							}
						}
					}
				}
			}
		}
	}
}

func bagSend(accessKey string, ruid, rroomid, bagID, id, num int) (int, string) {
	url := "https://api.live.bilibili.com/gift/v2/live/bag_send?" + other.Sign("access_key="+accessKey+"&actionKey=appkey&appkey=1d8b6e7d45233436&build="+build+"&channel="+channel+"&device=android&mobi_app=android&platform=android&statistics="+statistics+"&ts="+other.StrTime())
	data := fmt.Sprintf("uid=&ruid=%v&send_ruid=0&gift_id=%v&gift_num=%v&bag_id=%v&biz_id=%v&rnd=1%v&biz_code=live&data_behavior_id=&data_source_id=&jumpfrom=&version=5.54.0&click_id=&session_id=", ruid, id, num, bagID, rroomid, createCaptcha())
	b, c := other.AppReqPost(url, data)
	if !c {
		return -1, "失败"
	}
	//fmt.Println(string(b))
	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Printf("bagSend err: %v\n", err)
		return -1, "失败"
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("bagSend json code err: %v\n", err)
		return -1, "失败"
	}
	msg, err := js.Get("msg").String()
	if err != nil {
		fmt.Printf("bagSend json code err: %v\n", err)
		return -1, "失败"
	}
	return code, msg
}
func createCaptcha() string {
	return fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
}

type bagInfo struct {
	bagID  int
	name   string
	num    int
	id     int
	expire int
	mark   string
}

func getBagList(accessKey string) (userBagList []bagInfo) {
	userBagList = []bagInfo{}
	//access_key=dfae2c20e83871cdb121f8df21a5cd51&actionKey=appkey&appkey=1d8b6e7d45233436&build=5540500&channel=huawei&device=android&mobi_app=android&platform=android&room_id=7117440&statistics=%7B%22appId%22%3A1%2C%22platform%22%3A3%2C%22version%22%3A%225.54.0%22%2C%22abtest%22%3A%22%22%7D&ts=1588427969&sign=c13504e428c837f2557e454951929729
	url := "https://api.live.bilibili.com/xlive/app-room/v1/gift/bag_list?" + other.Sign("access_key="+accessKey+"&actionKey=appkey&appkey=1d8b6e7d45233436&build="+build+"&channel="+channel+"&device=android&mobi_app=android&platform=android&statistics="+statistics+"&ts="+other.StrTime())
	b, c := other.AppReqGet(url)
	if !c {
		return
	}
	//fmt.Println(string(b))
	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Printf("bagList err: %v\n", err)
		return
	}
	bagList, err := js.Get("data").Get("list").Array()
	if err != nil {
		fmt.Printf("bagList json data.list err: %v\n", err)
		return
	}
	for i := range bagList {
		bagID, err := js.Get("data").Get("list").GetIndex(i).Get("bag_id").Int()
		if err != nil {
			fmt.Printf("bagList json data.list.%v.bag_id err: %v\n", i, err)
			continue
		}
		name, err := js.Get("data").Get("list").GetIndex(i).Get("gift_name").String()
		if err != nil {
			fmt.Printf("bagList json data.list.%v.gift_name err: %v\n", i, err)
			continue
		}
		mark, err := js.Get("data").Get("list").GetIndex(i).Get("corner_mark").String()
		if err != nil {
			fmt.Printf("bagList json data.list.%v.corner_mark err: %v\n", i, err)
			continue
		}
		num, err := js.Get("data").Get("list").GetIndex(i).Get("gift_num").Int()
		if err != nil {
			fmt.Printf("bagList json data.list.%v.gift_num err: %v\n", i, err)
			continue
		}
		id, err := js.Get("data").Get("list").GetIndex(i).Get("gift_id").Int()
		if err != nil {
			fmt.Printf("bagList json data.list.%v.gift_id err: %v\n", i, err)
			continue
		}
		expire, err := js.Get("data").Get("list").GetIndex(i).Get("expire_at").Int()
		if err != nil {
			fmt.Printf("bagList json data.list.%v.expire_at err: %v\n", i, err)
			continue
		}
		userBagList = append(userBagList, bagInfo{bagID, name, num, id, expire, mark})
	}
	return
}

type medalInfo struct {
	iconCode     int //1=佩戴中 0=没
	isLighted    int //0=灰色 1=彩色
	medalName    string
	level        int
	intimacy     int //经验
	nextIntimacy int
	allIntimacy  int
	todayFeed    int
	dayLimit     int
	targetName   string
	targetID     int
	roomID       int
}

//PrintUserMedals ...
func PrintUserMedals(btoml *other.Btoml) {
	fmt.Println("用户勋章：")
	for userid := range btoml.BiliUser {
		fmt.Printf("  用户%v:\n", userid)
		medals := getUserMedalsPc(btoml.BiliUser[userid].Cookie)
		if len(medals) != 0 {
			for _, v := range medals {
				xunzhang := fmt.Sprintf("%v|%v", strLen(v.medalName, 6), v.level)
				if v.iconCode == 1 {
					xunzhang = "[" + xunzhang + "]"
				} else {
					xunzhang = " " + xunzhang + " "
				}
				if v.isLighted == 0 {
					xunzhang = xunzhang + "[X]"
				} else {
					xunzhang = xunzhang + "   "
				}
				fmt.Printf("    %v\t[%v%%]\t今日:%v/%v \n", xunzhang, v.intimacy*100/v.nextIntimacy, v.todayFeed, v.dayLimit)
			}
		} else {
			fmt.Println("    nil")
		}
	}
}

//PrintUserBags ...
func PrintUserBags(btoml *other.Btoml) {
	fmt.Println("用户背包：")
	for userid := range btoml.BiliUser {
		fmt.Printf("  用户%v:\n", userid)
		bags := getBagList(btoml.BiliUser[userid].AccessToken)
		if len(bags) != 0 {
			for _, v := range bags {
				fmt.Printf("    %v %v X %v\n", strLen("["+v.mark+"]", 8), strLen(v.name, 8), v.num)
			}
		} else {
			fmt.Println("    nil")
		}
	}
}

var allmedallevel []int = []int{0, 201, 501, 1001, 1701, 2701, 4201, 5801, 7501, 9401, 14901, 24901, 34901, 44901, 59901, 99901, 149901, 249901, 499901, 999901, 2002000, 2004500, 2007500, 2015000, 2030000, 2070000, 2160000, 2320000, 2600000, 3300000, 4500000, 6500000, 9000000, 12000000, 15500000, 19500000, 27000000, 37000000, 52000000, 102000000}

func getUserMedals(accessKey string) (userMedalList []medalInfo) {
	userMedalList = []medalInfo{}
	url := "https://api.live.bilibili.com/fans_medal/v2/HighQps/received_medals?" + other.Sign("access_key="+accessKey+"&actionKey=appkey&appkey=1d8b6e7d45233436&build="+build+"&channel="+channel+"&device=android&mobi_app=android&platform=android&statistics="+statistics+"&ts="+other.StrTime())
	b, c := other.AppReqGet(url)
	if !c {
		return
	}
	fmt.Println(string(b))
	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Printf("receivedMedals err: %v\n", err)
		return
	}
	medalList, err := js.Get("data").Get("list").Array()
	if err != nil {
		fmt.Printf("receivedMedail json data.list err: %v\n", err)
		return
	}
	for i := range medalList {
		medalName, err := js.Get("data").Get("list").GetIndex(i).Get("medal_name").String()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.medal_name err: %v\n", i, err)
			continue
		}
		level, err := js.Get("data").Get("list").GetIndex(i).Get("level").Int()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.level err: %v\n", i, err)
			continue
		}
		iconCode, err := js.Get("data").Get("list").GetIndex(i).Get("icon_code").Int() //1=佩戴中 0=没
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.icon_code err: %v\n", i, err)
			continue
		}
		roomID, err := js.Get("data").Get("list").GetIndex(i).Get("room_id").Int()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.room_id err: %v\n", i, err)
			continue
		}
		targetID, err := js.Get("data").Get("list").GetIndex(i).Get("target_id").Int()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.target_id err: %v\n", i, err)
			continue
		}
		targetName, err := js.Get("data").Get("list").GetIndex(i).Get("target_name").String()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.target_name err: %v\n", i, err)
			continue
		}
		intimacy, err := js.Get("data").Get("list").GetIndex(i).Get("intimacy").Int() //当前经验
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.intimacy err: %v\n", i, err)
			continue
		}
		nextIntimacy, err := js.Get("data").Get("list").GetIndex(i).Get("next_intimacy").Int() //总经验
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.next_intimacy err: %v\n", i, err)
			continue
		}
		allIntimacy := allmedallevel[level-1] + intimacy
		dayLimit, err := js.Get("data").Get("list").GetIndex(i).Get("day_limit").Int()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.day_limit err: %v\n", i, err)
			continue
		}
		todayFeed, err := js.Get("data").Get("list").GetIndex(i).Get("today_feed").Int()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.today_feed err: %v\n", i, err)
			continue
		}
		isLighted, err := js.Get("data").Get("list").GetIndex(i).Get("is_lighted").Int()
		if err != nil {
			fmt.Printf("receivedMedail json data.list.%v.today_feed err: %v\n", i, err)
			continue
		}
		//fmt.Printf("%v %v|%v 经验%v/%v 今日上限%v/%v %v(uid:%v)\n", iconCode, medalName, level, intimacy, nextIntimacy, todayFeed, dayLimit, targetName, targetID)
		userMedalList = append(userMedalList, medalInfo{iconCode, isLighted, medalName, level, intimacy, nextIntimacy, allIntimacy, todayFeed, dayLimit, targetName, targetID, roomID})
	}
	bubbleSort(userMedalList)
	return
}

func getUserMedalsPc(cookie string) (userMedalList []medalInfo) {
	userMedalList = []medalInfo{}
	url := "https://api.live.bilibili.com/i/api/medal?page=1&pageSize=100"
	b, c := other.PcReqGet(url, cookie)
	if !c {
		return
	}
	//fmt.Println(string(b))
	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Printf("receivedMedals err: %v\n", err)
		return
	}
	medalList, err := js.Get("data").Get("fansMedalList").Array()
	if err != nil {
		fmt.Printf("receivedMedailPc json data.fansMedalList err: %v\n", err)
		return
	}
	for i := range medalList {
		medalName, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("medal_name").String()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.medal_name err: %v\n", i, err)
			continue
		}
		level, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("level").Int()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.level err: %v\n", i, err)
			continue
		}
		iconCode, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("status").Int() //1=佩戴中 0=没
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.status err: %v\n", i, err)
			continue
		}
		roomID, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("roomid").Int()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.roomid err: %v\n", i, err)
			//continue
		}
		//realRoomID := getRealRoomID(roomID)
		targetID, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("target_id").Int()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.target_id err: %v\n", i, err)
			continue
		}
		targetName, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("target_name").String()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.target_name err: %v\n", i, err)
			continue
		}
		intimacy, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("intimacy").Int() //当前经验
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.intimacy err: %v\n", i, err)
			continue
		}
		nextIntimacy, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("next_intimacy").Int() //总经验
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.next_intimacy err: %v\n", i, err)
			continue
		}
		allIntimacy := allmedallevel[level-1] + intimacy
		dayLimit, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("day_limit").Int()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.day_limit err: %v\n", i, err)
			continue
		}
		todayFeed, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("today_feed").Int()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.today_feed err: %v\n", i, err)
			continue
		}
		isLighted, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("is_lighted").Int()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.is_lighted err: %v\n", i, err)
			continue
		}
		isReceive, err := js.Get("data").Get("fansMedalList").GetIndex(i).Get("is_receive").Int()
		if err != nil {
			fmt.Printf("receivedMedailPc json data.fansMedalList.%v.is_receive err: %v\n", i, err)
			continue
		}
		if isReceive != 1 {
			continue
		}
		//fmt.Printf("%v %v|%v 经验%v/%v 今日上限%v/%v %v(uid:%v)\n", iconCode, medalName, level, intimacy, nextIntimacy, todayFeed, dayLimit, targetName, targetID)
		userMedalList = append(userMedalList, medalInfo{iconCode, isLighted, medalName, level, intimacy, nextIntimacy, allIntimacy, todayFeed, dayLimit, targetName, targetID, roomID})
	}
	bubbleSort(userMedalList)
	return
}
func getUserMedalsPcRealRoomID(userMedalList []medalInfo) []medalInfo {
	newm := []medalInfo{}
	for k, v := range userMedalList {
		if v.roomID != 0 {
			userMedalList[k].roomID = getRealRoomID(v.roomID)
			newm = append(newm, userMedalList[k])
		}
	}
	//fmt.Println(userMedalList)
	//fmt.Println(newm)
	return newm
}
func bubbleSort(m []medalInfo) {
	var tmp medalInfo
	length := len(m)
	for i := 0; i < length; i++ {
		for j := length - 1; j > i; j-- {
			if m[j].allIntimacy > m[j-1].allIntimacy {
				tmp = m[j-1]
				m[j-1] = m[j]
				m[j] = tmp
			}
		}
	}
}
func lenStr(str string) int {
	a1 := strings.Count(str, "") - 1
	a2 := len(str)
	return (a1 + a2) / 2
}
func strLen(str string, l int) string {
	le := lenStr(str)
	if le < l {
		a := l - le
		for i := 0; i < a; i++ {
			str = str + " "
		}
	}
	return str
}
