package init_pkg1

var PackageName string

func init() {
	PackageName = "init_pkg1"
	println(PackageName, "包1不依赖")
}
