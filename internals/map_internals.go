package internals

import (
	"fmt"
	"unsafe"
)

// ============================================================
// Map 底层实现原理
// ============================================================

// hmap 是 Go map 的核心数据结构（简化版）
type hmap struct {
	count     int          // 元素个数
	flags     uint8        // 状态标志
	B         uint8        // bucket数量的对数 (buckets数量 = 2^B)
	noverflow uint16       // 溢出bucket数量
	hash0     uint32       // hash种子
	buckets   unsafe.Pointer // bucket数组指针
	oldbuckets unsafe.Pointer // 扩容时的旧bucket数组
	nevacuate uintptr       // 扩容进度
	extra     unsafe.Pointer // 溢出bucket链表
}

// bmap 是 bucket 的结构（简化版）
type bmap struct {
	tophash [8]uint8 // 存储key的hash高8位
	// 后面紧跟8个key
	// 再后面是8个value
	// 最后是overflow指针
}

// DemoMapStructure 演示map的内存结构
func DemoMapStructure() {
	fmt.Println("=== Map内存结构 ===")

	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3

	// 通过反射获取map信息
	fmt.Printf("map: %v\n", m)
	fmt.Printf("元素个数: %d\n", len(m))

	fmt.Println("\nhmap结构:")
	fmt.Println("+------------------+")
	fmt.Println("| count (int)      | 元素个数")
	fmt.Println("| flags (uint8)    | 状态标志")
	fmt.Println("| B (uint8)        | bucket数量=2^B")
	fmt.Println("| noverflow (uint16)| 溢出bucket数")
	fmt.Println("| hash0 (uint32)   | hash种子")
	fmt.Println("| buckets (*bucket)| bucket数组")
	fmt.Println("| oldbuckets       | 旧bucket(扩容)")
	fmt.Println("| nevacuate        | 扩容进度")
	fmt.Println("+------------------+")

	fmt.Println("\nbmap结构 (每个bucket存8个元素):")
	fmt.Println("+------------------+")
	fmt.Println("| tophash[8]       | hash高8位")
	fmt.Println("| keys[8]          | 8个key")
	fmt.Println("| values[8]        | 8个value")
	fmt.Println("| overflow         | 溢出指针")
	fmt.Println("+------------------+")
}

// DemoMapHash 演示map的hash计算
func DemoMapHash() {
	fmt.Println("\n=== Map Hash计算 ===")

	// hash函数根据key类型选择
	fmt.Println("Hash函数选择:")
	fmt.Println("  string: memhash (基于内存内容)")
	fmt.Println("  int:    memhash (直接hash内存)")
	fmt.Println("  float:  f64hash (特殊处理NaN)")
	fmt.Println("  指针:   memhash (hash地址值)")

	// hash值的用途
	fmt.Println("\nHash值用途:")
	fmt.Println("  低B位: 定位bucket (hash & (2^B-1))")
	fmt.Println("  高8位: 存储在tophash中，快速比较")

	// 示例
	key := "hello"
	hash := hashString(key) // 伪代码，实际由runtime计算
	fmt.Printf("\n示例: key=\"%s\"\n", key)
	fmt.Printf("  hash值: 0x%x\n", hash)
	fmt.Printf("  高8位(tophash): 0x%02x\n", uint8(hash>>56))
	fmt.Printf("  bucket索引(B=3): %d\n", hash&0b111)
}

// DemoMapLookup 演示map的查找过程
func DemoMapLookup() {
	fmt.Println("\n=== Map查找过程 ===")

	fmt.Println("查找步骤:")
	fmt.Println("1. 计算key的hash值")
	fmt.Println("2. 取低B位，定位bucket")
	fmt.Println("3. 取高8位，在bucket中遍历tophash")
	fmt.Println("4. 匹配则对比完整key")
	fmt.Println("5. 不匹配则继续查找overflow bucket")

	fmt.Println("\n查找流程图:")
	fmt.Println("key -> hash -> bucket索引")
	fmt.Println("         |")
	fmt.Println("         v")
	fmt.Println("    +---------+")
	fmt.Println("    | bucket  |")
	fmt.Println("    | tophash | --匹配--> 比较key --> 返回value")
	fmt.Println("    +---------+")
	fmt.Println("         |")
	fmt.Println("         v (不匹配)")
	fmt.Println("    +---------+")
	fmt.Println("    |overflow |")
	fmt.Println("    +---------+")
	fmt.Println("         |")
	fmt.Println("         v")
	fmt.Println("       ...继续")
}

