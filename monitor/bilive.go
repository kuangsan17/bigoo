package monitor

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/atwat/bigoo/other"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

func bilive(ch, mainch chan int, btoml *other.Btoml) {

	ws := btoml.Monitor.Bilive
	if ws == "" {
		return
	}
	//ws := "ws://39.108.158.17:20080/"
	dialer := websocket.Dialer{}
	p := strings.IndexAny(ws, "#")
	if p != -1 {
		dialer = websocket.Dialer{
			Subprotocols: []string{ws[p+1:]},
		}
	}
	header := map[string][]string{}
	//header["User-Agent"] = []string{"Yoki/" + createRand() + "/45"}
	header["User-Agent"] = []string{"bigoo/test/0.0.1"}
	c, _, err := dialer.Dial(ws, header)
	if err != nil {
		fmt.Println(ws, err)
		go func() {
			time.Sleep(time.Second * 5)
			bilive(ch, mainch, btoml)
		}()
		runtime.Goexit()
	}
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println(ws, err)
			go func() {
				time.Sleep(time.Second * 5)
				bilive(ch, mainch, btoml)
			}()
			runtime.Goexit()
		}
		//fmt.Print("\r", string(msg), "\n")
		def(msg, ch, mainch)
	}
}

func def(mesg []byte, ch, mainch chan int) {
	js, err := simplejson.NewJson(mesg)
	if err != nil {
		return
	}
	cmd, err := js.Get("cmd").String()
	if err != nil {
		return
	}
	if cmd == "sysmsg" {
		msg, err := js.Get("msg").String()
		if err != nil {
			return
		}
		fmt.Print(other.TI(), "来自ws的信息:", msg, "\n")
	}
	roomid, err := js.Get("roomID").Int()
	if err != nil {
		return
	}
	ty, err := js.Get("type").String()
	if err != nil {
		return
	}
	if ty == "beatStorm" {
		mainch <- roomid
		return
	}
	ch <- roomid
}

func createRand() string {
	r := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	//fmt.Println(r)
	return r
}
