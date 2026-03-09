package init_pkg4

var Package4Name string

func init() {
	Package4Name = "init_pkg4"
	println(Package4Name, "包4不依赖")
}
