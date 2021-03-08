package job

import (
	"time"

	"github.com/atwat/bigoo/other"
)

//FuckGray ...
func FuckGray(btoml *other.Btoml) {
	for userid := range btoml.BiliUser {
		go func(userid int) {
			for btoml.MoreSetting.Judgement {
				if other.NowHour() == 1 || other.NowHour() == 2 {
					randSleep(0, 300, (userid+1)*22)
					//go judgement(userid, btoml)

					time.Sleep(time.Hour * 20)
				} else {
					time.Sleep(time.Minute * 33)
				}
			}
		}(userid)
	}
}

func keepMedalNotGray(k int, btoml *other.Btoml) {

}
