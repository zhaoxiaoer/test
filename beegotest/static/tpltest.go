/*
header.tpl 文件内容:
{{define "header"}}
<html>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
  <head>
    <title>模板测试</title>
  </head>
{{end}}

content.tpl 文件内容:
{{define "content111"}}
{{template "header"}}
  <body>
    <h1>演示嵌套</h1>
    <ul>
	  <li>使用define定义模板</li>
	  <li>使用template调用模板</li>
	</ul>
  </body>
{{template "footer"}}
{{end}}

footer.tpl 文件内容:
{{define "footer"}}
</html>
{{end}}

hello.txt 文件内容:
hello, world
{{define "hello"}}
hello, world 222
{{end}}
*/

package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("header.tpl", "content.tpl", "footer.tpl", "hello.txt")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// t.ExecuteTemplate 的第二个参数要么是文件名，要么是模板名
	// 如果既是文件名，又是模板名，系统将报错误
	err = t.ExecuteTemplate(os.Stdout, "content111", nil)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = t.ExecuteTemplate(os.Stdout, "hello.txt", nil) // 文件名
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = t.ExecuteTemplate(os.Stdout, "hello", nil) // 模板名
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
