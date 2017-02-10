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
	fmt.Printf("params: %v\n", edit.Ctx.Params)
	edit.Ctx.WriteString(fmt.Sprintf("EditController get, id: %v", edit.Ctx.Params[":id"]))
	fmt.Printf("1 EditController Get\n")
}

func (edit *EditController) Post() {
	fmt.Printf("0 EditController Post\n")
	edit.Ctx.WriteString("EditController post")
	fmt.Printf("1 EditController Post\n")
}
