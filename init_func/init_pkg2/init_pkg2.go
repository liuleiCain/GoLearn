package init_pkg2

import "go-learn/init_func/init_pkg1"

var Pkg2 int

func init() {
	Pkg2 = 1
	pkg1Name := init_pkg1.PackageName
	println("包2依赖包1")
	println("init_pkg2: ", "pkg1=", pkg1Name, "pkg2=", Pkg2)
}
