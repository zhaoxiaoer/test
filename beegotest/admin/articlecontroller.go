package admin

import (
	"fmt"

	"github.com/astaxie/beego"
)

type ArticleController struct {
	beego.Controller
}

func (article *ArticleController) Get() {
	fmt.Printf("0 ArticleController Get\n")
	article.Ctx.WriteString("ArticleController Get")
	fmt.Printf("1 ArticleController Get\n")
}
