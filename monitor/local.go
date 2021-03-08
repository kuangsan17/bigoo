package monitor

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/atwat/bigoo/other"

	"github.com/bitly/go-simplejson"
	"github.com/valyala/fastjson"
	//"github.com/valyala/fastjson"
)

type biDanmu struct {
	conn   net.Conn
	roomid int
	live   bool
}

func (d *biDanmu) kill() {
	d.live = false
	d.conn.Close()
}
func (d *biDanmu) restart() {
	d.kill()
	time.Sleep(time.Second * 2)
	danmu(d.roomid)
}

func getroomidsfromallarea() (roomids [][]int) {
	roomids = [][]int{}
	for i := 1; i <= 6; i++ {
		roomidsarea := getroomidsfromarea(i)
		roomidsarea2 := []int{}
		for _, roomid := range roomidsarea {
			if len(roomidsarea2) >= 2 {
				break
			}
			if isliving(roomid) {
				roomidsarea2 = append(roomidsarea2, roomid)
			}
		}
		roomids = append(roomids, roomidsarea2)
	}
	return
}

func getroomidsfromarea(area int) (roomids []int) {
	roomids = []int{}
	b, c := other.PcReqGet("https://api.live.bilibili.com/room/v3/area/getRoomList?platform=web&cate_id=0&area_id=0&sort_type=online&page=1&page_size=20&tag_version=1&parent_area_id="+strconv.Itoa(area), "")
	if !c {
		return
	}

	js, err := simplejson.NewJson(b)
	if err != nil {
		return
	}
	list, err := js.Get("data").Get("list").Array()
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		roomid, err := js.Get("data").Get("list").GetIndex(i).Get("roomid").Int()
		if err != nil {
			return
		}
		roomids = append(roomids, roomid)
	}
	return
}

func isliving(roomid int) bool {
	b, c := other.PcReqGet("https://api.live.bilibili.com/room/v1/Room/get_info?room_id="+strconv.Itoa(roomid), "")
	if c != true {
		return false
	}
	js, err := simplejson.NewJson(b)
	if err != nil {
		return false
	}
	bb, err := js.Get("data").Get("live_status").Int()
	if err == nil && bb == 1 {
		return true
	}
	return false
}

type localRoom struct {
	*sync.RWMutex
	localRoomids [][]int
}

var localRoomids *localRoom

func reroomid(roomid int) {
	localRoomids.Lock()
	defer localRoomids.Unlock()
	roomids := localRoomids.localRoomids
	for k1, v1 := range roomids {
		for k2, v2 := range v1 {
			if v2 == roomid {
				area := k1 + 1
				roomids := getroomidsfromarea(area)
				for _, v3 := range roomids {
					if v3 != roomid && (v3 != v1[1-k2] && isliving(v3)) {
						danmu(v3)
						return
					}
				}
			}
		}
	}
}

var localCh chan int

func local(ch chan int, btoml *other.Btoml) {
	if !btoml.Monitor.Local {
		return
	}
	localCh = ch
	roomids := getroomidsfromallarea()
	localRoomids = &localRoom{&sync.RWMutex{}, roomids}
	for _, v := range roomids {
		for _, vv := range v {
			go danmu(vv)
		}
	}
}

func danmu(roomid int) {
	token, host, port := getDanmuInfo(roomid)
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(other.TI()+"dial failed, err", err)
		return
	}
	bdanmu := &biDanmu{conn, roomid, true}
	//defer conn.Close()
	go read(bdanmu)
	msg := string([]byte{0, 16, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0}) + `{"group":"","uid":0,"roomid":` + strconv.Itoa(roomid) + `,"key":"` + token + `","platform":"android","clientver":"5.52.1.5521100","hwid":"OQw5XDgLbg1oXGxZJVkl","protover":2}`
	data, err := Encode(msg)
	if err != nil {
		fmt.Println(other.TI()+"encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("握手err", err)
		bdanmu.restart()
	}
	go func() {
		for {
			time.Sleep(time.Second * 30)
			msg := string([]byte{0, 16, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0}) + `{}`
			data, err := Encode(msg)
			if err != nil {
				fmt.Println(other.TI()+"encode msg failed, err:", err)
				return
			}
			if !bdanmu.live {
				break
			}
			_, err = conn.Write(data)
			if err != nil {
				fmt.Println("心跳err", err)
				bdanmu.restart()
			}
		}
	}()
}

