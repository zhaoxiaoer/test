package sessionCtl

import (
	"fmt"

	"github.com/astaxie/beego"
)

type SessionCtl struct {
	beego.Controller
}

func (sessionCtl *SessionCtl) Get() {
	fmt.Printf("0 Session GET\n")
	username := sessionCtl.GetSession("username")
	sessionCtl.Data["username"] = username
	fmt.Printf("1 Session GET\n")
}

func (sessionCtl *SessionCtl) Post() {
	sessionCtl.Ctx.Request.ParseForm()
	username := sessionCtl.Ctx.Request.Form["username"][0]
	password := sessionCtl.Ctx.Request.Form["password"][0]
	fmt.Printf("username: %s, password: %s\n", username, password)
	sessionCtl.SetSession("username", username)
	sessionCtl.SetSession("password", password)
	sessionCtl.Redirect("/sessionctl", 301)
}
