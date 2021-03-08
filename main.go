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
	help     bool
	tomljson string
)

func init() {
	flag.StringVar(&tomljson, "toml", "", "配置")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.Usage = usage
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

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
	t := other.TomlStart(tomljson)
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
