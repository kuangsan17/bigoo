package other

import (
	"crypto/md5"
	"encoding/hex"
)

func calcSign(text, secret string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(text + secret))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//Sign ...
func Sign(text string) string {
	return text + "&sign=" + calcSign(text, "560c52ccd288fed045859ed18bffd973")
}

//TvSign ...
func TvSign(text string) string {
	return text + "&sign=" + calcSign(text, "59b43e04ad6965f34319062b478f83dd")
}

//BlueSign ...
func BlueSign(text string) string {
	return text + "&sign=" + calcSign(text, "25bdede4e1581c836cab73a48790ca6e")
}

//SpSign ...
func SpSign(text string) string {
	return text + "&sign=" + calcSign(text, "aHRmhWMLkdeMulLqORnYZocwMBpMEOdt")
}
