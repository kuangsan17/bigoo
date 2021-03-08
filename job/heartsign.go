package job

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	rand2 "math/rand"
	"github.com/atwat/bigoo/other"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/ebfe/keccak"
	"github.com/jzelinskie/whirlpool"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/ripemd160"
)

func randuuid() string {
	u1, err := uuid.NewV4()
	if err == nil {
		return u1.String()
	}
	return randuuid()
}

//SmallHeart ...
func SmallHeart(k int, btoml *other.Btoml) {
	roomID := "7117440"
	upID := "146563374"
	medals := getUserMedalsPc(btoml.BiliUser[k].Cookie)
	medals = getUserMedalsPcRealRoomID(medals)
	if len(medals) != 0 {
		roomID = strconv.Itoa(medals[0].roomID)
		upID = strconv.Itoa(medals[0].targetID)
	}
	for _, medal := range medals {
		if medal.iconCode == 1 {
			roomID = strconv.Itoa(medal.roomID)
			upID = strconv.Itoa(medal.targetID)
		}
	}
	/*for _, medal := range medals {
		go hmobileEntry(k, btoml, strconv.Itoa(medal.roomID), strconv.Itoa(medal.targetID))
	}*/
	hmobileEntry(k, btoml, roomID, upID)
}
func hmobileEntry(k int, btoml *other.Btoml, roomID, upID string) {
	/*roomID := "7117440"
	upID := "146563374"
	medals := getUserMedalsPc(btoml.BiliUser[k].Cookie)
	if len(medals) != 0 {
		roomID = strconv.Itoa(medals[0].roomID)
		upID = strconv.Itoa(medals[0].targetID)
	}*/
	tuuid := randuuid()
	buvid := randBuvid()
	url := "https://live-trace.bilibili.com/xlive/data-interface/v1/heartbeat/mobileEntry"
	data := other.Sign(`access_key=` + btoml.BiliUser[k].AccessToken + `&actionKey=appkey&appkey=1d8b6e7d45233436&area_id=0&build=6050500&buvid=` + buvid + `&c_locale=zh_CN&channel=bilibili140&client_ts=` + other.StrTime() + `&device=android&heart_beat=%5B%5D&is_patch=0&mobi_app=android&parent_id=0&platform=android&room_id=` + roomID + `&s_locale=zh_CN&seq_id=0&statistics=%7B%22appId%22%3A1%2C%22platform%22%3A3%2C%22version%22%3A%226.5.0%22%2C%22abtest%22%3A%22%22%7D&ts=` + other.StrTime() + `&uuid=` + tuuid)
	body, c := other.AppReqPost(url, data)
	if !c {
		return
	}
	//fmt.Println(0, string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在直播观看0时返回了非json字符串\n", other.TI(), k)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v hmobileEntry get message err: %v\n", k, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v hmobileEntry get code err: %v\n", k, err)
		return
	}
	if code != 0 {
		fmt.Printf("%v用户%v 在直播观看0时: %v\n", other.TI(), k, message)
		return
	}
	secret, err := js.Get("data").Get("secret_key").String()
	if err != nil {
		fmt.Printf("user %v hmobileEntry get data.secret_key err: %v\n", k, err)
		return
	}
	heart, err := js.Get("data").Get("heartbeat_interval").Int()
	if err != nil {
		fmt.Printf("user %v hmobileEntry get data.heartbeat_interval err: %v\n", k, err)
		return
	}
	timestamp, err := js.Get("data").Get("timestamp").Int()
	if err != nil {
		fmt.Printf("user %v hmobileEntry get data.timestamp err: %v\n", k, err)
		return
	}
	rule := []int{}
	ru, err := js.Get("data").Get("secret_rule").Array()
	if err != nil {
		fmt.Printf("user %v hmobileEntry get data.secret_rule err: %v\n", k, err)
		return
	}
	for i := range ru {
		ruid, err := js.Get("data").Get("secret_rule").GetIndex(i).Int()
		if err != nil {
			fmt.Printf("user %v hmobileEntry get data.secret_rule.%v err: %v\n", k, i, err)
			return
		}
		rule = append(rule, ruid)
	}
	guid := randStr(43)
	visitid := randStr(32)
	//fmt.Println(guid, visitid)
	hmobileHeartBeat(heart, timestamp, rule, secret, k, btoml, 1, tuuid, guid, visitid, buvid, roomID, upID)
}
func randStr(l int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	result := []byte{}
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func randBuvid() string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	for i := 0; i < 36; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return "X" + string(result)
}
func hmobileHeartBeat(hearttime, timestamp int, rule []int, secret string, k int, btoml *other.Btoml, sepid int, tuuid, guid, visitid, buvid, roomID, upID string) {
	time.Sleep(time.Second * time.Duration(hearttime))
	url := "https://live-trace.bilibili.com/xlive/data-interface/v1/heartbeat/mobileHeartBeat"
	ts := other.StrTime()
	playurl := ""
	upSession := ""
	clientSign := `{"platform":"android","uuid":"` + tuuid + `","buvid":"` + buvid + `","seq_id":"` + strconv.Itoa(sepid) + `","room_id":"` + roomID + `","parent_id":"1","area_id":"34","timestamp":"` + strconv.Itoa(timestamp) + `","secret_key":"` + secret + `","watch_time":"` + strconv.Itoa(hearttime) + `","up_id":"` + upID + `","up_level":"0","jump_from":"27007","gu_id":"` + guid + `","play_type":"0","play_url":"` + playurl + `","s_time":"0","data_behavior_id":"","data_source_id":"","up_session":"` + upSession + `","visit_id":"` + visitid + `","watch_status":"","click_id":"","session_id":"-99998","player_type":"0","client_ts":"` + ts + `"}`
	//fmt.Println(clientSign)
	for _, v := range rule {
		switch v {
		case 0:
			clientSign = hsign0(clientSign)
		case 1:
			clientSign = hsign1(clientSign)
		case 2:
			clientSign = hsign2(clientSign)
		case 3:
			clientSign = hsign3(clientSign)
		case 4:
			clientSign = hsign4(clientSign)
		case 5:
			clientSign = hsign5(clientSign)
		case 6:
			clientSign = hsign6(clientSign)
		case 7:
			clientSign = hsign7(clientSign)
		case 8:
			clientSign = hsign8(clientSign)
		case 9:
			clientSign = hsign9(clientSign)
		case 10:
			clientSign = hsign10(clientSign)
		case 11:
			clientSign = hsign11(clientSign)
		}
	}
	data := other.Sign(`access_key=` + btoml.BiliUser[k].AccessToken + `&actionKey=appkey&appkey=1d8b6e7d45233436&area_id=34&build=6050500&buvid=` + buvid + `&c_locale=zh_CN&channel=bilibili140&click_id=&client_sign=` + clientSign + `&client_ts=` + ts + `&data_behavior_id=&data_source_id=&device=android&gu_id=` + guid + `&jump_from=27007&mobi_app=android&parent_id=1&platform=android&play_type=0&play_url=&player_type=0&room_id=` + roomID + `&s_locale=zh_CN&s_time=0&secret_key=` + secret + `&seq_id=` + strconv.Itoa(sepid) + `&session_id=-99998&statistics={"appId":1,"platform":3,"version":"6.5.0","abtest":""}&timestamp=` + strconv.Itoa(timestamp) + `&ts=` + ts + `&up_id=` + upID + `&up_level=0&up_session=&uuid=` + tuuid + `&visit_id=` + visitid + `&watch_status=&watch_time=` + strconv.Itoa(hearttime))
	//fmt.Println(data)
	body, c := other.AppReqPost(url, data)
	if !c {
		return
	}
	//fmt.Println(sepid, string(body))
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Printf("%v用户%v 在直播观看时返回了非json字符串\n", other.TI(), k)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("user %v hmobileHeartBeat get message err: %v\n", k, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("user %v hmobileHeartBeat get code err: %v\n", k, err)
		return
	}
	if code != 0 && code != 1012003 {
		fmt.Printf("%v用户%v 在直播观看时: %v\n", other.TI(), k, message)
		//fmt.Println(sepid, string(body))
		//return
	}

	secret2, err := js.Get("data").Get("secret_key").String()
	if err != nil {
		fmt.Printf("user %v hmobileHeartBeat get data.secret_key err: %v\n", k, err)
		return
	}
	heart, err := js.Get("data").Get("heartbeat_interval").Int()
	if err != nil {
		fmt.Printf("user %v hmobileHeartBeat get data.heartbeat_interval err: %v\n", k, err)
		return
	}
	timestamp2, err := js.Get("data").Get("timestamp").Int()
	if err != nil {
		fmt.Printf("user %v hmobileHeartBeat get data.timestamp err: %v\n", k, err)
		return
	}
	rule2 := []int{}
	ru, err := js.Get("data").Get("secret_rule").Array()
	if err != nil {
		fmt.Printf("user %v hmobileHeartBeat get data.secret_rule err: %v\n", k, err)
		return
	}
	for i := range ru {
		ruid, err := js.Get("data").Get("secret_rule").GetIndex(i).Int()
		if err != nil {
			fmt.Printf("user %v hmobileHeartBeat get data.secret_rule.%v err: %v\n", k, i, err)
			return
		}
		rule2 = append(rule2, ruid)
	}
	hmobileHeartBeat(heart, timestamp2, rule2, secret2, k, btoml, sepid+1, tuuid, guid, visitid, buvid, roomID, upID)
}

