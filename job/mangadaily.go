package job

import (
	"fmt"
	"github.com/atwat/bigoo/other"
	"time"

	"github.com/bitly/go-simplejson"
)

//MangaDaily ...
func MangaDaily(btoml *other.Btoml) {
	for k := range btoml.BiliUser {
		go func(k int) {
			for btoml.Setting.MangaDailyJob {
				go mangaSign(k, btoml)
				go mangaShare(k, btoml)
				go mangaSvipMonth(k, btoml)
				time.Sleep(time.Hour * 8)
			}
		}(k)
	}
}

func mangaSign(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*16)
	url := "https://manga.bilibili.com/twirp/activity.v1.Activity/GetClockInInfo"
	data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&appkey=1d8b6e7d45233436&build=32000003&channel=huawei&device=android&mobi_app=android_comic&platform=android&ts=" + other.StrTime() + "&version=3.2.1")
	body, c := other.AppReqPost(url, data)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在漫画签到时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("msg").String()
	if err != nil {
		fmt.Printf("user %v mangaSign get msg err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v mangaSign get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在漫画签到时: %v\n", other.TI(), userid, message)
		return
	}
	isSigned, err := js.Get("data").Get("status").Int()
	if err != nil {
		fmt.Printf("user %v mangaSign get data.status err: %v\n", userid, err)
		return
	}
	dayCount, err := js.Get("data").Get("day_count").Int()
	if err != nil {
		fmt.Printf("user %v mangaSign get data.day_count err: %v\n", userid, err)
		return
	}
	if isSigned == 1 {
		//fmt.Printf("%v用户%v 在漫画已签到 已连续签到%v天\n", other.TI(), userid, dayCount)
	} else {
		//fmt.Printf("%v用户%v 在漫画未签到 已连续签到%v天\n", other.TI(), userid, dayCount)
		url := "https://manga.bilibili.com/twirp/activity.v1.Activity/ClockIn"
		data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&appkey=1d8b6e7d45233436&build=32000003&channel=huawei&device=android&mobi_app=android_comic&platform=android&ts=" + other.StrTime() + "&version=3.2.1")
		body, c := other.AppReqPost(url, data)
		if !c {
			return
		}
		js, err := simplejson.NewJson(body)
		if err != nil {
			fmt.Printf("%v用户%v 在漫画签到时返回了非json字符串\n", other.TI(), userid)
			return
		}
		message, err := js.Get("msg").String()
		if err != nil {
			fmt.Printf("user %v mangaSign get msg err: %v\n", userid, err)
			return
		}

		if message != "" {
			fmt.Printf("%v用户%v 在漫画签到时: %v\n", other.TI(), userid, message)
			return
		}
		fmt.Printf("%v用户%v 在漫画签到成功 已连续签到%v天\n", other.TI(), userid, dayCount)
		return
	}
}

func mangaShare(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*17)
	url := "https://manga.bilibili.com/twirp/activity.v1.Activity/ShareComic"
	data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&appkey=1d8b6e7d45233436&build=32000003&channel=huawei&device=android&mobi_app=android_comic&platform=android&ts=" + other.StrTime() + "&version=3.2.1")
	body, c := other.AppReqPost(url, data)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在漫画分享时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("msg").String()
	if err != nil {
		fmt.Printf("user %v mangaShare get msg err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v mangaShare get code err: %v\n", userid, err)
		return
	}
	if code == 1 {
		//fmt.Printf("%v用户%v 在漫画已分享\n", other.TI(), userid)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在漫画分享时: %v\n", other.TI(), userid, message)
		return
	}
	point, err := js.Get("data").Get("point").Int()
	if err != nil {
		fmt.Printf("user %v mangaShare get data.point err: %v\n", userid, err)
		return
	}
	fmt.Printf("%v用户%v 漫画分享成功 获得%v积分\n", other.TI(), userid, point)
}
func mangaSvipMonth(userid int, btoml *other.Btoml) {
	if time.Now().AddDate(0, 0, 1).In(time.FixedZone("CST", 28800)).Format("02") != "01" { //本月最后一天
		return
	}
	randSleep(0, 300, (userid+1)*24)
	url := "https://manga.bilibili.com/twirp/user.v1.User/GetVipRewardComics"
	data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&appkey=1d8b6e7d45233436&build=32000003&c_locale=&channel=huawei&device=android&is_teenager=0&machine=huawei&mobi_app=android_comic&platform=android&s_locale=&ts=" + other.StrTime() + "&version=3.2.1")
	body, c := other.AppReqPost(url, data)
	if !c {
		return
	}
	//fmt.Println(string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在漫画大会员福利检查时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("msg").String()
	if err != nil {
		fmt.Printf("user %v GetVipRewardComics get msg err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v GetVipRewardComics get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		if code != 1 {
			fmt.Printf("%v用户%v 在漫画大会员福利检查时: %v\n", other.TI(), userid, message)
		}
		return
	}
	gotsp, err := js.Get("data").Get("is_got_sp").Bool()
	if err != nil {
		fmt.Printf("user %v GetVipRewardComics get data.is_got_sp err: %v\n", userid, err)
		return
	}
	if !gotsp {
		url2 := "https://manga.bilibili.com/twirp/user.v1.User/GetVipReward"
		data2 := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&appkey=1d8b6e7d45233436&build=32000003&c_locale=&channel=huawei&device=android&is_teenager=0&machine=huawei&mobi_app=android_comic&platform=android&reason_id=1&s_locale=&ts=" + other.StrTime() + "&type=0&version=3.2.1")
		body2, c := other.AppReqPost(url2, data2)
		if !c {
			return
		}
		//fmt.Println(string(body2))
		js2, err := simplejson.NewJson(body2)
		if err != nil {
			fmt.Printf("%v用户%v 在漫画大会员福利领取时返回了非json字符串\n", other.TI(), userid)
			return
		}
		message2, err := js2.Get("msg").String()
		if err != nil {
			fmt.Printf("user %v GetVipReward get msg err: %v\n", userid, err)
			return
		}
		code2, err := js2.Get("code").Int()
		if err != nil {
			fmt.Printf("user %v GetVipReward get code err: %v\n", userid, err)
			return
		}
		if code2 != 0 {
			fmt.Printf("%v用户%v 在漫画大会员福利领取时: %v\n", other.TI(), userid, message2)
			return
		}
		amount, err := js2.Get("data").Get("amount").Int()
		if err != nil {
			fmt.Printf("user %v GetVipReward get data.amount err: %v\n", userid, err)
			return
		}
		fmt.Printf("%v用户%v 在漫画大会员福利领取获得: 福利券X%v\n", other.TI(), userid, amount)
	}
}

func mangaCouponsClean(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*25)
	if !btoml.MoreSetting.CleanExpiringMangaCoupons {
		return
	}

}

func mangaGetCoupons(userid int, btoml *other.Btoml) {

}
