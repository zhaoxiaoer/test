package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":6725")
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	//	defer conn.Close() // 关闭连接

	//	who := conn.RemoteAddr().String()

	go func(conn net.Conn) {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				break
			}
			//			str := who + ": "
			str := ""
			for i := 0; i < n; i++ {
				str += fmt.Sprintf("0x%02X ", buf[i])
			}
			fmt.Println(str)
		}
	}(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("33\n")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Printf("44\n")
		fmt.Print(line)

		n, err := conn.Write([]byte(line))
		fmt.Println(n, err)
		if err != nil {
			fmt.Printf("55\n")
			break
		}
	}
	fmt.Printf("66\n")
}
