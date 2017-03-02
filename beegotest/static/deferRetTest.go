// 1. defer是在return之前执行的  2. return x 语句不是原子指令。
// return x语句流程:
// 		1. 赋值指令。返回值 = x
// 		2. defer指令。调用defer函数
// 		3. RET指令。空的return
package main

import (
	"fmt"
)

func test1() (r int) {
	defer func() {
		r++
	}()
	return 0
}

func test2() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

func test3() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

func main() {
	fmt.Println("test1:", test1())
	fmt.Println("test2:", test2())
	fmt.Println("test3:", test3())
}

/*
输出:
test1: 1
test2: 5
test3: 1
*/
