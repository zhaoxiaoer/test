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
	"regexp"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

type CanMsg struct {
	To  int64 // Time Offset
	CID int64 // Can ID
	Val int64
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
	for {
		s, err := br.ReadString('\n')

		re := regexp.MustCompile(`\d+\)\s+(\d+).(\d{3})\s+\d\s+Rx\s+([[:xdigit:]]{4})\s+-\s+8\s+([[:xdigit:]]{2})\s([[:xdigit:]]{2})\s([[:xdigit:]]{2})\s([[:xdigit:]]{2})\s([[:xdigit:]]{2})\s([[:xdigit:]]{2})\s([[:xdigit:]]{2})\s([[:xdigit:]]{2})`)
		//		fmt.Printf("%q\n", re.FindStringSubmatch(s))
		ss := re.FindStringSubmatch(s)
		if len(ss) == 12 {
			to, err0 := strconv.ParseInt(ss[1]+ss[2], 10, 64)
			cid, err1 := strconv.ParseInt(ss[3], 16, 64)
			v, err2 := strconv.ParseInt(ss[4], 16, 64)
			if err0 == nil && err1 == nil && err2 == nil {
				fmt.Printf("to: %v, cid: %v, v: %v\n", to, cid, v)
				d := CanMsg{to, cid, v}
				err := websocket.JSON.Send(conn, d)
				if err != nil {
					return
				}
				time.Sleep(1000 * time.Millisecond)
			}
		}

		if err != nil {
			return
		}
	}
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
