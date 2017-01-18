package main

import (
	"fmt"
	"os"

	"github.com/russross/blackfriday"
)

var s = `# 一级标题

## 二级标题
### 三级标题
#### 四级标题
##### 五级标题
###### 六级标题

### 无序列表
* 1
* 2
* 3

### 有序列表
1. 1
2. 2
3. 3

### 引用
> 这里是引用

### 链接
[baidu](http://www.baidu.com)

### 图片
![1.jpg](1.jpg)

### 粗体
**这里是粗体**

### 斜体

*这里是斜体*

### 表格
|Table|Are|Cool|
|--------------|:-------------:|------:|
|col 3 is      |right-aligned  | $1600 |
|col 2 is      |centered       | $12   |
|zebra stripes |are neat       | $1    |

### 代码框
<code>代码区块</code>

### 分割线
***
`

func main() {
	// 原生markdown不支持表格，扩展的markdown才支持
	//	bOutput := blackfriday.MarkdownBasic([]byte(s))
	bOutput := blackfriday.MarkdownCommon([]byte(s))
	fmt.Printf("%s", string(bOutput))

	f, err := os.OpenFile("test.html", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	n, err := f.WriteString("<html>\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />\n<body>\n")
	n, err = f.Write(bOutput)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("n: %d\n", n)
	_, err = f.WriteString("</body></html>")
}
