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

//JoinPk ...
func JoinPk(roomid, gid, gtime int, giftname string, btoml *other.Btoml, Status *other.Status, userid int) {
	randSleep(0, gtime*10/11, gid*(userid+2))
	rand.Seed(time.Now().UnixNano() - int64((userid+2)*gid))
	if rand.Intn(100)+1 > btoml.Setting.Lottery.PkOdds {
		return //随机放弃
	}
	if !Status.IfUserOK(userid) {
		return //黑屋中
	}
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v1/pk/join"
	data := other.Sign(`access_key=` + btoml.BiliUser[userid].AccessToken + `&actionKey=appkey&appkey=1d8b6e7d45233436&build=` + build + `&channel=` + channel + `&device=android&id=` + strconv.Itoa(gid) + `&mobi_app=android&platform=android&roomid=` + strconv.Itoa(roomid) + `&statistics=` + statistics + `&ts=` + other.StrTime())
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
		fmt.Printf("user %v JoinPk roomid %v get message err: %v\n", userid, roomid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v JoinPk roomid %v get code err: %v\n", userid, roomid, err)
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
	gift, err := js.Get("data").Get("award_text").String()
	if err != nil {
		fmt.Printf("user %v JoinPk roomid %v get data.award_text err: %v\n", userid, roomid, err)
		return
	}
	g := strings.Index(gift, "X")
	giftName := gift[:g]
	giftNum, _ := strconv.Atoi(gift[g+1:])
	Status.OkUser(userid, giftName, giftNum)
	fmt.Printf("%v用户%v 在房间%v %v: %vX%v\n", other.TI(), userid, roomid, giftname, giftName, giftNum)
}
