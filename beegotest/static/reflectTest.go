// 计算机中提到的反射指的是程序借助某种手段检查自己结构的一种能力，通常借助的是编程语言中定义的各种类型。
// 因此反射是建立在类型系统上的
// Go语言的接口，内部存储了一个pair,即（value, type)，（值，具体类型）
// reflect包中的Type对应接口变量中的type，Value对应接口变量中的value
// 为了保持API简单，Value的setter和getter方法操作的是某个值的最大类型
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float32 = 3
	var ptrX = &x

	xT := reflect.TypeOf(x)
	xV := reflect.ValueOf(x)
	ptrXT := reflect.TypeOf(ptrX)
	ptrXV := reflect.ValueOf(ptrX)

	fmt.Println("type: ", xT)
	fmt.Println("value: ", xV)
	fmt.Println("v CanSet: ", xV.CanSet())
	fmt.Println("type: ", ptrXT)
	fmt.Println("value: ", ptrXV)
	fmt.Println("ptrV CanSet: ", ptrXV.CanSet())

	//	xVE := xV.Elem()
	ptrXVE := ptrXV.Elem()
	//	fmt.Println("xVE CanSet: ", xVE.CanSet())
	fmt.Println("ptrXVE CanSet: ", ptrXVE.CanSet())

	ptrXVE.SetFloat(5)
	fmt.Println("x: ", x)
}
