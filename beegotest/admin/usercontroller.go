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
	//	userj.Data2 = make(map[string]interface{})
	//	userj.Data2["name"] = "xiaoer"
	//	userj.Data2["password"] = "123456"
	//	userj.ServeJson()
	fmt.Printf("1 UserControllerJSON Get\n")
}