// DemoMapInsert 演示map的插入过程
func DemoMapInsert() {
	fmt.Println("\n=== Map插入过程 ===")

	fmt.Println("插入步骤:")
	fmt.Println("1. 检查是否需要扩容")
	fmt.Println("2. 计算key的hash值")
	fmt.Println("3. 定位bucket")
	fmt.Println("4. 查找空位或更新已有key")
	fmt.Println("5. 如果bucket满，创建overflow bucket")

	fmt.Println("\n插入位置选择优先级:")
	fmt.Println("1. 优先复用已删除的位置 (tophash=emptyOne)")
	fmt.Println("2. 使用bucket中的空位")
	fmt.Println("3. 创建overflow bucket")
}

// DemoMapDelete 演示map的删除过程
func DemoMapDelete() {
	fmt.Println("\n=== Map删除过程 ===")

	m := make(map[int]int)
	for i := 0; i < 10; i++ {
		m[i] = i * 10
	}

	fmt.Printf("删除前: len=%d\n", len(m))
	delete(m, 5)
	fmt.Printf("删除后: len=%d\n", len(m))

	fmt.Println("\n删除步骤:")
	fmt.Println("1. 计算key的hash值")
	fmt.Println("2. 定位bucket")
	fmt.Println("3. 查找key")
	fmt.Println("4. 将tophash标记为emptyOne")
	fmt.Println("5. 清零key和value")

	fmt.Println("\n删除状态标记:")
	fmt.Println("  emptyRest = 0    // 该位置及之后都为空")
	fmt.Println("  emptyOne  = 1    // 该位置为空")
	fmt.Println("  evacuatedX, evacuatedY // 扩容迁移标记")
}

// DemoMapGrow 演示map的扩容机制
func DemoMapGrow() {
	fmt.Println("\n=== Map扩容机制 ===")

	fmt.Println("扩容触发条件:")
	fmt.Println("1. 负载因子 > 6.5 (元素个数 / bucket数)")
	fmt.Println("   -> 等量扩容（bucket数量翻倍）")
	fmt.Println("2. overflow bucket过多")
	fmt.Println("   -> 当B < 15时，overflow > 2^B")
	fmt.Println("   -> 当B >= 15时，overflow > 2^15")
	fmt.Println("   -> 等量扩容（bucket数量不变，整理）")

	fmt.Println("\n扩容过程:")
	fmt.Println("1. 分配新的bucket数组")
	fmt.Println("2. 设置oldbuckets指向旧数组")
	fmt.Println("3. 渐进式迁移（每次操作迁移少量bucket）")
	fmt.Println("4. 迁移完成，释放旧数组")

	fmt.Println("\n渐进式迁移:")
	fmt.Println("- 每次插入/删除操作时，迁移1-2个bucket")
	fmt.Println("- 避免一次性大量数据迁移造成延迟")
	fmt.Println("- 查找时需要同时检查新旧bucket")
}

// DemoMapIterate 演示map的遍历原理
func DemoMapIterate() {
	fmt.Println("\n=== Map遍历原理 ===")

	m := map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}

	fmt.Println("遍历map（每次顺序可能不同）:")
	for i := 0; i < 3; i++ {
		fmt.Printf("第%d次: ", i+1)
		for k := range m {
			fmt.Printf("%d ", k)
		}
		fmt.Println()
	}

	fmt.Println("\n遍历随机性原因:")
	fmt.Println("1. 初始化迭代器时，随机选择起始bucket")
	fmt.Println("2. 遍历顺序由hash决定，不是插入顺序")
	fmt.Println("3. 扩容期间遍历更复杂，需要遍历新旧bucket")

	fmt.Println("\n迭代器实现（伪代码）:")
	fmt.Println("type hiter struct {")
	fmt.Println("    key       unsafe.Pointer")
	fmt.Println("    value     unsafe.Pointer")
	fmt.Println("    t         *maptype")
	fmt.Println("    h         *hmap")
	fmt.Println("    buckets   unsafe.Pointer")
	fmt.Println("    bptr      *bmap")
	fmt.Println("    startBucket uintptr  // 起始bucket（随机）")
	fmt.Println("    offset      uint8    // bucket内偏移（随机）")
	fmt.Println("}")
}

