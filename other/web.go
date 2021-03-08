package other

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var appheader [][2]string = [][2]string{
	{"Display-ID", "146563374-1583633292"},
	{"Buvid", "XZF6581FC00A30E3D62581F6644566F1AA4D0"},
	{"Device-ID", "LhwvGy4YeRgsT3cRbRFtA2ceeg8-Bz4OOgo6Cz8JMA"},
	{"env", "prod"},
	{"APP-KEY", "android"},
	{"User-Agent", "Mozilla/5.0 BiliDroid/5.54.0 (bbcallen@gmail.com)"},
	{"Connection", "keep-alive"},
}

var blueheader [][2]string = [][2]string{
	{"accept-encoding", "gzip"},
	{"cookie", "sid=a4v7ywe1"},
	{"user-agent", "Mozilla/5.0 BiliDroid/1.9.32 (bbcallen@gmail.com)"},
	//{"Connection", "keep-alive"},
}

var pcheader [][2]string = [][2]string{
	//{"Accept", "application/json, text/plain, */*"},
	//{"Accept-Encoding", "gzip, deflate, br"},
	{"Accept-Language", "zh-CN,zh;q=0.9"},
	{"Connection", "keep-alive"},
	{"Content-Type", "application/x-www-form-urlencoded; charset=utf-8"},
	{"User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.116 Safari/537.36"},
}

var tvheader [][2]string = [][2]string{
	{"Buvid", "XZc9b56c63d7c44e99476d2775cb4c968bb8c"},
	{"User-Agent", "Mozilla/5.0 BiliTV/1.2.3.1 (bbcallen@gmail.com)"},
	{"env", "prod"},
	{"APP-KEY", "android_tv_yst"},
	{"Content-Type", "application/x-www-form-urlencoded; charset=utf-8"},
	{"Connection", "keep-alive"},
	{"Accept-Encoding", "gzip"},
}

var netErr map[string]int = map[string]int{}
var netErrLock sync.RWMutex

//HTTPreq is http requests
func HTTPreq(method, url, data string, header ...[2]string) ([]byte, bool) {
	ur := url
	a := strings.Index(url, "?")
	if a != -1 {
		ur = ur[:a]
	}
	//fmt.Println(TI() + "request: " + ur)
	netErrLock.Lock()
	if _, ok := netErr[ur]; !ok {
		//fmt.Println(ur)
		netErr[ur] = 0
	}
	netError := netErr[ur]
	netErrLock.Unlock()

	if IntTime()-netError < 601 {
		return []byte{}, false
	}
	client := http.Client{
		Timeout: time.Duration(time.Second * 10),
	}
	for i := 1; i <= 10; i++ {
		var request *http.Request
		var err error
		if data == "" {
			request, err = http.NewRequest(method, url, nil)
		} else {
			request, err = http.NewRequest(method, url, strings.NewReader(data))
		}
		if err != nil {
			time.Sleep(time.Second / 5)
			continue
		}
		if method == "POST" {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		for _, v := range header {
			if v[0] == "" || v[1] == "" {
				continue
			}
			request.Header.Set(v[0], v[1])
		}
		//fmt.Println(request.Header)
		response, err := client.Do(request)
		if err != nil {
			txt := err.Error()
			reg := regexp.MustCompile(`access_key=.*?&`)
			txt = reg.ReplaceAllString(txt, "access_key=***&")
			reg = regexp.MustCompile(`sign=.*?"`)
			txt = reg.ReplaceAllString(txt, "sign=***\"")
			if false {
				fmt.Print("\r", txt, "\n\r正在重试中-")
			}
			time.Sleep(time.Second / 5)
			continue
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			txt := err.Error()
			reg := regexp.MustCompile(`access_key=.*?&`)
			txt = reg.ReplaceAllString(txt, "access_key=***&")
			reg = regexp.MustCompile(`sign=.*?"`)
			txt = reg.ReplaceAllString(txt, "sign=***\"")
			if false {
				fmt.Print("\r", txt, "\n\r正在重试中")
			}
			time.Sleep(time.Second / 5)
			continue
		}
		if response.StatusCode == 403 || response.StatusCode == 412 {
			netErrLock.Lock()
			netErr[ur]++
			if netErr[ur] > 2 {
				fmt.Println(TI() + ur + "频繁异常")
				netErr[ur] = IntTime()
			}
			netErrLock.Unlock()
			return []byte{}, false
		}
		netErrLock.Lock()
		netErr[ur] = 0
		netErrLock.Unlock()
		return body, true

	}
	fmt.Println("重试10次还不行？不管了")
	return []byte{}, false
}

//PcReq ...
func PcReq(method, url, data string, header ...[2]string) ([]byte, bool) {
	newheader := append(pcheader, header...)
	return HTTPreq(method, url, data, newheader...)
}

//PcReqGet ...
func PcReqGet(url, cookie string) ([]byte, bool) {
	return PcReq("GET", url, "", [2]string{"Cookie", cookie})
}

//PcReqPost ...
func PcReqPost(url, data, cookie string) ([]byte, bool) {
	return PcReq("POST", url, data, [2]string{"Cookie", cookie})
}

//AppReq ...
func AppReq(method, url, data string, header ...[2]string) ([]byte, bool) {
	newheader := append(appheader, header...)
	return HTTPreq(method, url, data, newheader...)
}

//AppReqGet ...
func AppReqGet(url string) ([]byte, bool) {
	return AppReq("GET", url, "", [2]string{})
}

//AppReqPost ...
func AppReqPost(url, data string) ([]byte, bool) {
	return AppReq("POST", url, data, [2]string{})
}

//BlueReq ...
func BlueReq(method, url, data string, header ...[2]string) ([]byte, bool) {
	newheader := append(blueheader, header...)
	return HTTPreq(method, url, data, newheader...)
}

//BlueReqGet ...
func BlueReqGet(url string) ([]byte, bool) {
	return AppReq("GET", url, "", [2]string{})
}

//BlueReqPost ...
func BlueReqPost(url, data string) ([]byte, bool) {
	return AppReq("POST", url, data, [2]string{})
}

//TvReq ...
func TvReq(method, url, data string, header ...[2]string) ([]byte, bool) {
	newheader := append(appheader, header...)
	return HTTPreq(method, url, data, newheader...)
}

//TvReqGet ...
func TvReqGet(url string) ([]byte, bool) {
	return TvReq("GET", url, "", [2]string{})
}

//TvReqPost ...
func TvReqPost(url, data string) ([]byte, bool) {
	return TvReq("POST", url, data, [2]string{})
}
