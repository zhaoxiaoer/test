package proemulator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//	"net/url"
	"strings"

	"github.com/astaxie/beego"
)

type Emulator struct {
	beego.Controller
}

func (emu *Emulator) Get() {
	fmt.Printf("0 Emulator Get\n")

	emu.Layout = "proemulator/layout.html"
	emu.TplNames = "proemulator/emulator.tpl"
	emu.Data["url"] = "/proemulator/emulator"
	emu.Data["posturl"] = "http://he.bobdz.com:4026/place/eshe/pukey"

	fmt.Printf("1 Emulator Get\n")
}

func (emu *Emulator) Post() {
	fmt.Printf("0 Emulator\n")

	fmt.Printf("0 %v\n", emu.Ctx.Request.Form)
	emu.Ctx.Request.ParseForm()
	fmt.Printf("1 %v\n", emu.Ctx.Request.Form)
	postUrl := emu.Ctx.Request.Form["posturl"][0]
	postBody := emu.Ctx.Request.Form["postbody"][0]

	//	body := httpPost(postUrl, postBody)
	body := httpDo(postUrl, postBody)
	//	emu.Ct.WriteString(body)

	emu.Layout = "proemulator/layout.html"
	emu.TplNames = "proemulator/emulator.tpl"
	emu.Data["url"] = "/proemulator/emulator"
	emu.Data["posturl"] = postUrl
	emu.Data["body"] = body

	fmt.Printf("1 Emulator\n")
}

func httpGet(myUrl string) string {
	resp, err := http.Get(myUrl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err.Error()
	}
	fmt.Printf("%s\n", string(body))

	return string(body)
}

func httpPost(myUrl string, postBody string) string {
	resp, err := http.Post(myUrl, "application/x-www-form-urlencoded", strings.NewReader(postBody))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err.Error()
	}
	fmt.Printf("%s\n", string(body))

	return string(body)
}

//func httpPostForm(myUrl string, body string) string {
//	resp, err := http.PostForm(myUrl, values)
//	if err != nil {
//		fmt.Printf("err: %v\n", err)
//		return
//	}
//	defer resp.Body.Close()

//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Printf("err: %v\n", err)
//		return
//	}
//	fmt.Printf("%s\n", string(body))

//	return string(body)
//}

func httpDo(myUrl string, postBody string) string {
	client := &http.Client{}

	req, err := http.NewRequest("POST", myUrl, strings.NewReader(postBody))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err.Error()
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=zhaoxiaoer")
	fmt.Printf("%v\n", req.Header)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err.Error()
	}
	fmt.Printf("%s\n", string(body))

	return string(body)
}
