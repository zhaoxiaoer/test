package admin

import (
	"fmt"

	"github.com2/astaxie/beego"
)

type EditController struct {
	beego.Controller
}

func (edit *EditController) Get() {
	fmt.Printf("0 EditController Get\n")
	fmt.Printf("params: %v\n", edit.Ct.Params)
	edit.Ct.WriteString(fmt.Sprintf("EditController get, id: %v", edit.Ct.Params[":id"]))
	fmt.Printf("1 EditController Get\n")
}

func (edit *EditController) Post() {
	fmt.Printf("0 EditController Post\n")
	edit.Ct.WriteString("EditController post")
	fmt.Printf("1 EditController Post\n")
}
