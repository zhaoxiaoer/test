package qrcode

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com2/astaxie/beego"
)

type Qrcode struct {
	beego.Controller
}

func (qr *Qrcode) Get() {
	fmt.Printf("0 Qrcode Get\n")

	qr.Layout = "qrcode/layout.html"
	qr.TplNames = "qrcode/qrcode.tpl"

	qr.Data["time"] = time.Now()
	qr.Data["format"] = "Y-m-d G:i:s"

	qr.Data["url"] = "/qrcode/qrcode"
	qr.Data["content"] = "http://192.168.1.7:8080/"
	//	qr.Data["image"] = "/static/homepage.png"
	png, err := qrcode.Encode("http://192.168.1.7:8080/", qrcode.Highest, 256)
	if err != nil {
		qr.Data["errinfo"] = err.Error()
	} else {
		pngStr := base64.StdEncoding.EncodeToString(png)
		qr.Data["image"] = pngStr
	}

	fmt.Printf("1 Qrcode Get\n")
}

func (qr *Qrcode) Post() {
	fmt.Printf("0 Qrcode Post\n")

	fmt.Printf("0 %v\n", qr.Ctx.Request.Form)
	qr.Ctx.Request.ParseForm()
	fmt.Printf("1 %v\n", qr.Ctx.Request.Form)
	content := qr.Ctx.Request.Form["content"][0]

	qr.Layout = "qrcode/layout.html"
	qr.TplNames = "qrcode/qrcode.tpl"

	qr.Data["time"] = time.Now()
	qr.Data["format"] = "Y-m-d G:i:s"

	qr.Data["url"] = "/qrcode/qrcode"
	qr.Data["content"] = content
	//	qr.Data["image"] = "/static/qr.png"

	//	err := qrcode.WriteFile(content, qrcode.Highest, 256, "./static/qr.png")
	png, err := qrcode.Encode(content, qrcode.Highest, 256)
	if err != nil {
		qr.Data["errinfo"] = err.Error()
	} else {
		pngStr := base64.StdEncoding.EncodeToString(png)
		qr.Data["image"] = pngStr
	}

	fmt.Printf("1 Qrcode Post\n")
}
