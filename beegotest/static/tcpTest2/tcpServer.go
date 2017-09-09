// Channel关闭原则：不要在接收端关闭channel，也不要在多个并发发送端中关闭channel。
// 注意：1 关闭已经关闭的channel会导致panic
//      2 发送值到已经关闭的channel会导致panic
// 解决方案1---打破“channel关闭原则”：defer、sync.Once、sync.Mutex
// 解决方案2---保持“channel关闭原则”：
//       (1) M个receivers，一个sender，sender通过关闭data channel说“不再发送”
//       (2) 一个receiver，N个senders，receiver通过关闭一个额外的signal channel说“请停止发送”
//       (3) M个receivers，N个senders，它们当中任意一个通过通知一个moderator（仲裁者）关闭额外的
//           signal channel来说“让我们结束游戏吧”
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Callback 是一个回调接口，用于连接的各种事件处理
type Callback interface {
	// 链接建立回调
	OnConnected(bc *BConn)
	// 消息处理回调
	OnMessage(bc *BConn, p []byte)
	// 链接断开回调
	OnDisconnected(bc *BConn)
	// 错误回调
	OnError(err error, bc *BConn)
}

type BCallback struct {
}

func (bcb *BCallback) OnConnected(bc *BConn) {
	fmt.Printf("OnConnected, connected\n")
}

func (bcb *BCallback) OnMessage(bc *BConn, p []byte) {
	fmt.Printf("OnMessage, message: %d\n", len(p))
}

func (bcb *BCallback) OnDisconnected(bc *BConn) {
	fmt.Printf("OnDisconnected, disconnected\n")
}

func (bcb *BCallback) OnError(err error, bc *BConn) {
	fmt.Printf("OnError, err: %v\n", err)
}

type BConn struct {
	conn      net.Conn
	writeChan chan []byte
	exitChan  chan struct{}
	rlQuit    chan struct{}
	wlQuit    chan struct{}
	closeOnce sync.Once
}

func NewBConn(conn net.Conn) *BConn {
	return &BConn{
		conn:      conn,
		writeChan: make(chan []byte),
		exitChan:  make(chan struct{}),
		rlQuit:    make(chan struct{}),
		wlQuit:    make(chan struct{}),
	}
}

// 当前这种写法，“已连接”的消息由 clientManager goroutine 发送，
// “已断开”的消息由 clientManager、readLoop、writeLoop 三个 goroutine 中的其中
// 一个发送（当 server 主动关闭连接时，“已断开”的消息由 clientManager 发送，
// 当 client 主动关闭连接时，“已断开”的消息由 readLoop 或 writeLoop 发送
// 这样就可能出现“已连接”消息和”已断开“消息不是由同一个 goroutine 发送的情况
func (bc *BConn) Serve() {
	//	obd.events <- "已连接"
	go bc.readLoop()
	go bc.writeLoop()
}

func (bc *BConn) readLoop() {
	defer bc.Close()
	defer close(bc.rlQuit)

	for {
		buf := make([]byte, 8)
		n, err := bc.conn.Read(buf)
		if err != nil {
			//			bc.callback.OnError(err, bc)
			fmt.Printf("err: %v\n", err)
			break
		}
		str := ""
		for i := 0; i < n; i++ {
			str += fmt.Sprintf("0x%02X ", buf[i])
		}
		fmt.Println(str)
		//		bc.callback.OnMessage(bc, buf[:n])
		//		obd.events <- "收到数据"
	}
}

func (bc *BConn) writeLoop() {
	defer bc.Close()
	defer close(bc.wlQuit)

	for {
		select {
		case data, ok := <-bc.writeChan:
			if !ok {
				fmt.Printf("writeChan has closed\n")
				return
			}
			n, err := bc.conn.Write(data)
			if err != nil {
				fmt.Printf("write err: %v\n", err)
				return
			}
			if n != len(data) {
				fmt.Printf("write err: n = %d, need = %d\n", n, len(data))
				return
			}
		}
	}
}

func (bc *BConn) Close() {
	fmt.Printf("close begin\n")
	bc.closeOnce.Do(func() {
		fmt.Printf("close begin 1\n")
		//		bc.callback.OnDisconnected(bc)
		//		obd.events <- "已断开"
		bc.conn.Close() // 关闭链接

		// 等待读循环goroutine退出
		<-bc.rlQuit

		// 等待写循环goroutine退出
		close(bc.writeChan)
		<-bc.wlQuit

		// 标志该链接已完全退出
		close(bc.exitChan)
	})
	fmt.Printf("close end\n")
}

