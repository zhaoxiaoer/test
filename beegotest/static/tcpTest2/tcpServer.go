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
	"strings"
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
	name      string
	msg       chan message
	writeChan chan []byte
	exitChan  chan struct{}
	rlQuit    chan struct{}
	wlQuit    chan struct{}
	closeOnce sync.Once
}

func NewBConn(conn net.Conn, name string, messages chan message) *BConn {
	return &BConn{
		conn:      conn,
		name:      name,
		msg:       messages,
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
		//		str := bc.conn.RemoteAddr().String() + ": "
		//		for i := 0; i < n; i++ {
		//			str += fmt.Sprintf("0x%02X ", buf[i])
		//		}
		//		fmt.Println(str)
		//		bc.callback.OnMessage(bc, buf[:n])
		//		obd.events <- "收到数据"
		bc.msg <- message{bc.name, buf[0:n]}
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
		return nil
	default:
		return fmt.Errorf("writeChan is full")
	}

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

type message struct {
	sender string
	value  []byte
}

type event struct {
	// 事件类型
	eType int
	// 事件描述字符串
	eDesc string
	// 事件附加数据
	eOptVal []byte
}

type OBD struct {
	commands chan event
	events   chan event

	// 命令goroutine通过该channel向clientManager goroutine发送"发往客户端的数据"
	messages chan message

	listener net.Listener
}

func NewOBD() *OBD {
	obd := &OBD{
		events:   make(chan event),
		commands: make(chan event),
		messages: make(chan message),
	}
	go obd.eventPro() // 发送事件的goroutine
	go obd.commPro()  // 接收命令的goroutine

	return obd
}

func (obd *OBD) init() {
	obd.commands <- event{1, "initServer\n", nil}
}

func (obd *OBD) uninit() {
	obd.commands <- event{9, "uninitServer\n", nil}
}

func (obd *OBD) write(data []byte) {
	obd.commands <- event{2, "write\n", data}
}

func (obd *OBD) clientManage(cAdd <-chan net.Conn, cmQuit chan<- bool) {
	bcs := make(map[*BConn]string)

	for {
		select {
		case conn, ok := <-cAdd:
			if !ok {
				fmt.Printf("cAdd channel has been closed\n")

				for bc := range bcs {
					bc.Close()
					// hooker的打开与关闭不处理
					if bc.name == "client" {
						obd.events <- event{-13, "client已关闭", nil}
					}
				}

				fmt.Printf("22222\n")
				close(cmQuit)
				return
			} else {
				// 将第一个连接当成client
				// 随后的连接当成hooker
				name := ""
				if len(bcs) == 0 {
					name = "client"
				} else {
					name = "hooker"
				}
				bc := NewBConn(conn, name, obd.messages)
				bc.Serve()
				bcs[bc] = name

				if bc.name == "client" {
					obd.events <- event{-14, "client已连接", nil}
				}
			}
		case <-time.After(100 * time.Millisecond):
			// 没办法通过channel来删除conn，所以只能采用主动查询的方式
			for bc := range bcs {
				if bc.IsClosed() {
					delete(bcs, bc)

					if bc.name == "client" {
						obd.events <- event{-13, "client已关闭", nil}
					}
				}
			}
		case msg := <-obd.messages:
			if msg.sender == "server" {
				// server发送的消息将发送给client和所有hooker
				for bc := range bcs {
					bc.Write(msg.value)
				}
			} else if msg.sender == "client" {
				// client发送的消息将发送给server和所有hooker
				// 发送给server
				obd.events <- event{-15, "收到client数据", msg.value}
				for bc, name := range bcs {
					// 发送给所有hooker
					if name == "hooker" {
						bc.Write(msg.value)
					}
				}
			} else {
				// hooker发送的消息将根据数据的第一个字节来确定
				// 应该将剩余数据发送给server还是client
				if msg.value[0] == 0x30 {
					obd.events <- event{-16, "收到hooker数据", msg.value[1:]}
				} else {
					for bc, name := range bcs {
						// 仅发送给client
						if name == "client" {
							bc.Write(msg.value[1:])
							break
						}
					}
				}
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

	for evt := range obd.commands {
		if evt.eType == 1 {
			if !initialized {
				fmt.Println("initServer1")

				obd.listener, err = net.Listen("tcp", ":6725")
				if err != nil {
					fmt.Println(err)
					obd.events <- event{-3, err.Error(), nil}
					return
				}

				lClose = make(chan bool)
				aQuit = make(chan bool)
				go obd.startAccept(lClose, aQuit)

				initialized = true
				fmt.Println("initServer2")
			} else {
				obd.events <- event{-1, "server has initialized", nil}
			}
		} else if evt.eType == 9 {
			if initialized {
				fmt.Println("uninitServer1")

				close(lClose)
				obd.listener.Close()

				// 等待 accept 循环退出
				<-aQuit

				initialized = false
				fmt.Println("uninitServer2")
			} else {
				obd.events <- event{-2, "server has uninitialized", nil}
			}
		} else if evt.eType == 2 {
			if initialized {
				obd.messages <- message{"server", evt.eOptVal}
			} else {
				obd.events <- event{-2, "server has uninitialized", nil}
			}
		}
	}
}

func (obd *OBD) eventPro() {
	for event := range obd.events {
		str := ""
		if event.eOptVal != nil {
			for i := 0; i < len(event.eOptVal); i++ {
				str += fmt.Sprintf("0x%02X ", event.eOptVal[i])
			}
		}
		fmt.Println(event.eType, event.eDesc, str)
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

		if line == "initServer\n" {
			obd.init()
		} else if line == "uninitServer\n" {
			obd.uninit()
		} else if strings.HasPrefix(line, "write ") {
			bs := []byte(line)
			obd.write(bs[6:])
		}
	}
}
