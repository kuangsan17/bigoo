package job

import (
	"fmt"
	"math/rand"
	"github.com/atwat/bigoo/other"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
)

//JoinStorm ...
func JoinStorm(roomid int, gid int64, btoml *other.Btoml, Status *other.Status, userid, stormnum, stormtime int) {
	rand.Seed(time.Now().UnixNano() - int64(userid+2)*gid)
	if rand.Intn(100)+1 > btoml.Setting.Lottery.StormOdds {
		return //随机放弃
	}
	if !Status.IfUserOK(userid) {
		return //黑屋中
	}
	starttime := other.IntTime()
	endtime := starttime + btoml.Setting.Lottery.StormSet[1] - 90 + stormtime
	if stormnum > 100 {
		endtime = starttime + btoml.Setting.Lottery.StormSet[1]*3/2 - 90 + stormtime
	}
	if endtime-starttime > stormtime {
		endtime = starttime + stormtime
	}
	for i := 0; other.IntTime() < endtime; i++ {
		time.Sleep(time.Millisecond * time.Duration(btoml.Setting.Lottery.StormSet[0]))
		url := "https://api.live.bilibili.com/xlive/lottery-interface/v1/storm/Join"
		data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&actionKey=appkey&appkey=1d8b6e7d45233436&build=" + build + "&channel=" + channel + "&device=android&id=" + strconv.FormatInt(gid, 10) + "&mobi_app=android&platform=android&statistics=" + statistics + "&ts=" + other.StrTime())
		body, c := other.AppReqPost(url, data)
		if !c {
			return
		}
		js, err := simplejson.NewJson(body)
		if err != nil {
			fmt.Printf("%v用户%v 在风暴%v err: not json\n", other.TI(), userid, gid/1000000)
			return
		}
		message, err := js.Get("msg").String()
		if err != nil {
			message, err = js.Get("data").Get("mobile_content").String()
			if err != nil {
				fmt.Printf("user %v JoinStorm %v get data.mobile_content err: %v\n", userid, gid/1000000, err)
				return
			}
			giftNum, err := js.Get("data").Get("gift_num").Int()
			if err != nil {
				fmt.Printf("user %v JoinStorm %v get data.gift_num err: %v\n", userid, gid/1000000, err)
				return
			}
			giftName, err := js.Get("data").Get("gift_name").String()
			if err != nil {
				fmt.Printf("user %v JoinStorm %v get data.gift_name err: %v\n", userid, gid/1000000, err)
				return
			}
			Status.OkUser(userid, giftName, giftNum)
			fmt.Printf("%v用户%v 在房间%v 节奏风暴: %v 获得%vX%v\n", other.TI(), userid, roomid, message, giftName, giftNum)
			return
		}
		code, err := js.Get("code").Int()
		if err != nil {
			fmt.Printf("user %v JoinStorm %v get code err: %v\n", userid, gid/1000000, err)
			return
		}
		if code == 400 {
			if message == "已经领取奖励" || message == "节奏风暴抽奖过期" {
				fmt.Printf("%v用户%v 在房间%v 节奏风暴: %v\n", other.TI(), userid, roomid, message)
				return
			} else if message == "你错过了奖励，下次要更快一点哦~" {
				continue
			}
		}
		_, err = js.Get("data").Array()
		if err == nil {
			fmt.Printf("%v用户%v 在房间%v 节奏风暴: %v\n", other.TI(), userid, roomid, "疑似黑屋")
			return
		}
		fmt.Printf("%v用户%v 在房间%v 节奏风暴: %v\n", other.TI(), userid, roomid, string(body))
	}
	fmt.Printf("%v用户%v 在房间%v 节奏风暴: %v\n", other.TI(), userid, roomid, "没抢到")
}
