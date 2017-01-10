package main

import (
	"bytes"
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
	fmt.Println(regexp.Match("He.*", []byte(s))) // false <nil>
	fmt.Println(regexp.Match("he.*", []byte(s))) // true <nil>

	fmt.Println(regexp.MatchString("He.*", s)) // false <nil>
	fmt.Println(regexp.MatchString("he.*", s)) // true <nil>

	fmt.Println(regexp.MatchReader("he.*", bytes.NewReader([]byte(s)))) // true <nil>

	// 解析并返回一个正则表达式。如果成功返回，该Regexp就可以用于匹配文本
	// 在匹配文本时，采用“leftmost-first”模式，还有一种模式为“leftmost-longest”
	// func Compile(expr string) (*Regexp, error)
	reg, _ := regexp.Compile("he.*")
	fmt.Println(reg.MatchString("hello, world")) // true
	// CompilePOSIX采用"leftmost-longest"
	//	reg, _ := regexp.CompilePOSIX("he.*")
	//	reg := regexp.MustCompile("he.*")
	//	reg := regexp.MustCompilePOSIX("he.*")

	// 返回用来编译出正则表达式的pattern字符串
	fmt.Printf("%s\n", reg.String()) // he.*

	// 返回一个字符串字面值，任何匹配本正则表达式的字符串都会以
	// 该字面值起始。如果该字符串字面值包含整个正则表达式，返回值
	// complete会设置为真
	fmt.Println(reg.LiteralPrefix()) // he false
	r, _ := regexp.Compile("he")
	fmt.Println(r.LiteralPrefix()) // he true
	r, _ = regexp.Compile(".*he")
	fmt.Println(r.LiteralPrefix()) // "" false

	//
	//	r, err := regexp.Compile("(?i:He.*),") // 匹配但不捕获的分组
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	//	fmt.Println(r.FindString("hello, world\n"))
	//	fmt.Printf("%v\n", r.FindStringSubmatch("hello3, world\n"))

	//	reg := regexp.MustCompile(`(\w)(\w)+`)
	//	fmt.Printf("%q\n", reg.FindSubmatch([]byte("Hello World!"))) // ["Hello" "H" "o"]
	//	reg = regexp.MustCompile(`(?:\w)(\w)(?U)+`)
	//	fmt.Printf("%q\n", reg.FindSubmatch([]byte("Hello World!"))) // ["He" "H" "o"]
}
