package qrcode

import (
	"fmt"

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
	qr.Data["url"] = "/qrcode/qrcode"
	qr.Data["content"] = "http://192.168.1.7:8080/"
	qr.Data["image"] = "/static/homepage.png"

	fmt.Printf("1 Qrcode Get\n")
}

func (qr *Qrcode) Post() {
	fmt.Printf("0 Qrcode Post\n")

	fmt.Printf("0 %v\n", qr.Ct.Request.Form)
	qr.Ct.Request.ParseForm()
	fmt.Printf("1 %v\n", qr.Ct.Request.Form)
	content := qr.Ct.Request.Form["content"][0]

	err := qrcode.WriteFile(content, qrcode.Highest, 256, "./static/qr.png")
	if err != nil {
		qr.Data["errinfo"] = err.Error()
	}

	qr.Layout = "qrcode/layout.html"
	qr.TplNames = "qrcode/qrcode.tpl"
	qr.Data["url"] = "/qrcode/qrcode"
	qr.Data["content"] = content
	qr.Data["image"] = "/static/qr.png"

	fmt.Printf("1 Qrcode Post\n")
}