func hsign0(s string) string {
	//0号算法 SHA224
	a := sha256.New224()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign1(s string) string {
	//1号算法 SHA256
	a := sha256.New()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign2(s string) string {
	//2号算法 SHA384
	a := sha512.New384()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign3(s string) string {
	//3号算法 SHA512
	a := sha512.New()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign4(s string) string {
	//4号算法 SHA3-224
	a := keccak.NewSHA3224()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign5(s string) string {
	//5号算法 SHA3-256
	a := keccak.NewSHA3256()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign6(s string) string {
	//6号算法 SHA3-384
	a := keccak.NewSHA3384()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign7(s string) string {
	//7号算法 SHA3-512
	a := keccak.NewSHA3512()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign8(s string) string {
	//8号算法 BLAKE2b512
	a, _ := blake2b.New512([]byte{})
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign9(s string) string {
	//9号算法 BLAKE2s_256
	a, _ := blake2s.New256([]byte{})
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign10(s string) string {
	//10号算法 ripemd160
	a := ripemd160.New()
	a.Write([]byte(s))
	return hex.EncodeToString(a.Sum(nil))
}
func hsign11(s string) string {
	// 11号算法 whirlpool
	w := whirlpool.New()
	w.Write([]byte(s))
	return hex.EncodeToString(w.Sum(nil))
}
