package job

import (
	"fmt"
	"github.com/atwat/bigoo/other"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

//LiveDaily ...
func LiveDaily(btoml *other.Btoml) {
	for k := range btoml.BiliUser {
		go func(k int) {
			for btoml.Setting.LiveDailyJob {
				go liveSign(k, btoml)
				//go doubleWatch(k, btoml)
				//go silverBox(k, btoml)
				go fanGroup(k, btoml)
				go dailyBag(k, btoml)
				time.Sleep(time.Hour * 8)
			}
		}(k)
	}
}

//LiveHeart ...
func LiveHeart(btoml *other.Btoml) {
	for k := range btoml.BiliUser {
		go func(k int) {
			randSleep(0, 300, (k+1)*15)
			go func() {
				for btoml.Setting.LiveOnlineHeart {
					randSleep(5, 10, (k+1)*15)
					SmallHeart(k, btoml)
					time.Sleep(time.Minute * 5)
				}
			}()
			if btoml.Setting.LiveOnlineHeart {
				fmt.Printf("%v用户%v 开始模拟观看直播\n", other.TI(), k)
				for btoml.Setting.LiveOnlineHeart {
					go func() {
						url := `https://api.live.bilibili.com/heartbeat/v1/OnLine/mobileOnline?` + other.Sign(`access_key=`+btoml.BiliUser[k].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
						data := `room_id=96136&scale=xxhdpi`
						body, c := other.AppReqPost(url, data)
						if !c {
							runtime.Goexit()
						}
						js, err := simplejson.NewJson(body)
						if err != nil {
							fmt.Printf("%v用户%v 在直播心跳时返回了非json字符串\n", other.TI(), k)
							runtime.Goexit()
						}
						message, err := js.Get("message").String()
						if err != nil {
							fmt.Printf("user %v liveHeart get message err: %v\n", k, err)
							runtime.Goexit()
						}
						code, err := js.Get("code").Int()
						if err != nil {
							fmt.Printf("user %v liveHeart get code err: %v\n", k, err)
							runtime.Goexit()
						}
						if code != 0 {
							fmt.Printf("%v用户%v 在直播心跳时: %v\n", other.TI(), k, message)
							runtime.Goexit()
						}
						//fmt.Printf("%v用户%v 直播心跳成功\n", other.TI(), k)
					}()
					time.Sleep(time.Minute * 5)
				}
			}
		}(k)
	}
}

func dailyBag(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*21)
	url := `https://api.live.bilibili.com/gift/v2/live/m_receive_daily_bag?` + other.Sign(`access_key=`+btoml.BiliUser[userid].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
	body, c := other.AppReqGet(url)
	if !c {
		return
	}
	//fmt.Println("user", userid, string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在获取直播勋章礼包时返回了非json字符串\n", other.TI(), userid)
		return
	}
	msg, err := js.Get("msg").String()
	if err != nil {
		fmt.Printf("user %v dailyBag get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v dailyBag get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在获取直播勋章礼包时: %v\n", other.TI(), userid, msg)
		return
	}
	data, err := js.Get("data").Array()
	if err != nil {
		fmt.Printf("user %v dailyBag get data err: %v\n", userid, err)
		return
	}
	for i := range data {
		bagName, err := js.Get("data").GetIndex(i).Get("bag_name").String()
		if err != nil {
			fmt.Printf("user %v dailyBag get data.%v.bag_name err: %v\n", userid, i, err)
			continue
		}
		bagSource, err := js.Get("data").GetIndex(i).Get("bag_source").String()
		if err != nil {
			fmt.Printf("user %v dailyBag get data.%v.bag_source err: %v\n", userid, i, err)
			continue
		}
		/*typee, err := js.Get("data").GetIndex(i).Get("type").Int()
		if err != nil {
			fmt.Printf("user %v dailyBag get data.%v.type err: %v\n", userid, i, err)
			continue
		}*/
		fmt.Printf("%v用户%v 领取%v的%v\n", other.TI(), userid, bagSource, bagName)
	}
}

func liveSign(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*11)
	url := `https://api.live.bilibili.com/rc/v2/Sign/getSignInfo?` + other.Sign(`access_key=`+btoml.BiliUser[userid].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
	body, c := other.AppReqGet(url)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在直播签到时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v liveSign get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v liveSign get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在直播签到时: %v\n", other.TI(), userid, message)
		return
	}
	isSigned, err := js.Get("data").Get("is_signed").Bool()
	if err != nil {
		fmt.Printf("user %v liveSign get data.is_signed err: %v\n", userid, err)
		return
	}
	/*days, err := js.Get("data").Get("days").Int()
	if err != nil {
		fmt.Printf("user %v liveSign get data.days err: %v\n", userid, err)
		return
	}
	signDays, err := js.Get("data").Get("sign_days").Int()
	if err != nil {
		fmt.Printf("user %v liveSign get data.sign_days err: %v\n", userid, err)
		return
	}*/
	if isSigned {
		//fmt.Printf("%v用户%v 在直播已签到 (本月%v/%v)\n", other.TI(), userid, signDays, days)
	} else {
		//fmt.Printf("%v用户%v 在直播未签到 (本月%v/%v)\n", other.TI(), userid, signDays, days)
		url := `https://api.live.bilibili.com/rc/v1/Sign/doSign?` + other.Sign(`access_key=`+btoml.BiliUser[userid].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
		body, c := other.AppReqGet(url)
		if !c {
			return
		}
		js, err := simplejson.NewJson(body)
		if err != nil {
			fmt.Printf("%v用户%v 在直播签到时返回了非json字符串\n", other.TI(), userid)
			return
		}
		message, err := js.Get("message").String()
		if err != nil {
			fmt.Printf("user %v liveSign get message err: %v\n", userid, err)
			return
		}
		code, err := js.Get("code").Int()
		if err != nil {
			fmt.Printf("user %v liveSign get code err: %v\n", userid, err)
			return
		}
		if code != 0 {
			fmt.Printf("%v用户%v 在直播签到时: %v\n", other.TI(), userid, message)
			return
		}
		text, err := js.Get("data").Get("text").String()
		if err != nil {
			fmt.Printf("user %v liveSign get data.text err: %v\n", userid, err)
			return
		}
		specialText, err := js.Get("data").Get("specialText").String()
		if err != nil {
			fmt.Printf("user %v liveSign get data.specialText err: %v\n", userid, err)
			return
		}
		allDays, err := js.Get("data").Get("allDays").Int()
		if err != nil {
			fmt.Printf("user %v liveSign get data.allDays err: %v\n", userid, err)
			return
		}
		hadSignDays, err := js.Get("data").Get("hadSignDays").Int()
		if err != nil {
			fmt.Printf("user %v liveSign get data.hadSignDays err: %v\n", userid, err)
			return
		}
		fmt.Printf("%v用户%v 在直播签到获得%v%v (本月%v/%v)\n", other.TI(), userid, text, specialText, hadSignDays, allDays)
	}
}
func doubleWatch(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*12)
	csrf := ""
	a1 := strings.Index(btoml.BiliUser[userid].Cookie, "bili_jct=")
	if a1 != -1 {
		s1 := btoml.BiliUser[userid].Cookie[a1+9:]
		a2 := strings.Index(s1, ";")
		if a2 != -1 {
			csrf = s1[:a2]
		}
	}
	_, c := other.PcReqPost(`https://api.live.bilibili.com/User/userOnlineHeart`, "csrf_token="+csrf+"&csrf="+csrf+"&visit_id= ", btoml.BiliUser[userid].Cookie)
	if !c {
		return
	}
	url := `https://api.live.bilibili.com/heartbeat/v1/OnLine/mobileOnline?` + other.Sign(`access_key=`+btoml.BiliUser[userid].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
	data := `room_id=96136&scale=xxhdpi`
	_, c = other.AppReqPost(url, data)
	if !c {
		return
	}
	url = `https://api.live.bilibili.com/activity/v1/task/receive_award?task_id=double_watch_task&` + other.Sign(`access_key=`+btoml.BiliUser[userid].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
	body, c := other.AppReqGet(url)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在领取双端奖励时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v doubleWatch get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v doubleWatch get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在领取双端奖励时: %v\n", other.TI(), userid, message)
		return
	}
	fmt.Printf("%v用户%v 领取双端奖励成功\n", other.TI(), userid)
}
func silverBox(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*13)
	for {
		url := `https://api.live.bilibili.com/lottery/v1/SilverBox/getCurrentTask?` + other.Sign(`access_key=`+btoml.BiliUser[userid].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
		body, c := other.AppReqGet(url)
		if !c {
			return
		}
		//	fmt.Println(string(body))
		js, err := simplejson.NewJson(body)
		if err != nil {
			fmt.Printf("%v用户%v 在检查银瓜子宝箱时返回了非json字符串\n", other.TI(), userid)
			return
		}
		message, err := js.Get("message").String()
		if err != nil {
			fmt.Printf("user %v checkSilverBox get message err: %v\n", userid, err)
			return
		}
		code, err := js.Get("code").Int()
		if err != nil {
			fmt.Printf("user %v checkSilverBox get code err: %v\n", userid, err)
			return
		}
		if code != 0 {
			if strings.Index(message, "已经") == -1 {
				fmt.Printf("%v用户%v 在检查银瓜子宝箱时: %v\n", other.TI(), userid, message)
			}
			return
		}
		timeEnd, err := js.Get("data").Get("time_end").Int()
		if err != nil {
			fmt.Printf("user %v checkSilverBox get data.time_end err: %v\n", userid, err)
			return
		}
		times, err := js.Get("data").Get("times").Int()
		if err != nil {
			fmt.Printf("user %v checkSilverBox get data.times err: %v\n", userid, err)
			return
		}
		maxTimes, err := js.Get("data").Get("max_times").Int()
		if err != nil {
			fmt.Printf("user %v checkSilverBox get data.max_times err: %v\n", userid, err)
			return
		}
		silver, err := js.Get("data").Get("silver").Int()
		if err != nil {
			fmt.Printf("user %v checkSilverBox get data.silver err: %v\n", userid, err)
			return
		}
		//fmt.Printf("%v用户%v 等待领取第%v(共%v)轮的%v银瓜子宝箱\n", other.TI(), userid, times, maxTimes, silver)
		randSleep(timeEnd-other.IntTime()+3, timeEnd-other.IntTime()+8, userid+1)
		url2 := `https://api.live.bilibili.com/lottery/v1/SilverBox/getAward?` + other.Sign(`access_key=`+btoml.BiliUser[userid].AccessToken+`&actionKey=appkey&appkey=1d8b6e7d45233436&build=`+build+`&channel=`+channel+`&device=android&mobi_app=android&platform=android&statistics=`+statistics+`&ts=`+other.StrTime())
		body2, c := other.AppReqGet(url2)
		if !c {
			return
		}
		//	fmt.Println(string(body2))
		js2, err := simplejson.NewJson(body2)
		if err != nil {
			fmt.Printf("%v用户%v 在打开第%v(共%v)轮的%v银瓜子宝箱时返回了非json字符串\n", other.TI(), userid, times, maxTimes, silver)
			return
		}
		message2, err := js2.Get("message").String()
		if err != nil {
			fmt.Printf("user %v openSilverBox get message err: %v\n", userid, err)
			return
		}
		code2, err := js2.Get("code").Int()
		if err != nil {
			fmt.Printf("user %v openSilverBox get code err: %v\n", userid, err)
			return
		}
		if code2 != 0 {
			fmt.Printf("%v用户%v 在打开第%v(共%v)轮的%v银瓜子宝箱时: %v\n", other.TI(), userid, times, maxTimes, silver, message2)
			return
		}
		awardSilver := js2.Get("data").Get("awardSilver").Interface()
		/*if err != nil {
			fmt.Printf("user %v openSilverBox get data.awardSilver err: %v\n", userid, err)
			return
		}*/
		fmt.Printf("%v用户%v 打开第%v(共%v)轮的银瓜子宝箱 获得银瓜子X%v\n", other.TI(), userid, times, maxTimes, awardSilver)
	}
}
func fanGroup(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*14)
	body, c := other.PcReqGet("https://api.vc.bilibili.com/link_group/v1/member/my_groups", btoml.BiliUser[userid].Cookie)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在获取应援团时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v getGroups get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v getGroups get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在获取应援团时: %v\n", other.TI(), userid, message)
		return
	}
	list, err := js.Get("data").Get("list").Array()
	if err != nil {
		fmt.Printf("user %v getGroups get data.list err: %v\n", userid, err)
		return
	}
	for i := range list {
		time.Sleep(time.Second * 5)
		groupID, err := js.Get("data").Get("list").GetIndex(i).Get("group_id").Int()
		if err != nil {
			fmt.Printf("user %v getGroups get data.list.%v.group_id err: %v\n", userid, i, err)
			continue
		}
		groupName, err := js.Get("data").Get("list").GetIndex(i).Get("group_name").String()
		if err != nil {
			fmt.Printf("user %v getGroups get data.list.%v.group_name err: %v\n", userid, i, err)
			continue
		}
		fansMedalName, err := js.Get("data").Get("list").GetIndex(i).Get("fans_medal_name").String()
		if err != nil {
			fmt.Printf("user %v getGroups get data.list.%v.fans_medal_name err: %v\n", userid, i, err)
			continue
		}
		ownerUID, err := js.Get("data").Get("list").GetIndex(i).Get("owner_uid").Int()
		if err != nil {
			fmt.Printf("user %v getGroups get data.list.%v.owner_uid err: %v\n", userid, i, err)
			continue
		}
		data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&group_id=" + strconv.Itoa(groupID) + "&owner_id=" + strconv.Itoa(ownerUID) + "&ts=" + other.StrTime())
		url := "https://api.vc.bilibili.com/link_setting/v1/link_setting/sign_in?" + data
		body, c := other.AppReqGet(url)
		if !c {
			continue
		}
		js2, err := simplejson.NewJson(body)
		if err != nil {
			fmt.Printf("%v用户%v 在%v 应援时返回了非json字符串\n", other.TI(), userid, groupName)
			continue
		}
		message, err := js2.Get("message").String()
		if err != nil {
			fmt.Printf("user %v signGroups %v get message err: %v\n", userid, groupName, err)
			continue
		}
		code, err := js2.Get("code").Int()
		if err != nil {
			fmt.Printf("user %v signGroups %v get code err: %v\n", userid, groupName, err)
			continue
		}
		if code != 0 {
			fmt.Printf("%v用户%v 在%v 应援时: %v\n", other.TI(), userid, groupName, message)
			continue
		}
		bstatus, err := js2.Get("data").Get("status").Int()
		if err != nil {
			fmt.Printf("user %v signGroups %v get data.status err: %v\n", userid, groupName, err)
			continue
		}
		if bstatus != 0 {
			//fmt.Printf("%v用户%v 在%v 已应援过\n", other.TI(), userid, groupName)
			continue
		}
		addNum, err := js2.Get("data").Get("add_num").Int()
		if err != nil {
			fmt.Printf("user %v signGroups %v get data.add_num err: %v\n", userid, groupName, err)
			continue
		}
		fmt.Printf("%v用户%v 在%v 应援获得: %v亲密度+%v\n", other.TI(), userid, groupName, fansMedalName, addNum)
	}
}
