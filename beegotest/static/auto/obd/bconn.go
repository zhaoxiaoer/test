package obd

import (
	"fmt"
	"net"
	"sync"
)

type BConn struct {
	conn      net.Conn
	name      string
	msg       chan Event
	writeChan chan []byte
	exitChan  chan struct{}
	rlQuit    chan struct{}
	wlQuit    chan struct{}
	closeOnce sync.Once
}

func NewBConn(conn net.Conn, name string, messages chan Event) *BConn {
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
// 一个发送（当 server 主动关闭连接时，“已断开”的消息由 connManage 发送，
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
		buf := make([]byte, 1024)
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
		bc.msg <- Event{1, bc.name, buf[:n]}
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

// 多个goroutine调用
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

	// 确保readLoop和writeLoop都退出
	// 如果没有 <-bc.exitChan，clientManager goroutine 先退出，readLoop 后退出的话，
	// readLoop 有可能阻塞在 bc.msg <- event{eType, bc.name, buf[:n]}
	<-bc.exitChan

	fmt.Printf("close end\n")
}

func (bc *BConn) Write(b []byte) (err error) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic: %v\n", p)
			err = fmt.Errorf("%v", p)
		}
	}()

	select {
	case bc.writeChan <- b:
		return nil
	default:
		return fmt.Errorf("writeChan is full")
	}

	return
}

func (bc *BConn) IsClosed() bool {
	select {
	case <-bc.exitChan:
		return true
	default:
		return false
	}
}