func read(bdanmu *biDanmu) {
	oldbyte := []byte{}
	for bdanmu.live {
		buff := [2560]byte{}
		n, err := bdanmu.conn.Read(buff[:])
		if !bdanmu.live {
			break
		}
		if err != nil {
			fmt.Println(other.TI()+"recv failed, err:", err)
			bdanmu.restart()
			return
		}
		//fmt.Println(string(buff[:n]))
		buf := append(oldbyte, buff[:n]...)
		//fmt.Println("buf:", string(buf))
		for {
			if len(buf) < 4 {
				oldbyte = buf
				break
			}
			alllength := int(buf[3]) + int(buf[2])*256 + int(buf[1])*256*256 + int(buf[0])*256*256*256
			//	fmt.Println("长度:", alllength)
			if len(buf) < alllength {
				//fmt.Println("半包:", string(buf))
				oldbyte = buf
				break
			} else {
				//fmt.Println("------", buf[11])
				//fmt.Println("全包:", string(buf[16:alllength]))
				switch buf[11] {
				case 3:
					//fmt.Println(other.TI() + "心跳成功")
				case 5:
					//go danmumsg(buf[16:alllength])
					gzipmsg(buf[16:alllength], bdanmu)
				case 8:
					fmt.Print(other.TI()+"房间", bdanmu.roomid, " 弹幕连接成功\n")
				default:
					fmt.Println(other.TI()+"----其他----", buf[11])
				}
				if len(buf) == alllength {
					oldbyte = []byte{}
					break
				} else {
					buf = buf[alllength:]
				}
			}
		}
	}
}

//DoZlibUnCompress 进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) ([]byte, error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return compressSrc, err
	}
	io.Copy(&out, r)
	return out.Bytes(), nil
}

func gzipmsg(msg []byte, bdanmu *biDanmu) []byte {
	buf, err := DoZlibUnCompress(msg)
	if err != nil {
		danmumsg(buf, bdanmu)
		return []byte{}
	}
	oldbyte := []byte{}
	for {
		if len(buf) < 4 {
			oldbyte = buf
			break
		}
		alllength := int(buf[3]) + int(buf[2])*256 + int(buf[1])*256*256 + int(buf[0])*256*256*256
		//	fmt.Println("长度:", alllength)
		if len(buf) < alllength {
			//fmt.Println("半包:", string(buf))
			oldbyte = buf
			break
		} else {
			//fmt.Println("------", buf[11])
			//fmt.Println("全包:", string(buf[16:alllength]))
			switch buf[11] {
			case 3:
				fmt.Println(other.TI() + "心跳成功w")
			case 5:
				go danmumsg(buf[16:alllength], bdanmu)
			case 8:
				fmt.Println(other.TI() + "弹幕连接成功w")
			default:
				fmt.Println(other.TI()+"----其他----w", buf[11])
			}
			if len(buf) == alllength {
				oldbyte = []byte{}
				break
			} else {
				buf = buf[alllength:]
			}
		}
	}
	if len(oldbyte) != 0 {
		fmt.Println("====>", oldbyte)
	}
	return oldbyte
}

func danmumsg(msg []byte, bdanmu *biDanmu) {
	//return
	cmd := fastjson.GetString(msg, "cmd")
	switch cmd {
	case "NOTICE_MSG": //msg_common msg_type
		msgType := fastjson.GetInt(msg, "msg_type")
		realRoomid := fastjson.GetInt(msg, "real_roomid")
		msgCommon := fastjson.GetString(msg, "msg_common")
		if msgType == 3 /*舰队*/ || msgType == 2 /*小电视摩天搂等*/ || msgType == 8 /*任意门*/ {
			//roomidch <- realRoomid
			if msgCommon != "" {
				//	fmt.Println(realRoomid, msgCommon)
				localCh <- realRoomid
			}
		}
		//fmt.Println(string(msg))
	case "GUARD_BUY": //data{username guard_level num price gift_name start_time end_time}
		//fmt.Println("====", string(msg), "====")
	case "SPECIAL_GIFT":
	case "PREPARING": //下播
		r := bdanmu.roomid
		fmt.Printf("%v房间%v 主播下播\n", other.TI(), r)
		bdanmu.kill()
		reroomid(r)
	case "ANCHOR_LOT_start": //禁言 data {uname operator}

	default:
		//fmt.Println(string(msg))
	}
}

// Encode 将消息编码
func Encode(message string) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.BigEndian, 4+length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.BigEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}
func getDanmuInfo(roomid int) (token string, host string, port int) {
	url := "https://api.live.bilibili.com/xlive/app-room/v1/index/getDanmuInfo"
	data := `actionKey=appkey&appkey=1d8b6e7d45233436&build=5521100&channel=huawei&device=android&mobi_app=android&platform=android&room_id=` + strconv.Itoa(roomid) + `&statistics=%7B%22appId%22%3A1%2C%22platform%22%3A3%2C%22version%22%3A%225.52.1%22%2C%22abtest%22%3A%22%22%7D&ts=` + other.StrTime()
	data = other.Sign(data)
	r, c := other.AppReqGet(url + "?" + data)
	if !c {
		return
	}
	token = fastjson.GetString(r, "data", "token")
	host = fastjson.GetString(r, "data", "ip_list", "0", "host")
	port = fastjson.GetInt(r, "data", "ip_list", "0", "port")
	return
}
