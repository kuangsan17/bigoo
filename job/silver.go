package job

import (
	"fmt"
	"github.com/atwat/bigoo/other"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

//UseSilver ...
func UseSilver(btoml *other.Btoml) {
	go BuyHuima(btoml)
	go SilverAndCoin(btoml)
}

//BuyHuima ...
func BuyHuima(btoml *other.Btoml) {
	for k := range btoml.BiliUser {
		go func(k int) {
			for btoml.MoreSetting.UseSilver.BuyHuima > 0 {
				go buyHuima(k, btoml)

				time.Sleep(time.Hour * 8)
			}
		}(k)
	}
}

//SilverAndCoin ...
func SilverAndCoin(btoml *other.Btoml) {
	for k := range btoml.BiliUser {
		go func(k int) {
			for btoml.MoreSetting.UseSilver.Sliver2Coin > 0 || btoml.MoreSetting.UseCoin.Coin2Sliver > 0 {
				go silverAndCoin(k, btoml)

				time.Sleep(time.Hour * 8)
			}
		}(k)
	}
}

func silverAndCoin(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*19)
	url := "https://api.live.bilibili.com/pay/v1/Exchange/getStatus"
	data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&actionKey=appkey&appkey=1d8b6e7d45233436&build=" + build + "&channel=" + channel + "&device=android&mobi_app=android&platform=android&statistics=" + statistics + "&ts=" + other.StrTime())
	body, c := other.AppReqGet(url + "?" + data)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在检查硬币银瓜子互换状态时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v checkCoinAndSilverExchange get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v checkCoinAndSilverExchange get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在检查检查硬币银瓜子互换状态时: %v\n", other.TI(), userid, message)
		return
	}
	silver2CoinLeft, err := js.Get("data").Get("silver_2_coin_left").Int()
	if err != nil {
		fmt.Printf("user %v checkCoinAndSilverExchange get data.silver_2_coin_left err: %v\n", userid, err)
		return
	}
	coin2SilverLeft, err := js.Get("data").Get("coin_2_silver_left").Int()
	if err != nil {
		fmt.Printf("user %v checkCoinAndSilverExchange get data.coin_2_silver_left err: %v\n", userid, err)
		return
	}
	status, err := js.Get("data").Get("vip").Int()
	if err != nil {
		fmt.Printf("user %v checkCoinAndSilverExchange get data.vip err: %v\n", userid, err)
		return
	}
	//	vip, err := js.Get("data").Get("vip").Int()
	//	if err != nil {
	//		fmt.Printf("user %v checkCoinAndSilverExchange get data.vip err: %v\n", userid, err)
	//		return
	//	}
	//老爷+大会员 1 1   非老爷+非大会员 0 0
	coin2SilverNum := btoml.MoreSetting.UseCoin.Coin2Sliver + coin2SilverLeft - 25
	if status > 0 {
		coin2SilverNum = btoml.MoreSetting.UseCoin.Coin2Sliver + coin2SilverLeft - 50
	}
	s2cnum := btoml.MoreSetting.UseSilver.Sliver2Coin
	if s2cnum > 0 && silver2CoinLeft > 0 {
		for is := 0; is < s2cnum; is++ {
			go func() {
				silver2Coin(userid, btoml)
			}()
		}
	}
	if s2cnum <= 0 && coin2SilverNum > 0 && btoml.MoreSetting.UseCoin.Coin2Sliver > 0 {
		coin2Silver(userid, coin2SilverNum, btoml)
	}
}

func silver2Coin(userid int, btoml *other.Btoml) {
	url := "https://api.live.bilibili.com/pay/v1/Exchange/silver2coin"
	data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&actionKey=appkey&appkey=1d8b6e7d45233436&build=" + build + "&channel=" + channel + "&device=android&mobi_app=android&num=1&platform=android&statistics=" + statistics + "&ts=" + other.StrTime())
	body, c := other.AppReqGet(url + "?" + data)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在银瓜子换硬币时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v silver2Coin get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v silver2Coin get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在银瓜子换硬币时: %v\n", other.TI(), userid, message)
		return
	}
	fmt.Printf("%v用户%v 银瓜子换硬币成功\n", other.TI(), userid)
}

func coin2Silver(userid, num int, btoml *other.Btoml) {
	url := "https://api.live.bilibili.com/pay/v1/Exchange/coin2silver"
	data := other.Sign("access_key=" + btoml.BiliUser[userid].AccessToken + "&actionKey=appkey&appkey=1d8b6e7d45233436&build=" + build + "&channel=" + channel + "&device=android&mobi_app=android&num=" + strconv.Itoa(num) + "&platform=android&statistics=" + statistics + "&ts=" + other.StrTime())
	body, c := other.AppReqGet(url + "?" + data)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在硬币换银瓜子时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v coin2Silver get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v coin2Silver get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在硬币换银瓜子时: %v\n", other.TI(), userid, message)
		return
	}
	fmt.Printf("%v用户%v 硬币换银瓜子成功\n", other.TI(), userid)
}

func buyHuima(userid int, btoml *other.Btoml) {
	randSleep(0, 300, (userid+1)*18)
	body, c := other.PcReqGet("https://api.live.bilibili.com/lottery/v1/Ema/index", btoml.BiliUser[userid].Cookie)
	if !c {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在检查绘马时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v checkHuima get message err: %v\n", userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v checkHuima get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在检查绘马时: %v\n", other.TI(), userid, message)
		return
	}
	num1 := fmt.Sprintf("%v", js.Get("data").Get("num").Interface())
	num, err := strconv.Atoi(num1)
	if err != nil {
		fmt.Printf("user %v checkHuima get data.num err: %v\n", userid, err)
		return
	}
	buynum := btoml.MoreSetting.UseSilver.BuyHuima - num
	if buynum <= 0 {
		fmt.Printf("%v用户%v 已购买绘马\n", other.TI(), userid)
		return
	}
	csrf := ""
	a1 := strings.Index(btoml.BiliUser[userid].Cookie, "bili_jct=")
	if a1 != -1 {
		s1 := btoml.BiliUser[userid].Cookie[a1+9:]
		a2 := strings.Index(s1, ";")
		if a2 != -1 {
			csrf = s1[:a2]
		}
	}
	url := "https://api.live.bilibili.com/lottery/v1/Ema/buy"
	data := "coinType=silver&num=" + strconv.Itoa(buynum) + "&csrf_token=" + csrf + "&csrf=" + csrf + "&visit_id="
	body, c = other.PcReqPost(url, data, btoml.BiliUser[userid].Cookie)
	if !c {
		return
	}
	js, err = simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在购买绘马时返回了非json字符串\n", other.TI(), userid)
		return
	}
	message, err = js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v buyHuima get message err: %v\n", userid, err)
		return
	}
	code, err = js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v buyHuima get code err: %v\n", userid, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在购买绘马时: %v\n", other.TI(), userid, message)
		return
	}
	fmt.Printf("%v用户%v 购买%v个绘马成功\n", other.TI(), userid, buynum)
}
