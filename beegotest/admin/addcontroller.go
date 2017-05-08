package admin

import (
	"fmt"

	"github.com/astaxie/beego"
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

	fmt.Printf("0 %v\n", add.Ctx.Request.Form)
	add.Ctx.Request.ParseForm()
	fmt.Printf("1 %v\n", add.Ctx.Request.Form)

	add.Ctx.WriteString("AddController Post")

	fmt.Printf("1 AddController\n")
}
