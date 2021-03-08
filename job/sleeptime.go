package job

import (
	"math/rand"
	"time"
)

func randSleep(timea, timeb int, aa int) {
	if timeb <= 1 {
		return
	}
	rand.Seed(time.Now().UnixNano() - int64(aa))
	if timea < 0 {
		a := rand.Intn(timeb*1000 - 800)
		time.Sleep(time.Duration(int64(a) * 1000000))
	}
	b := timeb - timea
	c := rand.Intn(b*1000-1000) + timea*1000 + 500
	time.Sleep(time.Duration(int64(c) * 1000000))
}
