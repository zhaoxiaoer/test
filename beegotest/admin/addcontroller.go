package admin

import (
	"fmt"

	"github.com2/astaxie/beego"
)

type AddController struct {
	beego.Controller
}

func (add *AddController) Get() {
	fmt.Printf("0 AddController Get\n")
	//	add.Ct.WriteString("AddController")
	add.Layout = "admin/layout.html"
	add.TplNames = "admin/add.tpl"
	add.Data["url"] = "/admin/addpkg"
	fmt.Printf("1 AddController Get\n")
}

func (add *AddController) Post() {
	fmt.Printf("0 AddController\n")

	fmt.Printf("0 %v\n", add.Ct.Request.Form)
	add.Ct.Request.ParseForm()
	fmt.Printf("1 %v\n", add.Ct.Request.Form)

	add.Ct.WriteString("AddController Post")

	fmt.Printf("1 AddController\n")
}
