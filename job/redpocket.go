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

//JoinRed ...
func JoinRed(roomid, redID int, btoml *other.Btoml, Status *other.Status, userid, redtime int) {
	if redtime <= 0 {
		return
	}
	giftname := "红包"
	randSleep(0, redtime*2/3, redID*(userid+2))
	rand.Seed(time.Now().UnixNano() - int64((userid+2)*redID))
	if rand.Intn(100)+1 > btoml.Setting.Lottery.RedOdds {
		return //随机放弃
	}
	if !Status.IfUserOK(userid) {
		return //黑屋中
	}
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v1/red_pocket/join"
	//access_key=300ed1f36c8e2a9f48d7d45b534c0c21&actionKey=appkey&appkey=1d8b6e7d45233436&build=6182200&c_locale=zh_CN&channel=pairui01&device=android&id=81&mobi_app=android&platform=android&roomId=22747055&s_locale=zh_CN&statistics=%7B%22appId%22%3A1%2C%22platform%22%3A3%2C%22version%22%3A%226.18.2%22%2C%22abtest%22%3A%22%22%7D&ts=1613042491&version=6.18.2&sign=db78f088cc3abbebfbba482899986308
	data := other.Sign(`access_key=` + btoml.BiliUser[userid].AccessToken + `&actionKey=appkey&appkey=1d8b6e7d45233436&build=` + build + `&c_locale=zh_CN&channel=` + channel + `&device=android&id=` + strconv.Itoa(redID) + `&mobi_app=android&platform=android&roomId=` + strconv.Itoa(roomid) + `&s_locale=zh_CN&statistics=` + statistics + `&ts=` + other.StrTime() + `&version=` + version)
	body, c := other.AppReqPost(url, data)
	if !c {
		return
	}
	fmt.Println(string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在房间%v %v err: not json\n", other.TI(), userid, roomid, giftname)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v JoinRed roomid %v get message err: %v\n", userid, roomid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v JoinRed roomid %v get code err: %v\n", userid, roomid, err)
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
	hash, err := js.Get("data").Get("server_hash").Int()
	if err != nil {
		fmt.Printf("user %v JoinRed roomid %v get data.server_hash err: %v\n", userid, roomid, err)
		return
	}
	fmt.Printf("%v用户%v 在房间%v %v 参与成功\n", other.TI(), userid, roomid, giftname)
	time.Sleep(time.Second * time.Duration(redtime))
	redResult(roomid, redID, btoml, Status, userid, redtime, hash)
}

func redResult(roomid, redID int, btoml *other.Btoml, Status *other.Status, userid, redtime, serverHash int) {
	giftname := "红包"
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v1/red_pocket/get_result"
	data := other.Sign(`access_key=` + btoml.BiliUser[userid].AccessToken + `&actionKey=appkey&appkey=1d8b6e7d45233436&build=` + build + `&c_locale=zh_CN&channel=` + channel + `&device=android&id=` + strconv.Itoa(redID) + `&mobi_app=android&platform=android&roomId=` + strconv.Itoa(roomid) + `&s_locale=zh_CN&serverHash=` + strconv.Itoa(serverHash) + `&statistics=` + statistics1 + `&ts=` + other.StrTime() + `&version=` + version)
	body, c := other.AppReqGet(url + "？" + data)
	if !c {
		return
	}
	fmt.Println(url+"？"+data, "\n", string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在房间%v %v err: not json\n", other.TI(), userid, roomid, giftname)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v redResult roomid %v get message err: %v\n", userid, roomid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v redResult roomid %v get code err: %v\n", userid, roomid, err)
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
	num, err := js.Get("data").Get("num").Int()
	if err != nil {
		fmt.Printf("user %v redResult roomid %v get data.num err: %v\n", userid, roomid, err)
		return
	}
	if num > 0 {
		Status.OkUser(userid, "金瓜子", num)
		fmt.Printf("%v用户%v 在房间%v %v 获奖%v\n", other.TI(), userid, roomid, giftname, num)
	}
}
