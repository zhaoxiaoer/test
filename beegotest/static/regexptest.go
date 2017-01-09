package main

import (
	"fmt"
	"regexp"
)

func main() {
	s := "hello, world"

	// 将字符串s中的所有正则表达式元字符进行转义
	// func QuoteMeta(s string) string
	fmt.Println(regexp.QuoteMeta(`[foo]+\d`)) // \[foo\]\+\\d

	// 检查b中是否存在匹配正则表达式pattern的子序列
	// func Match(pattern string, b []byte) (matched bool, err error)
	//	fmt.Println(regexp.Match("(?i:He.*)", []byte("hello, world")))
	fmt.Println(regexp.Match("He.*", []byte(s))) // false
	fmt.Println(regexp.Match("he.*", []byte(s))) // true

	//
	r, err := regexp.Compile("(?i:He.*),") // 匹配但不捕获的分组
	if err != nil {
		fmt.Println(err)
		return
	}
	//	fmt.Println(r.FindString("hello, world\n"))
	fmt.Printf("%v\n", r.FindStringSubmatch("hello3, world\n"))

	reg := regexp.MustCompile(`(\w)(\w)+`)
	fmt.Printf("%q\n", reg.FindSubmatch([]byte("Hello World!"))) // ["Hello" "H" "o"]
	reg = regexp.MustCompile(`(?:\w)(\w)(?U)+`)
	fmt.Printf("%q\n", reg.FindSubmatch([]byte("Hello World!"))) // ["He" "H" "o"]
}
