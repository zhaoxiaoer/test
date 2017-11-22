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
	"time"

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

	//	var lastTo int64
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
					//					time.Sleep(time.Duration(to-lastTo) / 1000 * time.Millisecond)
					//					lastTo = to
					time.Sleep(10 * time.Millisecond)
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
	ByteSPos   int     // Byte Start Position
	BitSPos    int     // Bit Start Position
	BitLen     int     // Bit Length
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
	var r uint64
	var m uint64 = 0xFFFFFFFFFFFFFFFF

	var numOB int
	var remainder int
	var bitORS uint
	var maskBitORS uint64

	numOB = (n.BitSPos + n.BitLen) / 8 // 总共需要的字节数
	remainder = (n.BitSPos + n.BitLen) % 8
	if remainder != 0 {
		numOB += 1
		bitORS = uint(8 - remainder) // Bits of right shift
	}
	maskBitORS = uint64(64 - n.BitLen) // 掩码需要右移的位数

	for i := 0; i < numOB; i++ {
		r <<= 8
		r |= uint64(d[n.ByteSPos+i])
	}
	r >>= bitORS
	m >>= maskBitORS
	r &= m

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

type MCU120 struct {
	DID       uint32  // CAN ID
	To        uint64  // Time Offset
	CtlTorque float64 // Controller Torque
}

type VCUP150 struct {
	DID      uint32
	To       uint64
	TorqueRQ float64 // TM电机目标扭矩
}

type VCUP151 struct {
	DID   uint32
	To    uint64
	Speed float64 // Vehicle Speed
}

func Parse(to uint64, did uint32, d [8]byte) (interface{}, error) {
	if did == 0x100 {
		u, _ := (&Numeric{0, 0, 16, 0.1, 0, false, 0, 0, 0, 0, "V", 1}).Decode(d)
		i, _ := (&Numeric{2, 0, 16, 0.1, -500, false, 0, 0, 0, 0, "A", 1}).Decode(d)
		fmt.Printf("u: %v, i: %v\n", u, i)
		bms100 := BMS100{
			DID:   did,
			To:    to,
			PackU: u,
			PackI: i,
		}
		return bms100, nil
	} else if did == 0x120 {
		ctlTorque, _ := (&Numeric{4, 1, 15, 0.05, -300, false, 0, 0, 0, 0, "NM", 2}).Decode(d)
		fmt.Printf("ctlTorque: %v\n", ctlTorque)
		mcu120 := MCU120{
			DID:       did,
			To:        to,
			CtlTorque: ctlTorque,
		}
		return mcu120, nil
	} else if did == 0x150 {
		torqueRQ, _ := (&Numeric{5, 1, 15, 0.05, 0, false, 0, 0, 0, 0, "NM", 2}).Decode(d)
		fmt.Printf("torqueRQ: %v\n", torqueRQ)
		vcup150 := VCUP150{
			DID:      did,
			To:       to,
			TorqueRQ: torqueRQ,
		}
		return vcup150, nil
	} else if did == 0x151 {
		speed, _ := (&Numeric{0, 0, 16, 0.1, 0, false, 0, 0, 0, 0, "KPH", 1}).Decode(d)
		fmt.Printf("speed: %v\n", speed)
		vcup151 := VCUP151{
			DID:   did,
			To:    to,
			Speed: speed,
		}
		return vcup151, nil
	}
	return nil, fmt.Errorf("Unknown did: 0x%04X", did)
}
