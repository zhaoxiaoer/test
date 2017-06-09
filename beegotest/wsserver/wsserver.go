package wsserver

import (
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"golang.org/x/net/websocket"
)

var users = &Users{
	pool: make(map[*websocket.Conn]struct{}),
}

type Users struct {
	mu   sync.RWMutex
	pool map[*websocket.Conn]struct{}
}

func (u *Users) Add(conn *websocket.Conn) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.pool[conn] = struct{}{}
}

func (u *Users) Remove(conn *websocket.Conn) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.pool, conn)
}

func (u *Users) Broadcast(m []byte) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for conn := range u.pool {
		conn.Write(m)
	}
}

func ChatHandler(conn *websocket.Conn) {
	fmt.Printf("11111 chat\n")
	users.Add(conn)
	defer users.Remove(conn)

	fullAddr := conn.RemoteAddr().String()
	begin := strings.LastIndex(fullAddr, "//") + 2
	end := strings.LastIndex(fullAddr, ":")
	addr := fullAddr[begin:end]

	for {
		var msg = make([]byte, 1024)
		n, err := conn.Read(msg)
		if err != nil {
			break
		}
		msg = []byte(fmt.Sprintf("%s: %s", addr, string(msg[:n])))
		users.Broadcast(msg)
	}
	fmt.Printf("22222 chat\n")
}

type Wsserver struct {
	beego.Controller
}

func (wss *Wsserver) Get() {
	wss.TplNames = "wsserver/wsserver.html"
}
