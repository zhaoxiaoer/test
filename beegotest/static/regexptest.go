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

	// 注意，分组分为：
	// 1. 捕获
	// 1.1 编号的捕获分组
	// 1.2 编号并命名的捕获分组
	// 2. 不捕获
	// 2.1 不捕获，匹配的分组
	// 2.1.1 不捕获，匹配，且设置特定分组的标志的分组
	// 2.2 不捕获，不匹配，且设置当前分组的标志的分组
	// 返回b中匹配正则表达式的第一个子序列以及（可能有的）分组匹配的子序列
	// func (re *Regexp) FindSubmatch(b []byte) [][]byte
	reg = regexp.MustCompile(`he(\w*)o w([[:alpha:]]*)d`)                            // 注：ASCII字符族需要用两个中括号
	fmt.Println(reg.FindSubmatch([]byte("hello world")))                             // [[104 101 108 108 111 32 119 111 114 108 100] [108 108] [111 114 108]]
	fmt.Println(reg.FindStringSubmatch("hello world"))                               // [hello world ll orl]
	fmt.Println(reg.FindSubmatchIndex([]byte("hello world")))                        // [0 11 2 4 7 10]
	fmt.Println(reg.FindStringSubmatchIndex("hello world"))                          // [0 11 2 4 7 10]
	fmt.Println(reg.FindReaderSubmatchIndex(bytes.NewReader([]byte("hello world")))) // [0 11 2 4 7 10]

	// 返回b中匹配正则表达式的不重叠的n个子序列。如果n为-1，则返回所有匹配的子序列
	// func (re *Regexp) FindAll(b []byte, n int) [][]byte
	reg = regexp.MustCompile(`he(\w*)o w([[:alpha:]]*)d`)
	fmt.Println(reg.FindAll([]byte("hello world, heeeeo wwwwd"), -1))      // [[104 101 108 108 111 32 119 111 114 108 100] [104 101 101 101 101 111 32 119 119 119 119 100]]
	fmt.Println(reg.FindAllString("hello world, heeeeo wwwwd", -1))        // [hello world heeeeo wwwwd]
	fmt.Println(reg.FindAllIndex([]byte("hello world, heeeeo wwwwd"), -1)) // [[0 11] [13 25]]
	fmt.Println(reg.FindAllStringIndex("hello world, heeeeo wwwwd", -1))   // [[0 11] [13 25]]

	// 匹配0到N次
	reg = regexp.MustCompile("a*")
	// a*匹配0到N个a，因此a*可以匹配空字符。默认情况下，采用贪婪模式。
	// 字符串的两端分别有个空字符，所以“abkoaabaccadjjaaae”实际上是"空abkoaabaccadjjaaae空"
	// 针对“空abkoaabaccadjjaaae空”:
	// 1 首先从头开始判断字节序列
	//   当判断完b字符时，即已判断的字节序列为“空ab”时，可以确认第1个匹配的子序列为“空a”，此时a匹配一次，即“a"
	// 2 剩余的字节序列为“koaabaccadjjaaae空”,头部加上“空",即”空koaabaccadjjaaae空“
	//   当判断完k字符时，即已判断的字节序列为“空k”时，可以确认第2个匹配的子序列为“空”，此时a匹配0次，即“”
	// 3 剩余的字节序列为“oaabaccadjjaaae空”,头部加上“空",即”空oaabaccadjjaaae空“
	//   当判断完o字符时，即已判断的字节序列为“空o”时，可以确认第3个匹配的子序列为“空”，此时a匹配0次，即“”
	// 4 剩余的字节序列为“aabaccadjjaaae空”,头部加上“空",即”空aabaccadjjaaae空“
	//   当判断完b字符时，即已判断的字节序列为“空aab”时，可以确认第4个匹配的子序列为“空aa”，此时a匹配2次，即“aa”
	// 5 剩余的字节序列为“accadjjaaae空”,头部加上“空",即”空accadjjaaae空“
	//   当判断完c字符时，即已判断的字节序列为“空ac”时，可以确认第5个匹配的子序列为“空a”，此时a匹配1次，即“a”
	// 6 剩余的字节序列为“cadjjaaae空”,头部加上“空",即”空cadjjaaae空“
	//   当判断完c字符时，即已判断的字节序列为“空c”时，可以确认第6个匹配的子序列为“空”，此时a匹配0次，即“”
	// 7 剩余的字节序列为“adjjaaae空”,头部加上“空",即”空adjjaaae空“
	//   当判断完d字符时，即已判断的字节序列为“空ad”时，可以确认第7个匹配的子序列为“空a”，此时a匹配1次，即“a”
	// 8 剩余的字节序列为“jjaaae空”,头部加上“空",即”空jjaaae空“
	//   当判断完j字符时，即已判断的字节序列为“空j”时，可以确认第8个匹配的子序列为“空”，此时a匹配0次，即“”
	// 9 剩余的字节序列为“jaaae空”,头部加上“空",即”空jaaae空“
	//   当判断完j字符时，即已判断的字节序列为“空j”时，可以确认第9个匹配的子序列为“空”，此时a匹配0次，即“”
	// 10 剩余的字节序列为“aaae空”,头部加上“空",即”空aaae空“
	//   当判断完e字符时，即已判断的字节序列为“空aaae”时，可以确认第10个匹配的子序列为“空aaa”，此时a匹配3次，即“aaa”
	// 11 剩余的字节序列为“空”
	//    可以确认第11个匹配的子序列为“空”，此时a匹配0次，即“”
	// 判断结束。
	fmt.Printf("%q\n", reg.FindAllString("abkoaabaccadjjaaae", -1)) // ["a" "" "" "aa" "a" "" "a" "" "" "aaa" ""]
	fmt.Printf("%q\n", reg.FindAllString("bkoaabaccadjjaaae", -1))  // ["" "" "" "aa" "a" "" "a" "" "" "aaa" ""]
	reg = regexp.MustCompile("a+")
	fmt.Printf("%q\n", reg.FindAllString("abkoaabaccadjjaaae", -1)) // ["a" "aa" "a" "a" "aaa"]
	// a{2,}匹配2到N个a。默认情况下，采用贪婪模式。
	// 字符串的两端分别有个空字符，所以“abkoaabaccadjjaaae”实际上是"空abkoaabaccadjjaaae空"
	// 针对“空abkoaabaccadjjaaae空”:
	// 1 首先从头开始判断字节序列
	//   当已判断的字节序列为“空abkoaab”时，可以确认第1个匹配的子序列为“aa”，此时a匹配2次，即“aa"
	// 2 剩余的字节序列为“accadjjaaae空”,头部加上“空",即”空accadjjaaae空“
	//   当已判断的字节序列为“空accadjjaaae”时，可以确认第2个匹配的子序列为“aaa”，此时a匹配3次，即“aaa”
	// 3 剩余的字节序列为“空”
	//   不符合正则表达式。
	// 判断结束。
	reg = regexp.MustCompile("a{2,}")
	fmt.Printf("%q\n", reg.FindAllString("abkoaabaccadjjaaae", -1)) // ["aa" "aaa" ]

	// 返回b中匹配正则表达式的不重叠的n个子序列以及（可能有的）分组匹配的子序列
	// func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte
	reg = regexp.MustCompile(`he(\w*)o w([[:alpha:]]*)d`)
	fmt.Printf("%q\n", reg.FindAllSubmatch([]byte("hello world, heeeeo wwwwd"), -1)) // [["hello world" "ll" "orl"] ["heeeeo wwwwd" "eee" "www"]]
	fmt.Println(reg.FindAllStringSubmatch("hello world, heeeeo wwwwd", -1))          // [[hello world ll orl] [heeeeo wwwwd eee www]]
	fmt.Println(reg.FindAllSubmatchIndex([]byte("hello world, heeeeo wwwwd"), -1))   // [[0 11 2 4 7 10] [13 25 15 18 21 24]]
	fmt.Println(reg.FindAllStringSubmatchIndex("hello world, heeeeo wwwwd", -1))     // [[0 11 2 4 7 10] [13 25 15 18 21 24]]

	// 将s中匹配正则表达式的不重叠的子序列作为分隔符，将s分割成n个字符串
	// func (re *Regexp) Split(s string, n int) []string
	reg = regexp.MustCompile("a*")                       // 匹配0到N次
	fmt.Printf("%q\n", reg.Split("abaabaccadddaaae", 5)) // ["" "b" "b" "c" "cadaaae"]
	// 首先将输入字符串根据正则表达式拆分为子字符串，然后替换即可
	// abaabaccadddaaae 根据判断步骤，拆分为下面的子字符串
	// 空ab 空aab 空ac 空c 空ad 空d 空d 空aaae 空
	fmt.Printf("%q\n", reg.Split("abaabaccadddaaae", -1)) // ["" "b" "b" "c" "c" "d" "d" "d" "e"] 注意最开始的空字符串

	// Expand 将模板处理后，添加到dst后面。Expand将模板中的变量替换为从src中正则匹配的结果，src正则匹配的起始与结束位置由match指定
	// func (re *Regexp) Expand(dst []byte, template []byte, src []byte, match []int) []byte
	reg = regexp.MustCompile(`he(?P<hello>\w+)o w(?P<world>\w+)d`)
	dst := []byte("say: ")
	tpl := []byte("你好$hello，世界$2") // 注意：不能使用template语法{{.hello}}
	src := []byte("aahello worldbb")
	match := reg.FindSubmatchIndex(src)
	fmt.Println(string(reg.Expand(dst, tpl, src, match))) // say: 你好ll，世界orl
	// ExpandString类似

	// 将src中所有的匹配结果替换为repl
	// func (re *Regexp) ReplaceAllLiteral(src, repl []byte) []byte
	reg = regexp.MustCompile(`he(?P<hello>\w+)o w(?P<world>\w+)d`)
	fmt.Println(string(reg.ReplaceAllLiteral([]byte("aahello worldbb eehello worldff"), []byte("cc")))) // aaccbb eeccff
	// ReplaceAllLiteralString 类似

	// 将src中所有的匹配结果替换为repl。在替换时，repl中的`$`符号会按照Expand方法的规则进行解释和替换，例如$1会被替换为第一个分组匹配的结果
}
