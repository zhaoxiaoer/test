package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	//	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type CanInfo struct {
	It  int8 // info type
	Val int32
}

func connState(c net.Conn, cs http.ConnState) {
	fmt.Printf("c: %v, state: %v\n", c.RemoteAddr().String(), cs)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, fh, err := r.FormFile("uploaded")
		if err == nil {
			data, err := ioutil.ReadAll(file)
			if err == nil {
				ioutil.WriteFile("canfiles/"+fh.Filename, data, 777)
			}
		}
	}

	fis := make([]os.FileInfo, 0)
	err := filepath.Walk("canfiles", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", path)

		if f.IsDir() == false {
			fis = append(fis, f)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}

	tpls := template.Must(template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/content.html"))
	tpls.ExecuteTemplate(w, "layout", fis)
}

func fileInfo(w http.ResponseWriter, r *http.Request) {
	tpls := template.Must(template.ParseFiles("templates/wsserver.html"))
	tpls.ExecuteTemplate(w, "wsserver", r.FormValue("file"))
}

func wsServer(conn *websocket.Conn) {
	r := conn.Request()
	f, err := os.Open("canfiles/" + r.FormValue("file"))
	if err != nil {
		conn.Write([]byte("Can not open the file!"))
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var i int = 0
	for {
		s, err := br.ReadString('\n')
		conn.Write([]byte(s))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Millisecond)
		i++
		if i > 300 {
			break
		}
	}
	//	for i := 0; i < 10000; i++ {
	//		var data CanInfo
	//		if i%2 == 0 {
	//			data = CanInfo{1, int32(i)}
	//		} else {
	//			data = CanInfo{2, int32(i)}
	//		}
	//		err := websocket.JSON.Send(conn, data)
	//		if err != nil {
	//			return
	//		}
	//		time.Sleep(1000 * time.Millisecond)
	//	}
}

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/public/", http.StripPrefix("/public/", files))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/fileinfo", fileInfo)
	mux.Handle("/wsserver", websocket.Handler(wsServer))

	server := http.Server{
		Addr:      ":8080",
		ConnState: connState,
		Handler:   mux,
	}

	server.ListenAndServe()
	//server.ListenAndServeTLS("cert.pem", "key.pem")
}
