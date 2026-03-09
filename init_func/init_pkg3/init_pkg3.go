package init_pkg1

import (
	"go-learn/init_func/init_pkg0"
	"go-learn/init_func/init_pkg2"
	"go-learn/init_func/init_pkg4"
)

var Pkg3Arr string

/*
对同一个go文件的 init() 调用顺序是从上到下的。
对同一个package中不同文件是按文件名字符串比较“从小到大”顺序调用各文件中的 init() 函数。
对于不同的 package ，如果不相互依赖的话，按照main包中"先 import 的后调用"的顺序调用其包中的init()
如果 package 存在依赖，则先调用最早被依赖的 package 中的 init() 。
*/
func init() {
	Pkg3Arr = "pkg3"
	pkg2 := init_pkg2.Pkg2
	pkg4 := init_pkg4.Package4Name
	pkg0 := init_pkg0.Package0Name
	println("包3依赖包2")
	println("init_pkg3: ", "pkg2=", pkg2, "pkg3=", Pkg3Arr, "pkg4=", pkg4, "pkg0=", pkg0)
}

func InitPkgSort() {
	println("finish init test")
}
