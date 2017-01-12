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

	// 返回正则表达式中捕获分组的个数（某些分组为不捕获分组）
	reg = regexp.MustCompile(`h(e.*)o wo(\w)d`)
	fmt.Println(reg.NumSubexp()) // 2
	reg = regexp.MustCompile(`h(?:e.*)o w(\w*)d`)
	fmt.Println(reg.NumSubexp()) // 1

	// 返回正则表达式中捕获分组的名字
	reg = regexp.MustCompile(`h(e.*)o w(\w*)d`)
	fmt.Printf("%q\n", reg.SubexpNames()) // ["" "" ""] 注意：此处为3个空字符串
	reg = regexp.MustCompile(`h(?P<first>e.*)o w(?P<second>\w*)d`)
	fmt.Println(reg.SubexpNames()) // [ first second] 注意：此处为3个字符串，第一个字符串为空

	// 让正则表达式采用“leftmost-longest”模式
	// reg.Longest()

	// 检查b中是否存在匹配正则表达式的子序列
	// func (re *Regexp) Match(b []byte) bool
	reg = regexp.MustCompile("he.*")
	fmt.Println(reg.Match([]byte("hello")))                        // true
	fmt.Println(reg.MatchString("hello"))                          // true
	fmt.Println(reg.MatchReader(bytes.NewReader([]byte("hello")))) // true

	// 返回b中匹配正则表达式的第一个子序列
	// func (re *Regexp) Find(b []byte) []byte
	reg = regexp.MustCompile(`he\w*`)                           // heee虽然也匹配正则表达式，但其是第二个子序列
	fmt.Println(string(reg.Find([]byte("hello world, heeee")))) // hello
	fmt.Println(reg.FindString("hello world"))                  // hello

	// 返回b中匹配正则表达式的第一个子序列的索引
	// func (re *Regexp) FindIndex(b []byte) (loc []int)
	fmt.Println(reg.FindIndex([]byte("hello world")))                        // [0 5]
	fmt.Println(reg.FindStringIndex("hello world"))                          // [0 5]
	fmt.Println(reg.FindReaderIndex(bytes.NewReader([]byte("hello world")))) // [0 5]

	// 返回b中匹配正则表达式的第一个子序列以及（可能有的）分组匹配的子序列
	// func (re *Regexp) FindSubmatch(b []byte) [][]byte
	reg = regexp.MustCompile(`he(\w*)o w([[:alpha:]]*)d`)                            // 注：ASCII字符族需要用两个中括号
	fmt.Println(reg.FindSubmatch([]byte("hello world")))                             // [[104 101 108 108 111 32 119 111 114 108 100] [108 108] [111 114 108]]
	fmt.Println(reg.FindStringSubmatch("hello world"))                               // [hello world ll orl]
	fmt.Println(reg.FindSubmatchIndex([]byte("hello world")))                        // [0 11 2 4 7 10]
	fmt.Println(reg.FindStringSubmatchIndex("hello world"))                          // [0 11 2 4 7 10]
	fmt.Println(reg.FindReaderSubmatchIndex(bytes.NewReader([]byte("hello world")))) // [0 11 2 4 7 10]

	// 返回b中匹配正则表达式的n个子序列。如果n为-1，则返回所有匹配的子序列
	// func (re *Regexp) FindAll(b []byte, n int) [][]byte
	reg = regexp.MustCompile(`he(\w*)o w([[:alpha:]]*)d`)
	fmt.Println(reg.FindAll([]byte("hello world, heeeeo wwwwd"), -1)) // [[104 101 108 108 111 32 119 111 114 108 100] [104 101 101 101 101 111 32 119 119 119 119 100]]
	fmt.Println(reg.FindAllString("hello world, heeeeo wwwwd", -1))   // [hello world heeeeo wwwwd]

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
