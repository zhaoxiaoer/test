// go中有pprof包来做代码的性能监控
// net/http/pprof
// runtime/pprof
// 其中net/http/pprof包只是将runtime/pprof包又封装了一下，并在http端口上暴露出来而已
package main

import (
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":6725", nil)
	log.Fatal(err)
}