// DemoMapConcurrency 演示map的并发问题
func DemoMapConcurrency() {
	fmt.Println("\n=== Map并发问题 ===")

	fmt.Println("map不是并发安全的!")
	fmt.Println("并发读写会导致:")
	fmt.Println("  1. 数据竞争")
	fmt.Println("  2. 程序崩溃 (fatal error: concurrent map read and write)")
	fmt.Println("  3. 数据损坏")

	fmt.Println("\n解决方案:")
	fmt.Println("1. 使用sync.Mutex")
	fmt.Println("   var mu sync.Mutex")
	fmt.Println("   mu.Lock(); m[k] = v; mu.Unlock()")

	fmt.Println("2. 使用sync.RWMutex (读多写少)")
	fmt.Println("   var rwmu sync.RWMutex")
	fmt.Println("   rwmu.RLock(); v := m[k]; rwmu.RUnlock()")

	fmt.Println("3. 使用sync.Map (读多写少，key相对稳定)")
	fmt.Println("   var sm sync.Map")
	fmt.Println("   sm.Store(key, value)")
	fmt.Println("   v, ok := sm.Load(key)")
}

// DemoMapKeyTypes 演示map的key类型限制
func DemoMapKeyTypes() {
	fmt.Println("\n=== Map Key类型限制 ===")

	fmt.Println("Key必须可比较 (comparable):")

	fmt.Println("\n可用类型:")
	fmt.Println("  ✓ 基本类型: int, float, string, bool")
	fmt.Println("  ✓ 指针: *T")
	fmt.Println("  ✓ 数组: [N]T (元素可比较)")
	fmt.Println("  ✓ 结构体: struct{...} (字段可比较)")

	fmt.Println("\n不可用类型:")
	fmt.Println("  ✗ 切片: []T")
	fmt.Println("  ✗ map: map[K]V")
	fmt.Println("  ✗ 函数: func(...)")
	fmt.Println("  ✗ 包含不可比较字段的结构体")

	// 示例
	m1 := make(map[[2]int]string) // 数组可以作为key
	m1[[2]int{1, 2}] = "array key"
	fmt.Printf("\n数组作为key: %v\n", m1)

	type Point struct{ X, Y int }
	m2 := make(map[Point]string) // 结构体可以作为key
	m2[Point{1, 2}] = "struct key"
	fmt.Printf("结构体作为key: %v\n", m2)
}

// DemoMapMemory 演示map的内存占用
func DemoMapMemory() {
	fmt.Println("\n=== Map内存占用 ===")

	fmt.Println("内存占用计算:")
	fmt.Println("  hmap结构: ~48字节")
	fmt.Println("  每个bucket: 8字节tophash + 8个key + 8个value + 8字节overflow指针")

	fmt.Println("\n示例 (map[int]int, B=0, 1个bucket):")
	fmt.Println("  hmap: 48字节")
	fmt.Println("  bucket: 8 + 8*8 + 8*8 + 8 = 168字节")
	fmt.Println("  总计: ~216字节")

	fmt.Println("\n内存优化建议:")
	fmt.Println("  1. 预分配容量: make(map[K]V, hint)")
	fmt.Println("  2. 使用指针value减少复制开销")
	fmt.Println("  3. 小map考虑用slice替代")
}

// hashString 模拟字符串hash（实际由runtime实现）
func hashString(s string) uint64 {
	h := uint64(0)
	for _, c := range s {
		h = h*31 + uint64(c)
	}
	return h
}
