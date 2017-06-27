package admin

import (
	"fmt"

	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (user *UserController) Get() {
	fmt.Printf("0 UserController Get\n")
	user.Ctx.WriteString("UserController Get")
	fmt.Printf("1 UserController Get\n")
}

type UserControllerJSON struct {
	beego.Controller
}

func (userj *UserControllerJSON) Get() {
	fmt.Printf("0 UserControllerJSON Get\n")
	type Info struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	info := Info{
		Name:     "xiaoer",
		Password: "123456",
	}
	userj.Data["json"] = &info
	userj.ServeJson()
	fmt.Printf("1 UserControllerJSON Get\n")
}

type UserControllerJSONP struct {
	beego.Controller
}

func (userjp *UserControllerJSONP) Get() {
	userjp.TplNames = "./admin/userjp.tpl"
}

type UserControllerJSONP2 struct {
	beego.Controller
}

func (userjp2 *UserControllerJSONP2) Get() {
	fmt.Printf("UserControllerJSONP2 begin\n")
	type Info struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Message string `json:"message"`
	}
	info := Info{
		Name:    "xiaoer",
		Address: "machang",
		Message: userjp2.Ctx.Request.Form.Get("message"),
	}
	userjp2.Data["jsonp"] = &info
	userjp2.ServeJsonp()
	fmt.Printf("UserControllerJSONP2 end\n")
}
