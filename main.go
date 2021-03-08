package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/atwat/bigoo/job"
	"github.com/atwat/bigoo/monitor"
	"github.com/atwat/bigoo/other"
)

var (
	user1    string
	user2    string
	user3    string
	pass1    string
	pass2    string
	pass3    string
	bilive   string
	yjaddr   string
	yjkey    string
	userPass [][]string = [][]string{}
	help     bool
)

func init() {
	flag.StringVar(&user1, "user1", "", "用户1")
	flag.StringVar(&user2, "user2", "", "用户2")
	flag.StringVar(&user3, "user3", "", "用户3")
	flag.StringVar(&pass1, "pass1", "", "密码1")
	flag.StringVar(&pass2, "pass2", "", "密码2")
	flag.StringVar(&pass3, "pass3", "", "密码3")
	flag.StringVar(&bilive, "bilive", "", "bilive监控")
	flag.StringVar(&yjaddr, "yjaddr", "", "yj监控addr")
	flag.StringVar(&yjkey, "yjkey", "", "yj监控key")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.Usage = usage
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if user1 != "" && pass1 != "" {
		userPass = append(userPass, []string{user1, pass1})
	}
	if user2 != "" && pass2 != "" {
		userPass = append(userPass, []string{user2, pass2})
	}
	if user3 != "" && pass3 != "" {
		userPass = append(userPass, []string{user3, pass3})
	}
	//fmt.Println(userPass)
}
func usage() {
	flag.PrintDefaults()
}
func main() {
	fmt.Println(`
           |    |
     ------------------
    |                  |
    |     |      |     |
    |                  |
    |仅供学习交流，请在|
    |下载后30分钟内删除|
     ------------------
`)
	t := other.TomlStart(userPass, bilive, yjaddr, yjkey)
	Status := other.StatusStart(t)
	go func() {
		for {
			time.Sleep(time.Hour * 12)
			Status.Clear2Day()
		}
	}()
	job.AllJobs(t)

	ch := make(chan int)
	go monitor.Monitor(ch, t)
	go func() {
		for {
			r := <-ch
			//fmt.Println(r, "+++++")
			go job.EnterRoom(r, t, Status, false)
		}
	}()
	for {
		line := 0
		fmt.Scanln(&line)
		switch line {
		case 0:
			//fmt.Println(Status)

		case 1:
			yesterdaygifts := Status.YesterTodayGifts()
			if len(yesterdaygifts) != 0 {
				fmt.Println("昨日礼物：")
				for k, v := range yesterdaygifts {
					fmt.Printf("  %v X %v\n", k, v)
				}
			}

			fmt.Println("今日礼物：")
			daygifts := Status.TodayGifts()
			if len(daygifts) == 0 {
				fmt.Println("  nil")
			}
			for k, v := range daygifts {
				fmt.Printf("  %v X %v\n", k, v)
			}
		case 2:
			ytmp := Status.YesterTodayUserAward()
			if len(ytmp) != 0 {
				fmt.Println("昨日获奖：")
				//fmt.Println(ytmp)
			}
			for userid := range ytmp {
				fmt.Printf("  用户%v:\n", userid)
				for k, v := range ytmp[userid] {
					fmt.Printf("    %v X %v\n", k, v)
				}
				if len(ytmp[userid]) == 0 {
					fmt.Println("    nil")
				}
			}

			fmt.Println("今日获奖：")
			tmp := Status.TodayUserAward()
			//fmt.Println(tmp)
			for userid := range tmp {
				fmt.Printf("  用户%v:\n", userid)
				for k, v := range tmp[userid] {
					fmt.Printf("    %v X %v\n", k, v)
				}
				if len(tmp[userid]) == 0 {
					fmt.Println("    nil")
				}
			}
			if len(tmp) == 0 {
				fmt.Println("  nil")
			}
		case 3:
			job.PrintUserBags(t)
		case 4:
			job.PrintUserMedals(t)
		case 5:
			for userid := range t.BiliUser {
				fmt.Printf("用户%v 账号: %v\n", userid, t.BiliUser[userid].UserName)
			}
		default:
			if line > 1000 {
				//ch <- line
				go job.EnterRoom(line, t, Status, true)
			} else {
				fmt.Println(">>", line)
			}
		}
	}
}
