package obd

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type event struct {
	// 事件类型
	eType int
	// 事件描述字符串
	eDesc string
	// 事件附加数据
	eOptVal []byte
}

type Event struct {
	// 事件类型
	EType int
	// 事件描述字符串
	EDesc string
	// 事件附加数据
	EOptVal []byte
}

type TCPServer struct {
	OutEvents chan Event

	sync.RWMutex
	hasInited bool

	connAdd      chan net.Conn
	cmQuit       chan bool
	events       chan event
	writeMsg     chan []byte // 写数据
	getStatusMsg chan int    // 获取client连接状态
	statusEvent  chan bool   // 返回client连接状态

	ln      net.Listener
	lnClose chan bool
	acQuit  chan bool
}

func NewTCPServer() *TCPServer {
	ts := &TCPServer{}
	return ts
}

func (ts *TCPServer) Init() error {
	ts.Lock()
	defer ts.Unlock()

	if !ts.hasInited {
		ts.OutEvents = make(chan Event)

		var err error
		ts.ln, err = net.Listen("tcp", ":6725")
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return err
		}

		ts.connAdd = make(chan net.Conn)
		ts.cmQuit = make(chan bool)
		ts.events = make(chan event)
		ts.writeMsg = make(chan []byte)
		ts.getStatusMsg = make(chan int)
		ts.statusEvent = make(chan bool)
		go ts.connManage()

		ts.lnClose = make(chan bool)
		ts.acQuit = make(chan bool)
		go ts.accept()

		ts.hasInited = true
		return nil
	} else {
		return fmt.Errorf("tcpserver has been initialized!")
	}
}

func (ts *TCPServer) Uninit() error {
	ts.Lock()
	defer ts.Unlock()

	if ts.hasInited {
		close(ts.lnClose)
		err := ts.ln.Close()
		if err != nil {
			return err
		}
		<-ts.acQuit

		// 释放connManage相关资源
		close(ts.connAdd)
		<-ts.cmQuit
		close(ts.events)
		close(ts.writeMsg)
		close(ts.getStatusMsg)
		close(ts.statusEvent)

		close(ts.OutEvents)

		ts.hasInited = false
		return nil
	} else {
		return fmt.Errorf("tcpserver not initialized!")
	}
}

func (ts *TCPServer) Write(data []byte) error {
	ts.RLock()
	defer ts.RUnlock()

	if ts.hasInited {
		ts.writeMsg <- data
		return nil
	} else {
		return fmt.Errorf("tcpserver not initialized!")
	}
}

func (ts *TCPServer) IsConnected() (bool, error) {
	ts.RLock()
	defer ts.RUnlock()

	if ts.hasInited {
		ts.getStatusMsg <- 1
		return <-ts.statusEvent, nil
	} else {
		return false, fmt.Errorf("tcpserver not initialized!")
	}
}

func (ts *TCPServer) accept() {
	for {
		conn, err := ts.ln.Accept()
		if err != nil {
			select {
			case <-ts.lnClose:
				close(ts.acQuit)
				return
			default:
				continue
			}
		}
		ts.connAdd <- conn
	}
}

func (ts *TCPServer) connManage() {
	bcs := make(map[*BConn]string)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case conn, ok := <-ts.connAdd:
			if !ok {
				//				fmt.Printf("connAdd channel has been closed\n")
				for bc := range bcs {
					bc.Close()
					// hooker的打开与关闭不处理
					if bc.name == "client" {
						ts.sendEvent(Event{2, "client未连接", nil})
					}
				}

				close(ts.cmQuit)
				return
			} else {
				// 没有client时将新连接当成client
				// 已存在client时，将新连接当成hooker
				name := ""
				hasClient := false
				for bc := range bcs {
					if bc.name == "client" {
						hasClient = true
						break
					}
				}
				if hasClient {
					name = "hooker"
				} else {
					name = "client"
				}

				// 最多允许3个连接，其中一个为client，两个为hooker
				if (name == "hooker") && (len(bcs) > 2) {
					conn.Close()
					continue
				}

				bc := NewBConn(conn, name, ts.events)
				bc.Serve()
				bcs[bc] = name

				if bc.name == "client" {
					ts.sendEvent(Event{3, "client已连接", nil})
				}
			}
		case <-ticker.C:
			// 没办法通过channel来删除conn，所以只能采用主动查询的方式
			for bc := range bcs {
				if bc.IsClosed() {
					delete(bcs, bc)

					if bc.name == "client" {
						ts.sendEvent(Event{2, "client未连接", nil})
					}
				}
			}
		case <-ts.getStatusMsg:
			findClient := false
			for bc := range bcs {
				if bc.name == "client" {
					findClient = true
				}
			}
			if !findClient {
				ts.statusEvent <- false
			} else {
				ts.statusEvent <- true
			}

		default:
		}

		select {
		case evt := <-ts.events:
			if evt.eType == 1 {
				if evt.eDesc == "client" {
					// 发送给server
					ts.sendEvent(Event{1, "data from the client", evt.eOptVal})

					for bc, name := range bcs {
						// 发送给所有hooker
						if name == "hooker" {
							bc.Write(evt.eOptVal)
						}
					}
				} else {
					// hooker发送的数据将根据数据的第一个字节来确定
					// 应该将剩余数据发送给server还是client
					if evt.eOptVal[0] == 'c' {
						ts.sendEvent(Event{1, "data from the hooker", evt.eOptVal[1:]})
					} else if evt.eOptVal[0] == 's' {
						for bc, name := range bcs {
							// 仅发送给client
							if name == "client" {
								bc.Write(evt.eOptVal[1:])
								break
							}
						}
					} else {
						fmt.Printf("ignore\n")
					}
				}
			} else {
				fmt.Printf("unknown type: %d\n", evt.eType)
			}
		case data := <-ts.writeMsg:
			// server发送的消息将发送给client和所有hooker
			findClient := false
			for bc := range bcs {
				if bc.name == "client" {
					findClient = true
				}
				bc.Write(data)
			}
			if !findClient {
				ts.sendEvent(Event{2, "client未连接", nil})
			}

		default:
		}
	}
}

func (ts *TCPServer) sendEvent(evt Event) error {
	select {
	case ts.OutEvents <- evt:
		return nil
	default:
		return fmt.Errorf("event can not be sent, Type: %d, Desc: %s", evt.EType, evt.EDesc)
	}
}
