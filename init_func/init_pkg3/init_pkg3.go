package init_pkg1

import "go-learn/init_func/init_pkg2"

var Pkg3Arr string

func init() {
	Pkg3Arr = "pkg3"
	pkg2 := init_pkg2.Pkg2
	println("包3依赖包2")
	println("init_pkg3: ", "pkg2=", pkg2, "pkg3=", Pkg3Arr)
}

func InitPkgSort() {
	println("finish init test")
}
