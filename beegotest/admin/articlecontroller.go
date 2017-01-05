package admin

import (
	"fmt"

	"github.com2/astaxie/beego"
)

type ArticleController struct {
	beego.Controller
}

func (article *ArticleController) Get() {
	fmt.Printf("0 ArticleController Get\n")
	article.Ct.WriteString("ArticleController Get")
	fmt.Printf("1 ArticleController Get\n")
}
