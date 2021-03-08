package other

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	rand2 "math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

var buvid string = "XY1093E13D66910EE046E31B3859313C6DA90"
var deviceid string = ""
var fingerprint string = ""

//Login ...
func Login(uname, upass, ca string) (access, refresh, cookie string, boo bool) {
	buvid = randStr(36)
	u, p := hashPass(uname, upass)
	url := "https://passport.bilibili.com/x/passport-tv-login/login"
	data := TvSign("appkey=4409e2ce8ffd12b8&bili_local_id=" + deviceid + "&build=102502&buvid=" + buvid + "&channel=dangbei&code=&device=HONOR&device_id=" + deviceid + "&device_name=HWYAL-O&device_platform=Android10HUAWEIYAL-AL50&fingerprint=" + fingerprint + "&guid=" + buvid + "&local_fingerprint=" + fingerprint + "&local_id=" + buvid + "&mobi_app=android_tv_yst&networkstate=&password=" + p + "&platform=android&token=c08d1927e69cac0e&ts=" + StrTime() + "&username=" + u)
	body, a := AppReqPost(url+"?"+data, "")
	if !a {
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println("登录时出错:", err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		return
	}
	if message == "CAPTCHA is not match" {
		ca1 := captcha()
		return Login(uname, upass, ca1)
	}

	if code != 0 {
		fmt.Println(TI() + message)
		return
	}
	access, err = js.Get("data").Get("access_token").String()
	if err != nil {
		fmt.Println("登录时出错:", err)
		return
	}
	refresh, err = js.Get("data").Get("refresh_token").String()
	if err != nil {
		fmt.Println("登录时出错:", err)
		return
	}
	access, refresh, cookie, boo = Refresh(access, refresh)
	//fmt.Println(TI() + "登录成功          ")
	return
}

//Iflogin return the time
func Iflogin(access, cookie string) (expiresin int) {
	buvid = randStr(36)
	url := "https://passport.bilibili.com/x/passport-login/oauth2/info"

	data := TvSign("access_key=" + access + "&appkey=4409e2ce8ffd12b8&bili_local_id=" + deviceid + "&build=102502&buvid=" + buvid + "&channel=dangbei&device=HONOR&device_id=" + deviceid + "&device_name=HWYAL-O&device_platform=Android10HUAWEIYAL-AL50&fingerprint=" + fingerprint + "&guid=" + buvid + "&local_fingerprint=" + fingerprint + "&local_id=" + buvid + "&mobi_app=android_tv_yst&networkstate=&platform=android&ts=" + StrTime())
	b, _ := AppReqGet(url + "?" + data)
	//fmt.Println(string(b))
	js, err := simplejson.NewJson(b)
	if err != nil {
		return
	}
	expiresin, _ = js.Get("data").Get("expires_in").Int()
	return
}

//Refresh refresh the cookie
func Refresh(access, refresh string) (newaccess, newrefresh, newcookie string, boo bool) {
	data := TvSign("access_token=" + access + "&appkey=4409e2ce8ffd12b8&build=102502&channel=dangbei&guid=" + buvid + "&mobi_app=android_tv_yst&platform=android&refresh_token=" + refresh + "&ts=" + StrTime())
	b, _ := AppReqPost("https://passport.snm0516.aisee.tv/api/v2/oauth2/refresh_token?"+data, "")
	js, err := simplejson.NewJson(b)
	if err != nil {
		return
	}
	newaccess, err = js.Get("data").Get("token_info").Get("access_token").String()
	if err != nil {
		return
	}
	newrefresh, err = js.Get("data").Get("token_info").Get("refresh_token").String()
	if err != nil {
		return
	}
	cookies, err := js.Get("data").Get("cookie_info").Get("cookies").Array()
	if err != nil {
		return
	}
	for k := range cookies {
		cookie1, err := js.Get("data").Get("cookie_info").Get("cookies").GetIndex(k).Get("name").String()
		if err != nil {
			return
		}
		cookie2, err := js.Get("data").Get("cookie_info").Get("cookies").GetIndex(k).Get("value").String()
		if err != nil {
			return
		}
		newcookie = newcookie + cookie1 + "=" + cookie2 + ";"
	}
	boo = true
	return
}
func randStr(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return "X" + string(result)
}
func hashPass(user, pass string) (username, b string) {
	buvid = randStr(36)
	url := "https://passport.bilibili.com/x/passport-tv-login/key?" + TvSign("appkey=4409e2ce8ffd12b8&bili_local_id="+deviceid+"&build=102502&buvid="+buvid+"&channel=dangbei&device=HONOR&device_id="+deviceid+"&device_name=HWYAL-O&device_platform=Android10HUAWEIYAL-AL50&fingerprint="+fingerprint+"&guid="+buvid+"&local_fingerprint="+fingerprint+"&local_id="+buvid+"&mobi_app=android_tv_yst&networkstate=&platform=android&ts="+StrTime())
	body, status := TvReqGet(url)
	if !status {
		time.Sleep(time.Second * 1200)
		return hashPass(user, pass)
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println("登录时getKey出错:", err)
		return
	}
	key, err := js.Get("data").Get("key").String()
	if err != nil {
		fmt.Println("登录时getKey出错:", err)
		return
	}
	hash, err := js.Get("data").Get("hash").String()
	if err != nil {
		fmt.Println("登录时getKey出错:", err)
		return
	}
	username = strings.Replace(user, "@", "%40", -1)
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		fmt.Println("登录时public key error")
		return user, pass
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("登录时rsa出错", err)
		return user, pass
	}
	pub := pubInterface.(*rsa.PublicKey)
	b1, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(hash+pass))
	if err != nil {
		fmt.Println("登录时rsa出错", err)
		return user, pass
	}
	b = base64.StdEncoding.EncodeToString(b1)
	b = strings.Replace(b, "%", "%25", -1)
	b = strings.Replace(b, "+", "%2B", -1)
	b = strings.Replace(b, " ", "%20", -1)
	b = strings.Replace(b, "/", "%2F", -1)
	b = strings.Replace(b, "?", "%3F", -1)
	b = strings.Replace(b, "#", "%23", -1)
	b = strings.Replace(b, "&", "%26", -1)
	b = strings.Replace(b, "=", "%3D", -1)
	b = strings.Replace(b, "@", "%40", -1)
	//fmt.Println(username, b)
	return username, b
}

