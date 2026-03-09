package init_pkg0

var Package0Name string

func init() {
	Package0Name = "init_pkg0"
	println(Package0Name, "包0不依赖")
}
