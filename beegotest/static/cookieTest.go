package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func helloServer(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		fmt.Println(err)
	}
	// 判断是否已经设置了cookie
	if cookie == nil {
		now := time.Now()
		cookie := http.Cookie{Name: "username", Value: "zhaoxiaoer", Expires: now.AddDate(0, 0, 1)}
		http.SetCookie(w, &cookie)
		io.WriteString(w, "hello, world")
	} else {
		//		cookie := http.Cookie{Name: "username", Value: "zhaoxiaoer", Expires: now.AddDate(0, 0, -1)}
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		io.WriteString(w, "欢迎回来，"+cookie.Value)
	}
}

func main() {
	http.HandleFunc("/hello", helloServer)
	log.Fatal(http.ListenAndServe(":6725", nil))
}
