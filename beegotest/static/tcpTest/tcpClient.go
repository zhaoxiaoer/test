package main

import (
	"log"
	"net"
	"time"
)

/*
// example1
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6725")
	if err != nil {
		fmt.Printf("dial error: %s\n", err)
		return
	}
	defer conn.Close()

	// read or write on conn
}
*/

/*
// example2
func main() {
	for i := 0; i < 1000; i++ {
		_, err := net.Dial("tcp", ":6725")
		if err != nil {
			fmt.Printf("dial error: %s\n", err)
			return
		}
		fmt.Printf("%d: connect to server ok\n", i)
	}
}

func main() {
	conn, err := net.DialTimeout("tcp", "www.baidu.com:80", 2*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Printf("dial ok\n")
}
*/

/*
// example3
func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("usage: go run client.go YOUR_CONTENT\n")
		return
	}

	c, err := net.Dial("tcp", ":6725")
	if err != nil {
		fmt.Printf("dial error: %s\n", err)
		return
	}
	defer c.Close()

	data := os.Args[1]
	c.Write([]byte(data))
}
*/

/*
// example4
func main() {
	c, err := net.Dial("tcp", ":6725")
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer c.Close()

	fmt.Println("dial ok")

	data := make([]byte, 65536)
	c.Write(data)
	fmt.Println("write ok")

	time.Sleep(100 * time.Second)
}
*/

/*
// example5 模拟写阻塞
func main() {
	c, err := net.Dial("tcp", ":6725")
	if err != nil {
		log.Printf("dial error: %v\n", err)
		return
	}
	defer c.Close()

	log.Println("dial ok")

	data := make([]byte, 65535)
	var total int
	for {
		n, err := c.Write(data)
		total += n
		if err != nil {
			log.Printf("write %d bytes, error: %v\n", n, err)
			break
		}
		log.Printf("write %d bytes this time, %d bytes in totaln", n, total)
	}
	log.Printf("write %d bytes in total\n", total)

	time.Sleep(100 * time.Second)
}
*/

/*
// example6 模拟写超时
func main() {
	c, err := net.Dial("tcp", ":6725")
	if err != nil {
		log.Printf("dial error: %v\n", err)
		return
	}
	defer c.Close()

	log.Println("dial ok")

	data := make([]byte, 65535)
	var total int
	for {
		c.SetWriteDeadline(time.Now().Add(10 * time.Microsecond))
		n, err := c.Write(data)
		total += n
		if err != nil {
			log.Printf("write %d bytes, error: %v\n", n, err)
			break
		}
		log.Printf("write %d bytes this time, %d bytes in totaln", n, total)
	}
	log.Printf("write %d bytes in total\n", total)

	time.Sleep(100 * time.Second)
}
*/

// example7 关闭连接
// 在己方已经关闭的socket上再进行read和write操作，会
// 得到“use of closed network connection”
func main() {
	c, err := net.Dial("tcp", ":6725")
	if err != nil {
		log.Printf("dial error: %v\n", err)
		return
	}
	c.Close() // 关闭连接
	log.Println("close ok")

	buf := make([]byte, 32)
	n, err := c.Read(buf)
	if err != nil {
		log.Println("read error: ", err)
	} else {
		log.Printf("read %d bytes", n)
	}

	n, err = c.Write(buf)
	if err != nil {
		log.Println("write error: ", err)
	} else {
		log.Printf("write %d bytes\n", n)
	}

	time.Sleep(10 * time.Second)
}