func captcha() string {
	url := "http://passport.snm0516.aisee.tv/api/captcha?token=8d526ed57efbbc2a"
	r, status := AppReqGet(url)
	if !status {
		time.Sleep(time.Second * 1200)
		return captcha()
	}

	f, err := os.Create("验证码.png")
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		f.Write(r)
	}

	encodeString := base64.StdEncoding.EncodeToString(r)
	data1 := `{"image": "` + encodeString + `"}`
	req, err := http.NewRequest("POST", "http://152.32.186.69:19951/captcha/v1", bytes.NewBuffer([]byte(data1)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(TI()+"识别验证码出错", err)
		f, err := os.Create("验证码.png")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			f.Write(r)
		}
		f.Close()
		ca := ""
		fmt.Print(TI() + "请手动输入验证码:")
		fmt.Scanln(&ca)
		return ca
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println(TI()+"识别验证码出错", err)
		f, err := os.Create("验证码.png")
		defer f.Close()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			f.Write(r)
		}
		ca := ""
		fmt.Print(TI() + "请手动输入验证码:")
		fmt.Scanln(&ca)
		return ca
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Println(TI()+"识别验证码出错", err)
		f, err := os.Create("验证码.png")
		defer f.Close()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			f.Write(r)
		}
		ca := ""
		fmt.Print(TI() + "请手动输入验证码:")
		fmt.Scanln(&ca)
		return ca
	}
	fmt.Println(TI() + "识别验证码:" + message)
	return message
}
