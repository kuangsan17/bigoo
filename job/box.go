package job

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/atwat/bigoo/other"

	"github.com/bitly/go-simplejson"
)

func joinBox(gid int, btoml *other.Btoml, Status *other.Status) {
	//fmt.Println(gid)
	for userid := range btoml.BiliUser {
		go func(userid int) {
			randSleep(0, 10, (userid+1)*25)
			boxGetStatus(gid, btoml, userid, Status)
		}(userid)
	}
}

/*
	-1 未开始
	0  开始未抽
	1  开始已抽
	2
	3  结束未抽
	4  结束已抽
*/
func boxGetStatus(gid int, btoml *other.Btoml, userid int, Status *other.Status) {
	url := other.Sign("https://api.live.bilibili.com/xlive/lottery-interface/v2/Box/getStatus?access_key=" + btoml.BiliUser[userid].AccessToken + "&actionKey=appkey&aid=" + strconv.Itoa(gid) + "&appkey=1d8b6e7d45233436&build=" + build + "&c_locale=zh-Hans_CN&channel=" + channel + "&device=android&mobi_app=android&platform=android&s_locale=zh-Hans_CN&statistics=" + statistics + "&ts=" + other.StrTime())
	body, c := other.AppReqGet(url)
	if !c {
		return
	}
	//fmt.Println(string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Print(other.TI(), "宝箱抽奖", gid, " 获取状态时返回了非json字符串\n")
		fmt.Println(string(body))
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v box %v GetStatus json message err :%v\n", userid, gid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v box %v GetStatus json code err :%v\n", userid, gid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在宝箱抽奖%v获取状态时: %v\n", other.TI(), userid, gid, message)
		return
	}
	boxTitle, err := js.Get("data").Get("title").String()
	if err != nil {
		fmt.Printf("user %v box %v GetStatus json data.title err :%v\n", userid, gid, err)
		return
	}
	typeB, err := js.Get("data").Get("typeB").Array()
	if err != nil {
		fmt.Printf("user %v box %v GetStatus json data.typeB err :%v\n", userid, gid, err)
		return
	}
	allid := len(typeB)
	for i := range typeB {
		status, err := js.Get("data").Get("typeB").GetIndex(i).Get("status").Int()
		if err != nil {
			fmt.Printf("user %v box %v GetStatus json data.typeB.%v.status err :%v\n", userid, gid, i, err)
			continue
		}
		if status != -1 && status != 0 {
			continue
		}
		joinStartTime, err := js.Get("data").Get("typeB").GetIndex(i).Get("join_start_time").Int()
		if err != nil {
			fmt.Printf("user %v box %v GetStatus json data.typeB.%v.join_start_time err :%v\n", userid, gid, i, err)
			continue
		}
		joinEndTime, err := js.Get("data").Get("typeB").GetIndex(i).Get("join_end_time").Int()
		if err != nil {
			fmt.Printf("user %v box %v GetStatus json data.typeB.%v.join_end_time err :%v\n", userid, gid, i, err)
			continue
		}
		nowTime := other.IntTime()
		start, end := joinStartTime-nowTime, joinEndTime-nowTime
		go boxJoin(gid, allid, btoml, userid, i+1, start, end, boxTitle, Status)
	}
}

