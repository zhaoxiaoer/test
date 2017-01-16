package main

import (
	"fmt"
	"net/http"

	"./admin"
	"./controllers"
	"./proemulator"
	"./qrcode"

	"github.com2/astaxie/beego"
	//	"gopkg.in/mgo.v2-unstable"
	//	"gopkg.in/mgo.v2-unstable/bson"
)

type TestHtmlMethod struct {
	beego.Controller
}

func (test *TestHtmlMethod) Get() {
	fmt.Printf("0 TestHtmlMethod: Get\n")
	test.Ct.Request.ParseForm()
	fmt.Printf("form: %v\n", test.Ct.Request.Form)

	test.Ct.WriteString("test html [GET] method")
	fmt.Printf("1 TestHtmlMethod: Get\n")
}

func (test *TestHtmlMethod) Post() {
	fmt.Printf("0 TestHtmlMethod: Post\n")
	test.Ct.Request.ParseForm()
	fmt.Printf("form: %v\n", test.Ct.Request.Form)

	test.Ct.WriteString("test html [POST] method")
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
	beego.BeeApp.RegisterController("/", &controllers.MainController{})

	//	fmt.Printf("main 1\n")
	//	beego.BeeApp.RegisterController("/testhtmlmethod/test.html", &TestHtmlMethod{})

	fmt.Printf("main 2\n")
	beego.BeeApp.RegisterController("/admin", &admin.UserController{})
	beego.BeeApp.RegisterController("/admin/index", &admin.ArticleController{})
	beego.BeeApp.RegisterController("/admin/addpkg", &admin.AddController{})

	fmt.Printf("main 22\n")
	beego.BeeApp.RegisterController("/admin/editpkg/:id([0-9]+)", &admin.EditController{})
	beego.BeeApp.RegisterController("/admin/delpkg/:id([0-9]+)", &admin.DelController{})

	fmt.Printf("main 2222233333\n")
	beego.BeeApp.RegisterController("/proemulator/emulator", &proemulator.Emulator{})

	fmt.Printf("main 4444444444\n")
	beego.BeeApp.RegisterController("/qrcode/qrcode", &qrcode.Qrcode{})

	fmt.Printf("main 33333\n")
	beego.BeeApp.RegisterController("/:pkg(.*)", &controllers.MainController{})

	fmt.Printf("main 3\n")
	beego.BeeApp.SetStaticPath("/public", "public")

	fmt.Printf("main 4\n")
	var FilterUser = func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("FilterUser\n")
		//		if r.URL.User == nil || r.URL.User.Username() != "admin" {
		//			http.Error(w, "user error", http.StatusUnauthorized)
		//		}
	}
	beego.BeeApp.Filter(FilterUser)
	beego.BeeApp.FilterParam("id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Filter id\n")
		id := r.URL.Query().Get(":id")
		fmt.Printf("id: %v\n", id)
		if id == "3" {
			s := fmt.Sprintf("id: [%s] error\n", id)
			http.Error(w, s, http.StatusUnauthorized)
		}
	})
	beego.BeeApp.FilterPrefixPath("/admin/delpkg", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Filter prefix path [/admin/delpkg]\n")
		id := r.URL.Query().Get(":id")
		if id == "4" {
			s := fmt.Sprintf("can not delete id: [%s]\n", id)
			http.Error(w, s, http.StatusUnauthorized)
		}
	})

	beego.BeeApp.Run()
}
