package job

import (
	"fmt"
	"math/rand"
	"github.com/atwat/bigoo/other"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

//JoinGuard ...
func JoinGuard(roomid, gid, gtime int, gkeyword, giftname string, btoml *other.Btoml, Status *other.Status, userid int) {
	if gtime >= 2000 {
		randSleep(0, 1800, gid*(userid+2))
	} else {
		randSleep(0, gtime*9/10, gid*(userid+2))
	}
	rand.Seed(time.Now().UnixNano() - int64((userid+2)*gid))
	if rand.Intn(100)+1 > btoml.Setting.Lottery.GuardOdds {
		return //随机放弃
	}
	if !Status.IfUserOK(userid) {
		return //黑屋中
	}
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v3/guard/join"
	data := other.Sign(`access_key=` + btoml.BiliUser[userid].AccessToken + `&actionKey=appkey&appkey=1d8b6e7d45233436&build=` + build + `&channel=` + channel + `&device=android&id=` + strconv.Itoa(gid) + `&mobi_app=android&platform=android&roomid=` + strconv.Itoa(roomid) + `&statistics=` + statistics + `&ts=` + other.StrTime() + `&type=` + gkeyword)
	body, c := other.AppReqPost(url, data)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在房间%v %v err: not json\n", other.TI(), userid, roomid, giftname)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v JoinGuard roomid %v get message err: %v\n", userid, roomid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v JoinGuard roomid %v get code err: %v\n", userid, roomid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在房间%v %v: %v\n", other.TI(), userid, roomid, giftname, message)
		if strings.Index(message, "拒绝") != -1 || strings.Index(message, "登录") != -1 {
			if Status.ErrUser(userid) {
				fmt.Printf("%v用户%v 多次异常暂停抽奖\n", other.TI(), userid)
			}
		}
		return
	}
	giftName, err := js.Get("data").Get("award_name").String()
	if err != nil {
		fmt.Printf("user %v JoinGuard roomid %v get data.award_name err: %v\n", userid, roomid, err)
		return
	}
	giftNum, err := js.Get("data").Get("award_num").Int()
	if err != nil {
		fmt.Printf("user %v JoinGuard roomid %v get data.award_num err: %v\n", userid, roomid, err)
		return
	}
	Status.OkUser(userid, giftName, giftNum)
	fmt.Printf("%v用户%v 在房间%v %v: %vX%v\n", other.TI(), userid, roomid, giftname, giftName, giftNum)
}
