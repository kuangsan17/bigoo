package other

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

//Btoml is all message from txt
type Btoml struct {
	Monitor     monitorInfo
	Setting     settingInfo
	MoreSetting moreSettingInfo
	Mail        mailInfo
	BiliUser    []biliUserInfo
}

type biliUserInfo struct {
	UserName     string
	PassWord     string
	AccessToken  string
	RefreshToken string
	Cookie       string
}

type monitorInfo struct {
	Bilive string
	YjAddr string
	YjKey  string
	Local  bool
	Lun    bool
}

type settingInfo struct {
	Lottery         lotteryInfo
	MainDailyJob    bool //主站直播漫画
	LiveDailyJob    bool
	MangaDailyJob   bool
	LiveOnlineHeart bool
}

type moreSettingInfo struct {
	Judgement                 bool
	KeepMedalColorful         bool
	CleanExpiringMangaCoupons bool
	UseCoin                   useCoinInfo
	UseSilver                 useSliverInfo
	UseBag                    sendBagInfo
}

type lotteryInfo struct {
	GuardOdds int
	GiftsOdds int
	PkOdds    int
	StormOdds int
	StormSet  [2]int
	BoxOdds   int
	BoxStart  int
	RedOdds   int
	SleepTime [][2]int
}

type mailInfo struct {
	Wxserver wxserverInfo
	Email    emailInfo
}

type wxserverInfo struct {
	Sckey string
}
type emailInfo struct {
	Recmail string
	User    string
	Pass    string
	Host    string
	Port    string
}

type useCoinInfo struct {
	SendCoinNum    int
	Send2SpecialUp bool
	//	Send2UpUids    []int
	Coin2Sliver int
}

type useSliverInfo struct {
	BuyHuima                     int
	Sliver2Coin                  int
	UseForWearingMedalIfNoLatiao bool
}

type sendBagInfo struct {
	Send2Wearing bool
	//	Send2OtherUpUids                     []int
	CleanExpiring2OtherMedal             bool
	CleanExpiring2OtherMedalFromHigh2Low bool
	//	CleanExpiring2Roomid                 []int
}