func boxJoin(gid, allnumber int, btoml *other.Btoml, userid, number, start, end int, boxTitle string, Status *other.Status) {
	//fmt.Println(userid, "JOINbox ...", start, end)
	if end <= 0 {
		return
	}
	if rand.Intn(100)+1 > btoml.Setting.Lottery.BoxOdds {
		return //随机放弃
	}
	go func() {
		randSleep(end+30, end+50, gid*(userid+3))
		boxGetWinner(gid, number, userid, btoml, Status)
	}()
	randSleep(start+(end-start)/5, start+(end-start)/2, gid*(userid+3))
	url := other.Sign("https://api.live.bilibili.com/xlive/lottery-interface/v2/Box/draw?access_key=" + btoml.BiliUser[userid].AccessToken + "&actionKey=appkey&aid=" + strconv.Itoa(gid) + "&appkey=1d8b6e7d45233436&build=" + build + "&c_locale=zh-Hans_CN&channel=" + channel + "&device=android&mobi_app=android&number=" + strconv.Itoa(number) + "&platform=android&s_locale=zh-Hans_CN&statistics=" + statistics + "&ts=" + other.StrTime())
	body, c := other.AppReqGet(url)
	if !c {
		return
	}
	//fmt.Println(string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在%v第%v(共%v)轮抽奖时:not json\n", other.TI(), userid, boxTitle, number, allnumber)
		fmt.Println(string(body))
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("%v用户%v 在%v第%v(共%v)轮抽奖时:json message %v\n", other.TI(), userid, boxTitle, number, allnumber, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("%v用户%v 在%v第%v(共%v)轮抽奖时:json code %v\n", other.TI(), userid, boxTitle, number, allnumber, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在%v第%v(共%v)轮抽奖时: %v\n", other.TI(), userid, boxTitle, number, allnumber, message)
		return
	}
	fmt.Printf("%v用户%v 在%v第%v(共%v)轮抽奖 参与成功\n", other.TI(), userid, boxTitle, number, allnumber)
}

func boxGetWinner(gid, number, userid int, btoml *other.Btoml, Status *other.Status) {
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v2/Box/getWinnerGroupInfo?aid=" + strconv.Itoa(gid) + "&number=" + strconv.Itoa(number)
	body, c := other.PcReqGet(url, btoml.BiliUser[userid].Cookie)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:not json\n", other.TI(), userid, gid, number)
		fmt.Println(string(body))
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:json message %v\n", other.TI(), userid, gid, number, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:json code %v\n", other.TI(), userid, gid, number, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时: %v\n", other.TI(), userid, gid, number, message)
		return
	}
	myuid, err := js.Get("data").Get("uid").Int()
	if err != nil {
		fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:json data.uid %v\n", other.TI(), userid, gid, number, err)
		return
	}
	groups, err := js.Get("data").Get("groups").Array()
	if err != nil {
		fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:json data.groups %v\n", other.TI(), userid, gid, number, err)
		return
	}
	iswiner := false
	gift := ""
	for groupid := range groups {
		giftTitle, err := js.Get("data").Get("groups").GetIndex(groupid).Get("giftTitle").String()
		if err != nil {
			fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:json data.groups.%v.giftTitle %v\n", other.TI(), userid, gid, number, groupid, err)
			continue
		}
		if gift == "" {
			gift = giftTitle
		} else {
			gift = gift + "、" + giftTitle
		}
		list, err := js.Get("data").Get("groups").GetIndex(groupid).Get("list").Array()
		if err != nil {
			fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:json data.groups.%v.list %v\n", other.TI(), userid, gid, number, groupid, err)
			continue
		}
		for listid := range list {
			owneruid, err := js.Get("data").Get("groups").GetIndex(groupid).Get("list").GetIndex(listid).Get("uid").Int()
			if err != nil {
				fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 检查获奖时:json data.groups.%v.list.%v.uid %v\n", other.TI(), userid, gid, number, groupid, listid, err)
				continue
			}
			if owneruid == myuid {
				iswiner = true
				Status.OkUser(userid, giftTitle, 1)
				fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 获得%v\n", other.TI(), userid, gid, number, giftTitle)
				other.SendMail(fmt.Sprintf("用户%v在宝箱抽奖中获奖", btoml.BiliUser[userid].UserName), btoml)
			} else {

			}
		}
	}
	if !iswiner {
		fmt.Printf("%v用户%v 在宝箱%v第%v轮抽奖 未获得%v\n", other.TI(), userid, gid, number, gift)
		//other.SendMail(fmt.Sprintf("用户%v在宝箱抽奖中未获奖", btoml.BiliUser[userid].UserName), btoml)
	}
}
