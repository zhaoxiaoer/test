// 从tcp socket诞生后，网络编程框架模型也几经演化，大致是：“每进程一个连接”->“每线程一个
// 连接”->“Non-Block + I/O多路复用”。目前主流web server均采用“Non-Block + I/O多路复用”。
// 不过Go的设计者似乎认为I/O多路复用通过回调机制割裂控制流的方式依旧复杂，且有悖于“一般逻辑”的
// 设计，为此Go语言将该“复杂性”隐藏到Runtime中了：Go开发者无需关注socket是否是non-block的，
// 也无需亲自注册文件描述符的回调，只需在每个连接对应的goroutine中以“block I/O”的方式对待
// socket处理即可，这大大降低了开发人员的负担。
// net.Conn是Goroutine safe的。TCPConn可以设置socket属性。golang中，listener socket
// 默认采用了SO_REUSEADDR属性。
package main

import (
	"log"
	"net"
	"time"
)

/*
// example1
func handleConn(c net.Conn) {
	defer c.Close()
	fmt.Printf("%s connected.\n", c.RemoteAddr().String())
	io.Copy(c, c)
	fmt.Printf("%s disconnected.\n", c.RemoteAddr().String())
}

func main() {
	l, err := net.Listen("tcp", ":6725")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}
	defer l.Close()

	for {
		time.Sleep(10 * time.Second)
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		// start a new goroutine to handle
		// the new connection
		go handleConn(c)
	}
}
*/

/*
// example2
func main() {
	l, err := net.Listen("tcp", ":6725")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}
	defer l.Close()

	var i int
	for {
		time.Sleep(10 * time.Second)
		_, err := l.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}
		fmt.Printf("%d: accept a new connection\n", i)
	}
}
*/

/*
// example3
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		var buf = make([]byte, 10)
		n, err := c.Read(buf)
		if err != nil {
			fmt.Printf("conn read error: %s\n", err)
			return
		}
		fmt.Printf("read %d bytes\n", n)
	}
}

func main() {
	l, err := net.Listen("tcp", ":6725")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(c)
	}
}
*/

/*
// example4 模拟读超时
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		time.Sleep(10 * time.Second)
		var buf = make([]byte, 65536)
		fmt.Println("start to read from conn")
		c.SetReadDeadline(time.Now().Add(10 * time.Microsecond))
		n, err := c.Read(buf)
		if err != nil {
			fmt.Printf("conn read %d bytes, error: %s\n", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			return
		}
		fmt.Printf("read %d bytes\n", n)
	}
}

func main() {
	l, err := net.Listen("tcp", ":6725")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(c)
	}
}
*/

/*
// example5
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		time.Sleep(5 * time.Second)
		var buf = make([]byte, 60000)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes, error: %s\n", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			if err == io.EOF {
				log.Printf("EOF\n")
				break
			}
			return
		}
		log.Printf("read %d bytes\n", n)
	}
}

func main() {
	l, err := net.Listen("tcp", ":6725")
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(c)
	}
}
*/

/*
// example6
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		time.Sleep(5 * time.Second)
		var buf = make([]byte, 60000)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes, error: %s\n", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			if err == io.EOF {
				log.Printf("EOF\n")
				break
			}
			return
		}
		log.Printf("read %d bytes\n", n)
	}
}

func main() {
	l, err := net.Listen("tcp", ":6725")
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(c)
	}
}
*/

// example7 关闭连接
// 在对方关闭的socket上执行read操作，会得到EOF error
// 但是write操作会成功，因为己方socket并未关闭，数据会成功写入己方的内
// 核socket缓冲区中，即便最终发送不到对方的socket缓冲区。
// 因此，当发现对方socket关闭后，己方应该正确合理的处理自己的socket，再
// 继续write已经无任何意义了。
func handleConn(c net.Conn) {
	defer c.Close()

	// read from the connection
	time.Sleep(5 * time.Second)
	var buf = make([]byte, 10)
	log.Println("start to read from conn")
	n, err := c.Read(buf)
	if err != nil {
		log.Printf("conn read error: %s\n", err)
	} else {
		log.Printf("conn read %d bytes\n", n)
	}

	n, err = c.Write(buf)
	if err != nil {
		log.Println("conn write error: ", err)
	} else {
		log.Printf("write %d bytes\n", n)
	}
}

func main() {
	l, err := net.Listen("tcp", ":6725")
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(c)
	}
}
