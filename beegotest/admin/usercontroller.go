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
	user.Ct.WriteString("UserController Get")
	fmt.Printf("1 UserController Get\n")
}
