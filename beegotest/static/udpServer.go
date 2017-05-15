package main

import (
	"fmt"
	"net"
	"runtime"
	"time"
)

func handleDatagram(conn *net.UDPConn, quit chan struct{}) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		fmt.Printf("read %d bytes\n", n)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			if e, ok := err.(net.Error); !ok || !e.Timeout() {
				// handle error, it's not a timeout
			}
			break
		}

		msg := ""
		for i := 0; i < n; i++ {
			msg += fmt.Sprintf("0x%02X ", buf[i])
		}
		fmt.Printf("read is\n%s\n", msg)

		n, err = conn.WriteToUDP(buf[:n], addr)
		fmt.Printf("write %d bytes\n", n)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}
	quit <- struct{}{}
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":6725")
	if err != nil {
		fmt.Printf("1 %s\n", err.Error())
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("2 %s\n", err.Error())
		return
	}
	defer conn.Close()

	// 设置超时时间
	conn.SetReadDeadline(time.Now().Add(300 * time.Second))

	// 利用多核(主线程占用一核)
	quit := make(chan struct{})
	for i := 0; i < runtime.NumCPU()-1; i++ {
		go handleDatagram(conn, quit)
	}
	<-quit // hang until an error
}
