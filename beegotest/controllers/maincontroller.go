package controllers

import (
	"fmt"

	"github.com2/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (main *MainController) Get() {
	fmt.Printf("0 MainController Get\n")
	//	fmt.Printf("params: %v\n", main.Ct.Params)
	//	main.Ct.WriteString("Hello, beego")
	main.Layout = "mainlayout.html"
	main.TplNames = "main.tpl"
	fmt.Printf("1 MainController Get\n")
}
