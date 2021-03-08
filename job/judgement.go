package job

import (
	"fmt"
	"math/rand"
	"github.com/atwat/bigoo/other"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

func judgementMain(btoml *other.Btoml) {
	for userid := range btoml.BiliUser {
		go func(userid int) {
			for btoml.MoreSetting.Judgement {
				if other.NowHour() == 8 || other.NowHour() == 9 {
					randSleep(0, 300, (userid+1)*22)
					go judgement(userid, btoml)
					time.Sleep(time.Hour * 20)
				} else {
					time.Sleep(time.Minute * 33)
				}
			}
		}(userid)
	}
}
func refererget(url, cookie, referer string) ([]byte, bool) {
	return other.PcReq("GET", url, "", [2]string{"Cookie", cookie}, [2]string{"referer", referer})
}

func refererpost(url, data, cookie, referer string) ([]byte, bool) {
	return other.PcReq("POST", url, data, [2]string{"Cookie", cookie}, [2]string{"referer", referer})
}

type opin struct {
	content string
	vote    int
	like    int
	hate    int
}

func judgementVote(id, voteid int, userid int, btoml *other.Btoml) {
	txt := ""
	cookie := btoml.BiliUser[userid].Cookie
	aqw := strings.Index(cookie, "bili_jct=")
	if aqw == -1 {
		return
	}
	bqw := cookie[aqw+9:]
	cqw := strings.Index(bqw, ";")
	if cqw == -1 {
		return
	}
	csrf := bqw[:cqw]
	switch voteid {
	case 1:
		txt = "封禁"
	case 2:
		txt = "不违规"
	case 3:
		txt = "放弃"
	case 4:
		txt = "删除"
	}
	url := "https://api.bilibili.com/x/credit/jury/vote"
	data := fmt.Sprintf("jsonp=jsonp&cid=%v&vote=%v&content=&likes&hates=&attr=0&csrf=%v", id, voteid, csrf)
	b, c := refererpost(url, data, cookie, "https://www.bilibili.com/judgement/vote/"+strconv.Itoa(id))
	if !c {
		return
	}
	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Printf("%vuser %v judgementVote json err: %v\n", other.TI(), userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("%vuser %v judgementVote json code err: %v\n", other.TI(), userid, err)
		return
	}
	if code == 0 {
		fmt.Printf("%v用户%v 案件%v 投票%v成功\n", other.TI(), userid, id, txt)
		return
	}
	message, err := js.Get("message").String()
	if err != nil {
		fmt.Printf("%vuser %v judgementVote json message err: %v\n", other.TI(), userid, err)
		return
	}
	fmt.Printf("%v用户%v 案件%v 投票%v失败:%v\n", other.TI(), userid, id, txt, message)
}

func createRand() string {
	r := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1e16))
	return r
}
func judgementOpinion(id int, userid int, btoml *other.Btoml) (opins []opin, type1, type2 int) {
	cookie := btoml.BiliUser[userid].Cookie
	opins = []opin{}
	for otype := 1; otype <= 2; otype++ {
		tss := time.Now().UnixNano() / 1e6
		url := fmt.Sprintf("https://api.bilibili.com/x/credit/jury/vote/opinion?jsonp=jsonp&callback=jQuery1720%v_%v&cid=%v&otype=%v&pn=1&ps=10&_=%v", createRand(), tss, id, otype, tss)
		b, c := refererget(url, cookie, "https://www.bilibili.com/judgement/vote/"+strconv.Itoa(id))
		if !c {
			continue
		}
		//fmt.Println(string(b))
		lb := len(b)
		for k, v := range b {
			if v == byte(40) {
				//fmt.Println(k, b[k+1:lb-1])
				b = b[k+1 : lb-1]
				break
			}
		}
		js, err := simplejson.NewJson(b)
		if err != nil {
			fmt.Printf("%vuser %v judgementOpinion json err: %v\n", other.TI(), userid, err)
			fmt.Println(string(b))
			continue
		}
		typenum1 := fmt.Sprintf("%v", js.Get("data").Get("count").Interface())
		typenum, err := strconv.Atoi(typenum1)
		if err != nil {
			fmt.Printf("%vuser %v judgementOpinion data.count err: %v\n", other.TI(), userid, err)
			fmt.Println(string(b))
		}
		if otype == 1 {
			type1 = typenum
		} else {
			type2 = typenum
		}
		ops, err := js.Get("data").Get("opinion").Array()
		if err != nil && typenum != 0 {
			fmt.Printf("%vuser %v judgementOpinion json data.opinion err: %v\n", other.TI(), userid, err)
			fmt.Println(string(b))
			continue
		}
		for i := range ops {
			content, err := js.Get("data").Get("opinion").GetIndex(i).Get("content").String()
			if err != nil {
				fmt.Printf("%vuser %v judgementOpinion json data.opinion.%v.content err: %v\n", other.TI(), userid, i, err)
				fmt.Println(string(b))
				continue
			}
			vote, err := js.Get("data").Get("opinion").GetIndex(i).Get("vote").Int()
			if err != nil {
				fmt.Printf("%vuser %v judgementOpinion json data.opinion.%v.vote err: %v\n", other.TI(), userid, i, err)
				fmt.Println(string(b))
				continue
			}
			like, err := js.Get("data").Get("opinion").GetIndex(i).Get("like").Int()
			if err != nil {
				fmt.Printf("%vuser %v judgementOpinion json data.opinion.%v.like err: %v\n", other.TI(), userid, i, err)
				fmt.Println(string(b))
				continue
			}
			hate, err := js.Get("data").Get("opinion").GetIndex(i).Get("hate").Int()
			if err != nil {
				fmt.Printf("%vuser %v judgementOpinion json data.opinion.%v.hate err: %v\n", other.TI(), userid, i, err)
				fmt.Println(string(b))
				continue
			}
			opins = append(opins, opin{content, vote, like, hate})
		}
	}
	return
}

