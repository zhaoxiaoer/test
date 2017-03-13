package sessionCtl

import (
	"fmt"

	"github.com2/astaxie/beego"
)

type SessionCtl struct {
	beego.Controller
}

func (sessionCtl *SessionCtl) Get() {
	fmt.Printf("0 Session GET\n")
	session := beego.GlobalSessions.SessionStart(sessionCtl.Ctx.ResponseWriter, sessionCtl.Ctx.Request)
	username := session.Get("username")
	sessionCtl.Data["username"] = username
	fmt.Printf("1 Session GET\n")
}

func (sessionCtl *SessionCtl) Post() {
	sessionCtl.Ctx.Request.ParseForm()
	username := sessionCtl.Ctx.Request.Form["username"][0]
	password := sessionCtl.Ctx.Request.Form["password"][0]
	fmt.Printf("username: %s, password: %s\n", username, password)
	session := beego.GlobalSessions.SessionStart(sessionCtl.Ctx.ResponseWriter, sessionCtl.Ctx.Request)
	session.Set("username", username)
	session.Set("password", password)
	sessionCtl.Redirect("/sessionctl", 301)
}