func (bc *BConn) Write(b []byte) error {
	str := ""

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("adfasdfasdfasdf\n")
			str = err.(string)
		}
	}()

	select {
	case bc.writeChan <- b:
		fmt.Printf("11111111111\n")
		return nil
	default:
		fmt.Printf("222222222222\n")
		return fmt.Errorf("writeChan is full")
	}

	fmt.Printf("33333333333\n")
	return fmt.Errorf(str)
}

func (bc *BConn) IsClosed() bool {
	select {
	case <-bc.exitChan:
		return true
	default:
		return false
	}
}

type OBD struct {
	commands chan string
	events   chan string

	// 命令goroutine通过该channel向clientManager goroutine发送"发往客户端的数据"
	cData chan []byte

	listener net.Listener
}

func NewOBD() *OBD {
	obd := &OBD{
		events:   make(chan string),
		commands: make(chan string),
		cData:    make(chan []byte),
	}
	go obd.eventPro() // 发送事件的goroutine
	go obd.commPro()  // 接收命令的goroutine

	return obd
}

func (obd *OBD) init() {
	obd.commands <- "initServer\n"
}

func (obd *OBD) uninit() {
	obd.commands <- "uninitServer\n"
}

func (obd *OBD) write() {
	obd.commands <- "write\n"
}

func (obd *OBD) clientManage(cAdd <-chan net.Conn, cmQuit chan<- bool) {
	// 当前只允许存在一个连接
	var bc *BConn

	for {
		select {
		case conn, ok := <-cAdd:
			if !ok {
				fmt.Printf("cAdd channel has been closed\n")
				if bc != nil {
					bc.Close()
				}
				fmt.Printf("22222\n")
				close(cmQuit)
				return
			} else {
				// 为了提高效率，用新连接代替老连接
				if bc != nil {
					bc.Close()
				}

				bc = NewBConn(conn)
				bc.Serve()
			}
		case <-time.After(100 * time.Millisecond):
			// 没办法通过channel来删除conn，所以只能采用主动查询的方式
			if (bc != nil) && (bc.IsClosed()) {
				fmt.Printf("连接已断开\n")
				bc = nil
			}
		case d := <-obd.cData:
			if bc != nil {
				bc.Write(d)
			} else {
				obd.events <- "未连接，无法发送数据"
			}
		}
	}
}

func (obd *OBD) startAccept(lClose <-chan bool, aQuit chan<- bool) {
	cAdd := make(chan net.Conn) // 添加client
	cmQuit := make(chan bool)   // clientManage 退出标志

	// 增加客户端管理goroutine
	go obd.clientManage(cAdd, cmQuit)

	for {
		conn, err := obd.listener.Accept()
		if err != nil {
			fmt.Println(err)
			select {
			case <-lClose:
				fmt.Println("lClose channel has closed")

				// 关闭客户端添加channel
				close(cAdd)
				// 等待clientManage退出
				<-cmQuit

				close(aQuit)
				return
			default:
				continue
			}
		}

		cAdd <- conn
	}
}

// 使用goroutine接收命令，可以使命令序列化，避免竞态
func (obd *OBD) commPro() {
	var err error
	initialized := false

	var lClose chan bool // listener 关闭标志
	var aQuit chan bool  // accept 循环退出标志

	for comm := range obd.commands {
		if comm == "initServer\n" {
			if !initialized {
				fmt.Println("initServer1")

				obd.listener, err = net.Listen("tcp", ":6725")
				if err != nil {
					fmt.Println(err)
					obd.events <- err.Error()
					return
				}

				lClose = make(chan bool)
				aQuit = make(chan bool)
				go obd.startAccept(lClose, aQuit)

				initialized = true
				fmt.Println("initServer2")
			} else {
				obd.events <- "server has initialized"
			}
		} else if comm == "uninitServer\n" {
			if initialized {
				fmt.Println("uninitServer1")

				close(lClose)
				obd.listener.Close()

				// 等待 accept 循环退出
				<-aQuit

				initialized = false
				fmt.Println("uninitServer2")
			} else {
				obd.events <- "server has uninitialized"
			}
		} else if comm == "write\n" {
			if initialized {
				data := []byte{0x01, 0x02, 0x0A}
				obd.cData <- data
			} else {
				obd.events <- "server has uninitialized"
			}
		}
	}
}

func (obd *OBD) eventPro() {
	for event := range obd.events {
		fmt.Println(event)
	}
}

func main() {
	obd := NewOBD()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}

		if line == "init\n" {
			obd.init()
		} else if line == "uninit\n" {
			obd.uninit()
		} else if line == "write\n" {
			obd.write()
		}
	}
}
