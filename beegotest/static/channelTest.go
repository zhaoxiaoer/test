// go 语言有两种方式用于 goroutine 之间的同步与通信
// 1. 通过 sync 包
// 1.1 Once 和 WaitGroup 主要用于较上层的编程
// 1.2 Mutex, RWMutex, Cond 和 Pool 主要用于较底层的编程
// 2. 通过 channel

// 本例摘抄自：《The Go Programming Language》
// goroutines 使用 channel 来同步与通信
package main

import (
	"fmt"
	"time"
)

func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out) // 写入端关闭
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	// 写 channel 的函数，应该负责关闭 channel
	// 针对本函数，in 是只读 channel，所以不用关闭
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)
	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares) // 0 1 4 9 16 ...
}
