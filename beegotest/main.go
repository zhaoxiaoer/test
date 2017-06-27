package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/zhaoxiaoer/test/beegotest/admin"
	"github.com/zhaoxiaoer/test/beegotest/controllers"
	"github.com/zhaoxiaoer/test/beegotest/proemulator"
	"github.com/zhaoxiaoer/test/beegotest/qrcode"
	"github.com/zhaoxiaoer/test/beegotest/sessionCtl"
	"github.com/zhaoxiaoer/test/beegotest/uploadFile"
	"github.com/zhaoxiaoer/test/beegotest/wsserver"

	"github.com/astaxie/beego"
	"golang.org/x/net/websocket"
	//	"gopkg.in/mgo.v2-unstable"
	//	"gopkg.in/mgo.v2-unstable/bson"
)

type TestHtmlMethod struct {
	beego.Controller
}

func (test *TestHtmlMethod) Get() {
	fmt.Printf("0 TestHtmlMethod: Get\n")
	test.Ctx.Request.ParseForm()
	fmt.Printf("form: %v\n", test.Ctx.Request.Form)

	test.Ctx.WriteString("test html [GET] method")
	fmt.Printf("1 TestHtmlMethod: Get\n")
}

func (test *TestHtmlMethod) Post() {
	fmt.Printf("0 TestHtmlMethod: Post\n")
	test.Ctx.Request.ParseForm()
	fmt.Printf("form: %v\n", test.Ctx.Request.Form)

	test.Ctx.WriteString("test html [POST] method")
	fmt.Printf("1 TestHtmlMethod: Post\n")
}

type Person struct {
	Name  string
	Phone string
}

// beego 先判断URL是否符合静态路径，然后是固定路径，最后是动态路径（正则表达式路径）
// 找到路由后，先进行滤波，然后才会执行prepare, get等函数
func main() {
	fmt.Printf("main 000000000000\n")
	//	session, err := mgo.Dial("192.168.1.7:27017")
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer session.Close()
	//	// Optional. Switch the session to a monotonic behavior.
	//	session.SetMode(mgo.Monotonic, true)
	//	c := session.DB("foobar").C("people")
	//	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"}, &Person{"Cla", "+55 53 8402 8510"})
	//	if err != nil {
	//		fmt.Printf("DB error: %v\n", err)
	//	}
	//	result := Person{}
	//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	//	if err != nil {
	//		fmt.Printf("DB error 2: %v\n", err)
	//	}
	//	fmt.Printf("Phone: %v\n", result.Phone)

	fmt.Printf("main 0\n")
	beego.Router("/", &controllers.MainController{})

	//	fmt.Printf("main 1\n")
	//	beego.BeeApp.RegisterController("/testhtmlmethod/test.html", &TestHtmlMethod{})

	fmt.Printf("main 2\n")
	beego.Router("/admin", &admin.UserController{})
	beego.Router("/adminjson", &admin.UserControllerJSON{})
	beego.Router("/admin/index", &admin.ArticleController{})
	beego.Router("/admin/addpkg", &admin.AddController{})

	fmt.Printf("main 22\n")
	beego.Router("/admin/editpkg/:id([0-9]+)", &admin.EditController{})
	beego.Router("/admin/delpkg/:id([0-9]+)", &admin.DelController{})

	fmt.Printf("main 2222233333\n")
	beego.Router("/proemulator/emulator", &proemulator.Emulator{})

	fmt.Printf("main 4444444444\n")
	beego.Router("/qrcode/qrcode", &qrcode.Qrcode{})

	fmt.Printf("main 555555555\n")
	beego.Router("/sessionctl", &sessionCtl.SessionCtl{})

	fmt.Printf("main 666666666\n")
	beego.Router("/uploadFile", &uploadFile.UploadFile{})

	//	fmt.Printf("main 33333\n")
	//	beego.Router("/:pkg(.*)", &controllers.MainController{})
	fmt.Printf("add wsserver\n")
	beego.Router("/wsserver", &wsserver.Wsserver{})
	beego.RouterHandler("/chat", websocket.Handler(wsserver.ChatHandler))

	beego.Router("/jsonptest", &admin.UserControllerJSONP{})
	beego.Router("/jsonptest2", &admin.UserControllerJSONP2{})

	fmt.Printf("main 3\n")
	beego.BeeApp.SetStaticPath("/public", "public")

	fmt.Printf("main 4\n")
	var FilterUser = func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("FilterUser\n")
		//		if r.URL.User == nil || r.URL.User.Username() != "admin" {
		//			http.Error(w, "user error", http.StatusUnauthorized)
		//		}
	}
	beego.Filter(FilterUser)
	beego.FilterParam("id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Filter id\n")
		id := r.URL.Query().Get(":id")
		fmt.Printf("id: %v\n", id)
		if id == "3" {
			s := fmt.Sprintf("id: [%s] error\n", id)
			http.Error(w, s, http.StatusUnauthorized)
		}
	})
	beego.FilterPrefixPath("/admin/delpkg", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Filter prefix path [/admin/delpkg]\n")
		id := r.URL.Query().Get(":id")
		if id == "4" {
			s := fmt.Sprintf("can not delete id: [%s]\n", id)
			http.Error(w, s, http.StatusUnauthorized)
		}
	})

	beego.Run()
}
