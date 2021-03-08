package job

import "github.com/atwat/bigoo/other"

const (
	build       = "6182200"
	channel     = "huawei"
	statistics  = `{"appId":1,"platform":3,"version":"6.18.2","abtest":""}`
	statistics1 = `%7B%22appId%22%3A1%2C%22platform%22%3A3%2C%22version%22%3A%226.18.2%22%2C%22abtest%22%3A%22%22%7D`
	version     = `6.18.2`
)

//AllJobs ...
func AllJobs(btoml *other.Btoml) {
	go LiveDaily(btoml)
	go LiveHeart(btoml)
	go MangaDaily(btoml)
	go MainSvipDaily(btoml)
	go UseSilver(btoml)
	go ClearBag(btoml)
	go judgementMain(btoml)
}
