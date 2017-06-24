package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strconv"
	"sync"
	"syscall"
)

const (
	FDKey = "LISTENER_FD"
)

var wg sync.WaitGroup
var listener *net.TCPListener

type myHandler struct {
}

func (mH *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("begin: %d\n", os.Getpid())
	w.Write([]byte("Hello, world!"))
	fmt.Printf("end\n")
}

// 其实新进程也可以直接 ListenerTCP，但是这样的话，新进程启动前的时间，会拒绝连接
func getListener(addr string, port int) (net.Listener, error) {
	fdStr := os.Getenv(FDKey)
	if fdStr == "" {
		fmt.Printf("new socket\n")
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			fmt.Printf("%v\n", err)
			return nil, err
		}
		return net.ListenTCP("tcp", addr)
	} else {
		fmt.Printf("origin socket")
		fd, err := strconv.Atoi(fdStr)
		if err != nil {
			return nil, err
		}
		f := os.NewFile(uintptr(fd), "listen socket")
		l, err := net.FileListener(f)
		if err != nil {
			return nil, err
		}
		return l, nil
	}
}

type myConn struct {
	net.Conn
}

func (mc myConn) Close() error {
	err := mc.Conn.Close()
	wg.Done()
	fmt.Printf("done\n")
	return err
}

type myListener struct {
	net.Listener
}

func (ml *myListener) Accept() (net.Conn, error) {
	conn, err := ml.Listener.Accept()
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	fmt.Printf("add\n")

	return myConn{Conn: conn}, err
}

func waitSignal() error {
	fmt.Printf("111111111\n")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP)
	fmt.Printf("2222222\n")
	for {
		fmt.Printf("333333\n")
		sig := <-ch
		fmt.Printf("4444444\n")
		fmt.Println(sig.String())
		switch sig {
		case syscall.SIGTERM:
			return nil
		case syscall.SIGHUP:
			return reStart()
		}
	}

	return nil
}

//func closeSelf() error {
//	pid := os.Getpid()
//	p, err := os.FindProcess(pid)
//	if err != nil {
//		return err
//	}
//	return p.Kill()
//}

func reStart() error {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("path: %v\n", path)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Printf("dir: %s\n", dir)

	// listener.File()是底层os.File的拷贝，所以描述符会增加
	//	f, err := listener.File()
	//	fd := f.Fd()
	//	fmt.Printf("fd: %d\n", fd)

	v := reflect.ValueOf(listener).Elem().FieldByName("fd").Elem()
	fd2 := uintptr(v.FieldByName("sysfd").Int())
	fmt.Printf("fd2: %d\n", fd2)

	// fd2可以不传递，如果不传递的话，老进程关闭listen后，新进程开始listener前，服务器将拒绝连接
	p, err := os.StartProcess(path, os.Args, &os.ProcAttr{
		Dir:   dir,
		Env:   append(os.Environ(), fmt.Sprintf("%s=%d", FDKey, fd2)),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr, os.NewFile(fd2, "listen socket")},
	})
	if err != nil {
		return err
	}
	fmt.Printf("spawn child: %d\n", p.Pid)

	return nil
}

func main() {
	server := &http.Server{Handler: &myHandler{}}
	l, err := getListener("", 6725)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	listener = l.(*net.TCPListener)

	// 负责重启的goroutine
	go func() {
		waitSignal()
		l.Close()
	}()

	ml := &myListener{Listener: l}
	server.Serve(ml)

	// 等待已连接的client关闭
	wg.Wait()

	//	// 杀死自己
	//	closeSelf()
	fmt.Printf("server close\n")
}
