package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"golang.org/x/net/websocket"
)

type ContParam struct {
	Host string
	Fis  []os.FileInfo
}

type WsParam struct {
	Host string
	File string
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
			to, toErr := strconv.ParseInt(ss[1]+ss[2], 10, 64)
			did, didErr := strconv.ParseInt(ss[3], 16, 64)
			v0, err0 := strconv.ParseInt(ss[4], 16, 64)
			v1, _ := strconv.ParseInt(ss[5], 16, 64)
			v2, _ := strconv.ParseInt(ss[6], 16, 64)
			v3, _ := strconv.ParseInt(ss[7], 16, 64)
			v4, _ := strconv.ParseInt(ss[8], 16, 64)
			v5, _ := strconv.ParseInt(ss[9], 16, 64)
			v6, _ := strconv.ParseInt(ss[10], 16, 64)
			v7, _ := strconv.ParseInt(ss[11], 16, 64)
			if toErr == nil && didErr == nil && err0 == nil {
				var d [8]uint8 = [8]uint8{uint8(v0), uint8(v1), uint8(v2), uint8(v3), uint8(v4), uint8(v5), uint8(v6), uint8(v7)}
				i, err := Parse(uint64(to), uint32(did), d)
				if err == nil {
					err := websocket.JSON.Send(conn, i)
					if err != nil {
						return
					}
				}
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

// CAN数据解析
type Numeric struct {
	NumOB      int     // Number Of Bytes
	ByteSPos   int     // ByteStart Position
	ScalFac    float64 // Scaling Factor
	ScalOffset float64 // Scaling Offset
	Signed     bool
	Amin       float64 // Absolute Min
	Amax       float64 // Absolute Max
	Omin       float64 // Operating Min
	Omax       float64 // Operating Max
	Units      string
	Trunc      int // 截取小数点后几位
}

func (n *Numeric) Decode(d [8]byte) (float64, error) {
	var r uint32

	for i := 0; i < n.NumOB; i++ {
		r <<= 8
		r |= uint32(d[n.ByteSPos+n.NumOB-1-i])
	}

	r_f := float64(r) * n.ScalFac
	r_f += n.ScalOffset

	pow10_t := math.Pow10(n.Trunc)
	r_f = math.Trunc((r_f+0.5/pow10_t)*pow10_t) / pow10_t

	return r_f, nil
}

type BMS100 struct {
	DID   uint32  // CAN ID
	To    uint64  // Time Offset
	PackU float64 // 电池电压
	PackI float64 // 电池电流
}

func Parse(to uint64, did uint32, d [8]byte) (interface{}, error) {
	if did == 0x100 {
		u, _ := (&Numeric{2, 0, 0.1, 0, false, 0, 0, 0, 0, "V", 1}).Decode(d)
		i, _ := (&Numeric{2, 2, 0.1, -500, false, 0, 0, 0, 0, "A", 1}).Decode(d)
		bms100 := BMS100{
			DID:   did,
			To:    to,
			PackU: u,
			PackI: i,
		}
		return bms100, nil
	}
	return nil, fmt.Errorf("Unknown did: 0x%04X", did)
}
