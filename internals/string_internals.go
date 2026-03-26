package internals

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
	"unsafe"
)

// ============================================================
// String 底层实现原理
// ============================================================

// StringHeader 字符串头部结构
type StringHeader struct {
	Data uintptr // 指向字节数组的指针
	Len  int     // 字节长度
}

// DemoStringStructure 演示string的内存结构
func DemoStringStructure() {
	fmt.Println("=== String内存结构 ===")

	s := "Hello世界"

	// 获取字符串头部
	header := (*StringHeader)(unsafe.Pointer(&s))

	fmt.Printf("字符串: %s\n", s)
	fmt.Printf("Data指针: 0x%x\n", header.Data)
	fmt.Printf("字节长度: %d\n", header.Len)
	fmt.Printf("字符数量: %d\n", utf8.RuneCountInString(s))

	fmt.Println("\n内存布局:")
	fmt.Println("+--------+--------+")
	fmt.Println("|  Data  |  Len   |")
	fmt.Printf("| 0x%x |   %d    |\n", header.Data, header.Len)
	fmt.Println("+--------+--------+")
	fmt.Println("         |")
	fmt.Println("         v")
	fmt.Println("+--+--+--+--+--+--+--+--+--+--+--+")
	fmt.Println("|H |e |l |l |o |世       |界       |")
	fmt.Println("+--+--+--+--+--+--+--+--+--+--+--+")
	fmt.Println(" 1  1  1  1  1    3字节    3字节")
	fmt.Println("  <- 5字节 ->  <-  6字节  ->")
	fmt.Println("  总共: 11字节")
}

// DemoStringImmutable 演示string的不可变性
func DemoStringImmutable() {
	fmt.Println("\n=== String不可变性 ===")

	s := "hello"
	fmt.Printf("字符串: %s\n", s)

	// s[0] = 'H' // 编译错误：cannot assign to s[0]

	fmt.Println("字符串是不可变的:")
	fmt.Println("  1. 不能修改字符串的某个字节")
	fmt.Println("  2. 字符串字面量存储在只读段")
	fmt.Println("  3. 修改会导致panic或编译错误")

	fmt.Println("\n修改字符串的方法:")
	fmt.Println("  1. 转换为[]byte修改")
	fmt.Println("     b := []byte(s)")
	fmt.Println("     b[0] = 'H'")
	fmt.Println("     s = string(b)")

	fmt.Println("  2. 使用strings包")
	fmt.Println("     s = strings.Replace(s, \"h\", \"H\", 1)")
}

// DemoStringBytes 演示string与[]byte的转换
func DemoStringBytes() {
	fmt.Println("\n=== String与[]byte转换 ===")

	s := "hello"

	// 标准转换（会复制）
	b := []byte(s)
	fmt.Printf("[]byte: %v\n", b)

	// 转换回string（会复制）
	s2 := string(b)
	fmt.Printf("string: %s\n", s2)

	fmt.Println("\n标准转换特点:")
	fmt.Println("  优点: 安全，不会影响原数据")
	fmt.Println("  缺点: 有内存复制开销")

	fmt.Println("\n零拷贝转换（unsafe）:")
	fmt.Println("  // []byte -> string")
	fmt.Println("  func bytesToString(b []byte) string {")
	fmt.Println("      return *(*string)(unsafe.Pointer(&b))")
	fmt.Println("  }")
	fmt.Println("")
	fmt.Println("  // string -> []byte")
	fmt.Println("  func stringToBytes(s string) []byte {")
	fmt.Println("      return *(*[]byte)(unsafe.Pointer(")
	fmt.Println("          &reflect.SliceHeader{")
	fmt.Println("              Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data,")
	fmt.Println("              Len:  len(s),")
	fmt.Println("              Cap:  len(s),")
	fmt.Println("          }))")
	fmt.Println("  }")

	fmt.Println("\n注意事项:")
	fmt.Println("  1. 零拷贝转换不安全，可能导致数据损坏")
	fmt.Println("  2. 只在性能关键场景使用")
	fmt.Println("  3. 确保转换后的数据不会被修改")
}

// DemoStringConcat 演示字符串拼接
func DemoStringConcat() {
	fmt.Println("\n=== 字符串拼接 ===")

	// 方法1: + 操作符
	s1 := "Hello" + " " + "World"
	fmt.Printf("+ 操作符: %s\n", s1)

	// 方法2: fmt.Sprintf
	s2 := fmt.Sprintf("%s %s", "Hello", "World")
	fmt.Printf("fmt.Sprintf: %s\n", s2)

	// 方法3: strings.Builder
	var builder strings.Builder
	builder.WriteString("Hello")
	builder.WriteString(" ")
	builder.WriteString("World")
	s3 := builder.String()
	fmt.Printf("strings.Builder: %s\n", s3)

	// 方法4: []byte
	buf := make([]byte, 0, 20)
	buf = append(buf, "Hello"...)
	buf = append(buf, ' ')
	buf = append(buf, "World"...)
	s4 := string(buf)
	fmt.Printf("[]byte: %s\n", s4)

	fmt.Println("\n性能对比:")
	fmt.Println("  + 操作符: 每次拼接都会创建新字符串")
	fmt.Println("  fmt.Sprintf: 使用反射，较慢")
	fmt.Println("  strings.Builder: 推荐，高效")
	fmt.Println("  []byte: 高效，但需要预分配容量")

	fmt.Println("\nstrings.Builder优化:")
	fmt.Println("  1. 预分配容量: builder.Grow(n)")
	fmt.Println("  2. 避免多次内存分配")
	fmt.Println("  3. 内部使用[]byte")
}

