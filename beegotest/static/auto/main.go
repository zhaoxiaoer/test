package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zhaoxiaoer/test/beegotest/static/auto/obd"
)

func main() {
	tcpserver := obd.NewTCPServer()
	tcpserver.Init()
	go func() {
		for {
			select {
			case evt := <-tcpserver.OutEvents:
				//				fmt.Printf("type: %d, desc: %s\n", evt.EType, evt.EDesc)
				if evt.EOptVal != nil {
					str := ""
					for i := 0; i < len(evt.EOptVal); i++ {
						str += fmt.Sprintf("0x%02X ", evt.EOptVal[i])
					}
					fmt.Println(str)
				}
			}
		}
	}()

	// 大数据测试
	for {
		tcpserver.Write([]byte("12345678"))
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}

		//		if line == "init\n" {
		//			tcpserver.Init()
		//		} else if line == "uninit\n" {
		//			tcpserver.Uninit()
		//		} else
		if strings.HasPrefix(line, "write ") {
			bs := []byte(line)
			tcpserver.Write(bs[6 : len(bs)-1])
		} else if strings.HasPrefix(line, "status") {
			status, _ := tcpserver.IsConnected()
			fmt.Printf("status: %t\n", status)
		} else if line == "exit\n" {
			break
		}
	}

	tcpserver.Uninit()
}
