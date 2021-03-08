package monitor

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/atwat/bigoo/other"

	"github.com/bitly/go-simplejson"
)

var yjch chan int
var yjmainch chan int
var tomlc *other.Btoml

func yj(ch1, mainch chan int, Tomlc1 *other.Btoml) {
	jkroomids = other.NewSafeDict(map[string]int{})
	yjch = ch1
	yjmainch = mainch
	tomlc = Tomlc1
	if tomlc.Monitor.YjKey == "" {
		return
	}
	thetime := other.IntTime()
	jkroomids.Put(tomlc.Monitor.YjKey, thetime)
	yjDanmu(tomlc.Monitor.YjKey, thetime)

}

type yjLive struct {
	conn      *net.TCPConn
	lastbyte  []byte
	roomid    string
	live      bool
	starttime int
}

var jkroomids *other.SafeDict //使用中
func yjDanmu(roomid string, t int) {
	jk := yjLive{roomid: roomid, live: true, starttime: t}
	jk.hello(roomid)
	go jk.getmessage()
	go jk.heart()
}

func (c *yjLive) restart() {
	if c.live == false {
		return
	}
	c.kill()
	time.Sleep(time.Second)
	fmt.Print(other.TI() + "重新连接    ")
	thetime := other.IntTime()
	jkroomids.Put(c.roomid, thetime)
	yjDanmu(c.roomid, thetime)
}

func (c *yjLive) kill() {
	c.live = false
	c.conn.Close()
}
func (c *yjLive) hello(roomid string) {
	tcpaddr, err := net.ResolveTCPAddr("tcp4", tomlc.Monitor.YjAddr)
	c.conn, err = net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		//fmt.Println("KEY", c.roomid, "还能有这种错误？", err)
		r := c.roomid
		fmt.Print(other.TI() + "连接时错误    ")
		time.Sleep(time.Second)
		//fmt.Println("hello重新连接", r)
		go yjDanmu(r, c.starttime)
		runtime.Goexit()
	}
	data, err := encode(`{"code": 0,"type":"ask","data": {"key": "` + roomid + `"}}`)
	if err != nil {
		fmt.Println(other.TI(), err)
		return
	}
	_, err = c.conn.Write(data)
	if err != nil {
		fmt.Print(other.TI() + "握手时错误")
		c.restart()
		return
	}
}

func (c *yjLive) heart() {
	for {
		if c.live == false {
			return
		}
		ttt, _ := jkroomids.Get(c.roomid)
		if c.starttime != ttt {
			fmt.Print(other.TI() + "关闭     ")
			c.kill()
			return
		}
		data, err := encode("")
		if err != nil {
			fmt.Println(other.TI(), err)
			return
		}
		_, err = c.conn.Write(data)
		if err != nil {
			fmt.Print(other.TI()+"心跳失败", "    ")
			c.restart()
			return
		}
		time.Sleep(30 * time.Second)
		//fmt.Println("心跳成功")
	}
}

func (c *yjLive) getmessage() {
	for c.live {
		read := bufio.NewReader(c.conn)
		s, err := decode(read)
		if err == nil && len(s) != 0 {
			c.mesgtype(s)
		}
	}
}

func (c *yjLive) mesgtype(msg []byte) {
	if len(msg) < 20 {
		return
	}
	//fmt.Println(string(msg))
	js, err := simplejson.NewJson(msg)
	if err != nil {
		fmt.Println(other.TI(), err)
		return
	}
	thetype, err := js.Get("type").String()
	if err != nil {
		return
	}
	if thetype == "raffle" {
		roomid, err := js.Get("data").Get("room_id").Int()
		if err != nil {
			return
		}
		raffleType, _ := js.Get("data").Get("raffle_type").String()
		if raffleType == "STORM" {
			yjmainch <- roomid
			return
		}
		yjch <- roomid
	} else {
		if thetype == "error" { //"entered"
			fmt.Println(other.TI()+"Yj返回", thetype)
			c.kill()
			if sleept > 30 {
				fmt.Println(other.TI() + "Yj多次尝试连接失败！请检测key")
				return
			}
			sleept = sleept * 2
			time.Sleep(time.Second * time.Duration(sleept))
			thetime := other.IntTime()
			jkroomids.Put(c.roomid, thetime)
			go yjDanmu(c.roomid, thetime)
			return
		}
		fmt.Println(other.TI() + "Yj已连接")
	}
}

var sleept int = 1

func byte2num(a []byte) int {
	l := len(a)
	num := 0
	for i := l; i > 0; i-- {
		num = num + int(a[i-1])*int(math.Pow(float64(256), float64(l-i)))
	}
	return num
}
func generatePacket(a int, body string) []byte {
	b := []byte{0, 0, 0, byte(len(body) + 16), 0, 16, 0, 1, 0, 0, 0, byte(a), 0, 0, 0, 1}
	return append(b, []byte(body)...)
}

func inttime() int {
	tt := time.Now().Unix()
	strInt64 := strconv.FormatInt(tt, 10)
	id16, _ := strconv.Atoi(strInt64)
	return id16
}

func encode(message string) ([]byte, error) {
	// 读取消息的长度
	var length int32 = int32(len(message))
	var pkg *bytes.Buffer = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.BigEndian, length)
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

func decode(reader *bufio.Reader) ([]byte, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.BigEndian, &length)
	if err != nil {
		return []byte{}, err
	}
	//fmt.Println(int32(reader.Buffered()), "?", length+4)
	if int32(reader.Buffered()) < length+4 {
		return []byte{}, err
	}

	// 读取消息真正的内容
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return []byte{}, err
	}
	return pack[4:], nil
}
