package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:6725")
	if err != nil {
		fmt.Printf("1 %s\n", err.Error())
		return
	}
	defer conn.Close()
	n, err := conn.Write([]byte{0x17, 0x5C, 0xCF, 0x7F, 0x80, 0xDE, 0x98, 0xC0, 0xA8, 0x1F, 0x65})
	fmt.Printf("write %d bytes\n", n)
	if err != nil {
		fmt.Printf("2 %s\n", err.Error())
	}

	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	fmt.Printf("read %d bytes\n", n)
	if err != nil {
		fmt.Printf("3 %s\n", err.Error())
	}

	msg := ""
	for i := 0; i < n; i++ {
		msg += fmt.Sprintf("0x%02X ", buf[i])
	}
	fmt.Printf("msg is\n%s\n", msg)
}