//WriteToml ...
func (T *Btoml) WriteToml() {
	f, err := os.Create(tomlFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := toml.NewEncoder(f).Encode(T); err != nil {
		fmt.Println(err)
		return
	}
}

//AddUser add user
func (T *Btoml) AddUser(uname, upass string) bool {
	for k, v := range T.BiliUser {
		if v.UserName == uname {
			if v.PassWord == upass {
				fmt.Println("已存在该用户")
				return true
			}
			line := ""
			fmt.Print("已存在该用户，是否使用新密码[Y/n]:")
			fmt.Scanln(&line)
			if line == "n" || line == "N" {
				return true
			}
			T.BiliUser[k].PassWord = upass
			var boo bool
			T.BiliUser[k].AccessToken, T.BiliUser[k].RefreshToken, T.BiliUser[k].Cookie, boo = Login(uname, upass, "")
			if boo {
				T.WriteToml()
				return true
			}
			return false
		}
	}
	var newuser biliUserInfo
	var boo bool
	newuser.UserName = uname
	newuser.PassWord = upass
	newuser.AccessToken, newuser.RefreshToken, newuser.Cookie, boo = Login(uname, upass, "")
	if boo {
		fmt.Printf("%v用户%v登录成功\n", TI(), uname)
		T.BiliUser = append(T.BiliUser, newuser)
		T.WriteToml()
		return true
	}
	return false
}

//AddUserByPrint add user
func (T *Btoml) AddUserByPrint() {
	for {
		var uname, upass string
		for {
			fmt.Print("请输入账号:")
			fmt.Scanln(&uname)
			if uname != "" {
				break
			}
		}
		for {
			fmt.Print("请输入密码:")
			fmt.Scanln(&upass)
			if upass != "" {
				break
			}
		}
		if T.AddUser(uname, upass) {
			break
		}
	}
}

//CheckCookie ...
func (T *Btoml) CheckCookie(k int) {
	tt := Iflogin(T.BiliUser[k].AccessToken, T.BiliUser[k].Cookie)
	if tt < 86400 {
		fmt.Printf("%v用户%v Token可能过期\n", TI(), k)
		T.refresh(k)
	} else {
		fmt.Printf("%v用户%v Token于%v.%v天后过期并将自动更新\n", TI(), k, tt/86400, tt/864-tt/86400*100)
		go func(tt, k int) {
			time.Sleep(time.Second * time.Duration(tt-86400))
			T.refresh(k)
		}(tt, k)
	}
}

func (T *Btoml) refresh(k int) {
	var boo bool
	T.BiliUser[k].AccessToken, T.BiliUser[k].RefreshToken, T.BiliUser[k].Cookie, boo = Refresh(T.BiliUser[k].AccessToken, T.BiliUser[k].RefreshToken)
	if boo {
		fmt.Printf("%v用户%v Token更新成功\n", TI(), k)
		T.WriteToml()
		go func() {
			time.Sleep(time.Hour * 696)
			T.refresh(k)
		}()
	} else {
		fmt.Printf("%v用户%v Token更新失败\n", TI(), k)
		T.relogin(k)
	}
}

func (T *Btoml) relogin(k int) {
	var boo bool
	T.BiliUser[k].AccessToken, T.BiliUser[k].RefreshToken, T.BiliUser[k].Cookie, boo = Login(T.BiliUser[k].UserName, T.BiliUser[k].PassWord, "")
	if boo {
		fmt.Printf("%v用户%v 重新登录成功\n", TI(), k)
		T.WriteToml()
		go func() {
			time.Sleep(time.Hour * 696)
			T.refresh(k)
		}()
	} else {
		fmt.Printf("%v用户%v 重新登录失败\n", TI(), k)
	}
}

var tomlFilePath string = "./bigoo.toml"

func firstStart(tomljson string) {
	PutPayImg()
	monitor := monitorInfo{"", "", "", false, false}
	setting := settingInfo{
		lotteryInfo{100, 100, 100, 100, [2]int{90, 20}, 0, 661, 0, [][2]int{{163000, 170500}}},
		true, true, true, true,
	}
	moresetting := moreSettingInfo{
		false,
		true,
		true,
		useCoinInfo{0, false, 0},
		useSliverInfo{0, 0, false},
		sendBagInfo{false, false, true},
	}
	mail := mailInfo{wxserverInfo{""}, emailInfo{"", "", "", "smtp.qq.com", "465"}}
	biliuser := []biliUserInfo{}
	t := Btoml{monitor, setting, moresetting, mail, biliuser}
	if tomljson != "" {
		var nt Btoml
		err := json.Unmarshal([]byte(tomljson), &nt)
		if err != nil {
			fmt.Println("传入toml参数有误！")
			fmt.Println(strings.ReplaceAll(tomljson, "[", "【"))
			time.Sleep(time.Second * 3)
			os.Exit(0)
		}
		t = nt
	}
	t.WriteToml()
}

//TomlStart return the struct of toml
func TomlStart(tomljson string) *Btoml {
	if _, err := os.Stat(tomlFilePath); err != nil {
		firstStart(tomljson)
	}
	var config *Btoml
	if _, err := toml.DecodeFile(tomlFilePath, &config); err != nil {
		fmt.Println(TI() + "请检查配置文件！")
		time.Sleep(time.Second)
		panic(err)
	}
	if len(config.BiliUser) == 0 {
		config.AddUserByPrint()
	}
	for i, v := range config.BiliUser {
		if v.AccessToken == "" || v.RefreshToken == "" || v.Cookie == "" {
			config.relogin(i)
		} else {
			config.CheckCookie(i)
		}
	}
	return config
}
