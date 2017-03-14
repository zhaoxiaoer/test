package admin

import (
	"fmt"

	"github.com2/astaxie/beego"
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
