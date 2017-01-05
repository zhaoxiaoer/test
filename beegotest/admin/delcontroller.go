package admin

import (
	"fmt"

	"github.com2/astaxie/beego"
)

type DelController struct {
	beego.Controller
}

func (del *DelController) Get() {
	fmt.Printf("0 DelController Get\n")
	fmt.Printf("params: %v\n", del.Ct.Params)
	del.Ct.WriteString(fmt.Sprintf("DelController get, id: %v", del.Ct.Params[":id"]))
	fmt.Printf("1 DelController Get\n")
}

func (del *DelController) Post() {
	fmt.Printf("0 DelController Post\n")
	fmt.Printf("params: %v\n", del.Ct.Params)
	del.Ct.WriteString(fmt.Sprintf("DelController post, id: %v", del.Ct.Params[":id"]))
	fmt.Printf("1 DelController Post\n")
}
