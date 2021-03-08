package job

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/atwat/bigoo/other"

	"github.com/bitly/go-simplejson"
)

//JoinGifts ...
func JoinGifts(roomid, gid, gtime, gwait int, gtype, giftname string, btoml *other.Btoml, Status *other.Status, userid int) {
	randSleep(gwait, gwait+(gtime-gwait)*4/5, gid*(userid+1))
	rand.Seed(time.Now().UnixNano() - int64((userid+1)*gid))
	if rand.Intn(100)+1 > btoml.Setting.Lottery.GiftsOdds {
		return //随机放弃
	}
	if !Status.IfUserOK(userid) {
		return //黑屋中
	}
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v4/smalltv/Getaward"
	data := other.Sign(`access_key=` + btoml.BiliUser[userid].AccessToken + `&actionKey=appkey&appkey=1d8b6e7d45233436&build=` + build + `&channel=huawei&device=android&mobi_app=android&platform=android&raffleId=` + strconv.Itoa(gid) + `&roomid=` + strconv.Itoa(roomid) + `&statistics=` + statistics + `&ts=` + other.StrTime() + `&type=` + gtype)
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
		fmt.Printf("user %v joinGifts roomid %v get message err: %v\n", userid, roomid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v joinGifts roomid %v get code err: %v\n", userid, roomid, err)
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
	giftName, err := js.Get("data").Get("gift_name").String()
	if err != nil {
		fmt.Printf("user %v joinGifts roomid %v get data.gift_name err: %v\n", userid, roomid, err)
		return
	}
	giftNum, err := js.Get("data").Get("gift_num").Int()
	if err != nil {
		fmt.Printf("user %v joinGifts roomid %v get data.gift_num err: %v\n", userid, roomid, err)
		return
	}
	Status.OkUser(userid, giftName, giftNum)
	fmt.Printf("%v用户%v 在房间%v %v: %vX%v\n", other.TI(), userid, roomid, giftname, giftName, giftNum)
	if strings.Index(giftName, "瓜子") != -1 || strings.Index(giftName, "电视") != -1 {
		other.SendMail(fmt.Sprintf("用户%v获得%vX%v", btoml.BiliUser[userid].UserName, giftName, giftNum), btoml)
	}
}
