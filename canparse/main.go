package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/websocket"
)

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
	fmt.Printf("file: %s\n", r.FormValue("file"))
	tpls := template.Must(template.ParseFiles("templates/wsserver.html"))
	tpls.ExecuteTemplate(w, "wsserver", r.FormValue("file"))
}

func wsServer(ws *websocket.Conn) {
	fmt.Printf("22222222\n")
	buf := make([]byte, 1024)
	ws.Read(buf)
	ws.Write(buf)
}

func main() {
	mux := http.NewServeMux()
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
