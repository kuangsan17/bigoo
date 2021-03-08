package other

import (
	"fmt"
)

//SendMail ...
func SendMail(text string, btoml *Btoml) {
	sendMailWx(text, btoml)
	fmt.Println(TI()+"发送：", text)
}

func sendMailWx(text string, btoml *Btoml) {
	if btoml.Mail.Wxserver.Sckey == "" {
		return
	}
	url := "http://sc.ftqq.com/" + btoml.Mail.Wxserver.Sckey + ".send?text=" + text
	PcReqGet(url, "")
}