// DemoStringRune 演示string与rune的关系
func DemoStringRune() {
	fmt.Println("\n=== String与Rune ===")

	s := "Hello世界"

	fmt.Println("字符串遍历方式:")

	// 按字节遍历
	fmt.Println("\n按字节遍历 (for i := 0; i < len(s); i++):")
	for i := 0; i < len(s); i++ {
		fmt.Printf("  s[%d] = %d (0x%02x)\n", i, s[i], s[i])
	}

	// 按rune遍历
	fmt.Println("\n按rune遍历 (for i, r := range s):")
	for i, r := range s {
		fmt.Printf("  s[%d] = %c (0x%04x, %d字节)\n", i, r, r, utf8.RuneLen(r))
	}

	fmt.Println("\nRune与字节的关系:")
	fmt.Println("  ASCII字符: 1字节")
	fmt.Println("  中文等Unicode: 2-4字节")
	fmt.Println("  len(string) 返回字节数")
	fmt.Println("  utf8.RuneCountInString() 返回字符数")

	// rune切片
	runes := []rune(s)
	fmt.Printf("\n[]rune: %v (len=%d)\n", runes, len(runes))
}

// DemoStringIntern 演示字符串驻留
func DemoStringIntern() {
	fmt.Println("\n=== 字符串驻留 ===")

	// 字面量字符串会驻留
	s1 := "hello"
	s2 := "hello"
	fmt.Printf("字面量: s1==s2: %v, 地址相同: %v\n",
		s1 == s2,
		(*StringHeader)(unsafe.Pointer(&s1)).Data == (*StringHeader)(unsafe.Pointer(&s2)).Data)

	// 运行时创建的字符串不会驻留
	s3 := "hel" + "lo"
	s4 := string([]byte{'h', 'e', 'l', 'l', 'o'})
	fmt.Printf("运行时: s3==s4: %v, 地址相同: %v\n",
		s3 == s4,
		(*StringHeader)(unsafe.Pointer(&s3)).Data == (*StringHeader)(unsafe.Pointer(&s4)).Data)

	fmt.Println("\n字符串驻留规则:")
	fmt.Println("  1. 字面量会驻留")
	fmt.Println("  2. 常量表达式会驻留")
	fmt.Println("  3. 运行时创建的字符串不会驻留")

	fmt.Println("\n手动驻留:")
	fmt.Println("  Go没有内置的字符串池")
	fmt.Println("  可以使用map实现:")
	fmt.Println("  var internPool = make(map[string]string)")
	fmt.Println("  func intern(s string) string {")
	fmt.Println("      if v, ok := internPool[s]; ok {")
	fmt.Println("          return v")
	fmt.Println("      }")
	fmt.Println("      internPool[s] = s")
	fmt.Println("      return s")
	fmt.Println("  }")
}

// DemoStringCompare 演示字符串比较
func DemoStringCompare() {
	fmt.Println("\n=== 字符串比较 ===")

	s1 := "hello"
	s2 := "hello"
	s3 := "world"

	fmt.Printf("s1 == s2: %v\n", s1 == s2)
	fmt.Printf("s1 < s3: %v\n", s1 < s3)
	fmt.Printf("s1 > s3: %v\n", s1 > s3)

	fmt.Println("\n比较实现原理:")
	fmt.Println("  1. 先比较长度")
	fmt.Println("  2. 长度相同则比较内存内容")
	fmt.Println("  3. 使用memcmp实现")

	fmt.Println("\n比较优化:")
	fmt.Println("  - 短字符串: 直接比较")
	fmt.Println("  - 长字符串: 先比较前几个字节")
	fmt.Println("  - 相同指针: 直接返回true")

	fmt.Println("\n性能建议:")
	fmt.Println("  1. 使用 == 而非 strings.Compare")
	fmt.Println("  2. 比较前先比较长度")
	fmt.Println("  3. 使用 strings.EqualFold 进行忽略大小写比较")
}

// DemoStringSubstr 演示子字符串
func DemoStringSubstr() {
	fmt.Println("\n=== 子字符串 ===")

	s := "Hello, World!"
	sub := s[0:5]
	fmt.Printf("原字符串: %s (len=%d)\n", s, len(s))
	fmt.Printf("子字符串: %s (len=%d)\n", sub, len(sub))

	// 查看内部指针
	h1 := (*StringHeader)(unsafe.Pointer(&s))
	h2 := (*StringHeader)(unsafe.Pointer(&sub))

	fmt.Printf("\n原字符串指针: 0x%x\n", h1.Data)
	fmt.Printf("子字符串指针: 0x%x (偏移0)\n", h2.Data)

	fmt.Println("\n子字符串特点:")
	fmt.Println("  1. 共享底层数组")
	fmt.Println("  2. 只复制头部（指针+长度）")
	fmt.Println("  3. 可能导致内存泄漏")

	fmt.Println("\n内存泄漏场景:")
	fmt.Println("  s := make([]byte, 1<<20) // 1MB")
	fmt.Println("  sub := string(s[:10])    // 只需要10字节")
	fmt.Println("  // 但整个1MB无法被GC回收!")

	fmt.Println("\n解决方法:")
	fmt.Println("  sub := string([]byte(s[:10])) // 复制一份")
}

