package job

import (
	"fmt"
	"github.com/atwat/bigoo/other"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

//MainSvipDaily ...
func MainSvipDaily(btoml *other.Btoml) {
	for k := range btoml.BiliUser {
		go func(k int) {
			for btoml.Setting.MangaDailyJob {
				go svipReceive(k, btoml)
				time.Sleep(time.Hour * 25)
			}
		}(k)
	}
}

func svipReceive(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*15)
	body0, c0 := other.PcReqGet("https://api.bilibili.com/x/vip/web/user/info?jsonp=jsonp", btoml.BiliUser[userid].Cookie)
	if !c0 {
		return
	}
	js0, err := simplejson.NewJson(body0)
	if err != nil {
		return
	}
	vipType, err := js0.Get("data").Get("vip_type").Int()
	if err != nil {
		return
	}
	if vipType != 2 {
		return
	}
	url := "https://api.bilibili.com/x/vip/privilege/my"
	body, c := other.PcReqGet(url, btoml.BiliUser[userid].Cookie)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在检查年度大会员卡券包时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v svipReceive get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v svipReceive get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在检查年度大会员卡券包时: %v\n", other.TI(), userid, message)
		return
	}
	list, err := js.Get("data").Get("list").Array()
	if err != nil {
		fmt.Printf("user %v svipReceive get data.list err: %v\n", userid, err)
		return
	}
	for i := range list {
		itype, err := js.Get("data").Get("list").GetIndex(i).Get("type").Int()
		if err != nil {
			fmt.Printf("user %v svipReceive get data.list.%v.type err: %v\n", userid, i, err)
			return
		}
		istate, err := js.Get("data").Get("list").GetIndex(i).Get("state").Int()
		if err != nil {
			fmt.Printf("user %v svipReceive get data.list.%v.state err: %v\n", userid, i, err)
			return
		}
		if istate == 1 {
			if itype == 1 {
				//fmt.Printf("%v用户%v 已领取过年度大会员专享B币\n", other.TI(), userid)
			} else if itype == 2 {
				//fmt.Printf("%v用户%v 已领取过年度大会员专享会员购优惠券\n", other.TI(), userid)
			} else {
				//fmt.Printf("%v用户%v 已领取过年度大会员专享礼包(type:%v)\n", other.TI(), userid, itype)
			}
		} else {
			svipReceiveRec(userid, itype, btoml)
		}
	}
}

func svipReceiveRec(userid, itype int, btoml *other.Btoml) {
	csrf := ""
	a1 := strings.Index(btoml.BiliUser[userid].Cookie, "bili_jct=")
	if a1 != -1 {
		s1 := btoml.BiliUser[userid].Cookie[a1+9:]
		a2 := strings.Index(s1, ";")
		if a2 != -1 {
			csrf = s1[:a2]
		}
	}
	url := "https://api.bilibili.com/x/vip/privilege/receive"
	data := fmt.Sprintf("type=%v&csrf=%v&csrf_token=%v", itype, csrf, csrf)
	body, c := other.PcReqPost(url, data, btoml.BiliUser[userid].Cookie)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在领取年度大会员卡券包时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v svipReceiveRec get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v svipReceiveRec get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在领取年度大会员卡券包时: %v\n", other.TI(), userid, message)
		return
	}
	if itype == 1 {
		fmt.Printf("%v用户%v 成功领取年度大会员专享B币\n", other.TI(), userid)
	} else if itype == 2 {
		fmt.Printf("%v用户%v 成功领取年度大会员专享会员购优惠券\n", other.TI(), userid)
	} else {
		fmt.Printf("%v用户%v 成功领取年度大会员专享礼包(type:%v)\n", other.TI(), userid, itype)
	}
}