func judgementCaseObtain(userid int, btoml *other.Btoml) (id int, messagetxt string) {
	cookie := btoml.BiliUser[userid].Cookie
	aqw := strings.Index(cookie, "bili_jct=")
	if aqw == -1 {
		return
	}
	bqw := cookie[aqw+9:]
	cqw := strings.Index(bqw, ";")
	if cqw == -1 {
		return
	}
	csrf := bqw[:cqw]
	b, c := refererpost("https://api.bilibili.com/x/credit/jury/caseObtain", "jsonp=jsonp&csrf="+csrf, cookie, "https://www.bilibili.com/judgement/index")
	if !c {
		return
	}
	//fmt.Println(string(b))
	js, err := simplejson.NewJson(b)
	if err != nil {
		fmt.Printf("%vuser %v judgementCaseObtain json err: %v\n", other.TI(), userid, err)
		return
	}
	code, err := js.Get("code").Int()
	if err != nil {
		fmt.Printf("%vuser %v judgementCaseObtain json code err: %v\n", other.TI(), userid, err)
		return
	}
	message, err := js.Get("message").String()
	messagetxt = message
	if err != nil {
		fmt.Printf("%vuser %v judgementCaseObtain json message err: %v\n", other.TI(), userid, err)
		return
	}
	if code == 25014 {
		return -1, messagetxt
	} else if code != 0 {
		//fmt.Println(other.TI() + message)
		return
	}
	id, _ = js.Get("data").Get("id").Int()
	return
}

func judgement(userid int, btoml *other.Btoml) {
	for {
		id, mesg := judgementCaseObtain(userid, btoml)
		if id == 0 {
			if strings.Index(mesg, "请成为") != -1 {
				fmt.Printf("%v用户%v 不是风纪委员\n", other.TI(), userid)
				break
			} else {
				fmt.Printf("%v用户%v 没有接到案件: %v\n", other.TI(), userid, mesg)
			}
			time.Sleep(time.Minute * 10)
			continue
		} else if id == -1 {
			fmt.Printf("%v用户%v 今日案件已审满\n", other.TI(), userid)
			return
		}

		fmt.Printf("%v用户%v 接到案件%v\n", other.TI(), userid, id)
		var kill, del, no int
		for i := 0; i < 3; i++ {
			votes, yesnum, nonum := judgementOpinion(id, userid, btoml)
			//fmt.Println(votes, yesnum, nonum)
			for _, v := range votes {
				switch v.vote {
				case 1:
					kill += 2
					kill += v.like
					kill -= v.hate
				case 4:
					del += 2
					del += v.like
					del -= v.hate
				case 2:
					no += 2
					no += v.like
					no -= v.hate
				}
			}
			if yesnum != 0 {
				kill *= yesnum
				del *= yesnum
			}
			if nonum != 0 {
				no *= nonum
			}
			if yesnum+nonum > 3 {
				fmt.Printf("%v用户%v 案件%v 投票比 %v:%v:%v\n", other.TI(), userid, id, kill, del, no)
				break
			} else {
				if i < 2 {
					fmt.Printf("%v用户%v 案件%v 投票比 %v:%v:%v 等待5分钟\n", other.TI(), userid, id, kill, del, no)
					time.Sleep(time.Minute * 5)
				} else {
					fmt.Printf("%v用户%v 案件%v 投票比 %v:%v:%v 不等了\n", other.TI(), userid, id, kill, del, no)
				}
				continue
			}
		}
		if kill+del >= no/2 {
			if del >= kill {
				randSleep(1, 3, (userid+1)*23)
				judgementVote(id, 4, userid, btoml)
			} else {
				randSleep(1, 3, (userid+1)*23)
				judgementVote(id, 1, userid, btoml)
			}
		} else {
			randSleep(1, 3, (userid+1)*23)
			judgementVote(id, 2, userid, btoml)
		}

	}
}