// DemoStringSwitch 演示字符串switch优化
func DemoStringSwitch() {
	fmt.Println("\n=== 字符串Switch优化 ===")

	s := "hello"

	switch s {
	case "hello":
		fmt.Println("匹配hello")
	case "world":
		fmt.Println("匹配world")
	default:
		fmt.Println("不匹配")
	}

	fmt.Println("\nSwitch优化策略:")
	fmt.Println("  1. 少量case: 二分查找")
	fmt.Println("  2. 大量case: 使用hash表")
	fmt.Println("  3. 编译器自动选择")

	fmt.Println("\n编译后伪代码:")
	fmt.Println("  // 二分查找")
	fmt.Println("  switch hash(s) {")
	fmt.Println("  case hash(\"hello\"):")
	fmt.Println("      if s == \"hello\" { ... }")
	fmt.Println("  case hash(\"world\"):")
	fmt.Println("      if s == \"world\" { ... }")
	fmt.Println("  }")
}

// DemoStringUTF8 演示UTF-8编码
func DemoStringUTF8() {
	fmt.Println("\n=== UTF-8编码 ===")

	fmt.Println("UTF-8编码规则:")
	fmt.Println("  0xxxxxxx: 1字节 (ASCII, 0-127)")
	fmt.Println("  110xxxxx 10xxxxxx: 2字节 (128-2047)")
	fmt.Println("  1110xxxx 10xxxxxx 10xxxxxx: 3字节 (2048-65535)")
	fmt.Println("  11110xxx 10xxxxxx 10xxxxxx 10xxxxxx: 4字节")

	// 示例
	chars := []string{"A", "中", "😀"}
	for _, c := range chars {
		r := []rune(c)[0]
		fmt.Printf("\n字符: %c\n", r)
		fmt.Printf("  Unicode: U+%04X\n", r)
		fmt.Printf("  UTF-8: %v\n", []byte(c))
		fmt.Printf("  字节数: %d\n", utf8.RuneLen(r))
	}

	fmt.Println("\nUTF-8优点:")
	fmt.Println("  1. ASCII兼容")
	fmt.Println("  2. 变长编码，节省空间")
	fmt.Println("  3. 可以从任意字节开始解码")
	fmt.Println("  4. Go的string原生支持UTF-8")
}

// DemoStringPool 演示字符串内存池
func DemoStringPool() {
	fmt.Println("\n=== 字符串内存池 ===")

	fmt.Println("Go运行时字符串处理:")
	fmt.Println("  1. 字面量存储在二进制文件中")
	fmt.Println("  2. 运行时字符串在堆上分配")
	fmt.Println("  3. 小字符串可能使用tiny allocator")

	fmt.Println("\n内存优化技巧:")
	fmt.Println("  1. 复用字符串")
	fmt.Println("  2. 使用strings.Builder")
	fmt.Println("  3. 避免频繁转换")
	fmt.Println("  4. 使用[]byte替代频繁修改的字符串")

	fmt.Println("\n示例 - 构建大字符串:")
	fmt.Println("  // 不推荐")
	fmt.Println("  var s string")
	fmt.Println("  for i := 0; i < n; i++ {")
	fmt.Println("      s += str  // 每次分配新内存")
	fmt.Println("  }")
	fmt.Println("")
	fmt.Println("  // 推荐")
	fmt.Println("  var builder strings.Builder")
	fmt.Println("  builder.Grow(estimatedSize)")
	fmt.Println("  for i := 0; i < n; i++ {")
	fmt.Println("      builder.WriteString(str)")
	fmt.Println("  }")
	fmt.Println("  s := builder.String()")
}

// DemoStringReflection 演示字符串反射
func DemoStringReflection() {
	fmt.Println("\n=== 字符串反射 ===")

	s := "Hello世界"

	// 使用反射获取字符串信息
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	fmt.Printf("类型: %v\n", t)
	fmt.Printf("种类: %v\n", t.Kind())
	fmt.Printf("长度: %d\n", v.Len())

	// 转换为[]byte
	b := []byte(s)
	fmt.Printf("字节: %v\n", b)

	fmt.Println("\n反射操作:")
	fmt.Println("  reflect.ValueOf(s).Len()  // 长度")
	fmt.Println("  []byte(s) // 字节切片")
	fmt.Println("  reflect.ValueOf(s).String() // 字符串")

	fmt.Println("\n注意:")
	fmt.Println("  - v.Bytes() 返回的是原字符串的字节切片")
	fmt.Println("  - 不要修改返回的切片（字符串不可变）")
}
