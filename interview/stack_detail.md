# Go 栈管理与抢占详解

本文档讲解 goroutine 栈的增长、拷贝、收缩以及抢占机制。

---

## 一、为什么需要可增长栈

- Goroutine 数量通常远大于系统线程，固定大栈会造成巨大内存浪费。
- Go 使用“按需增长”的栈设计：初始小栈，遇到需要时再扩容。

---

## 二、栈增长机制（概念流程）

1. 函数序言会检查 `stackguard0`，判断当前栈是否有足够空间。
2. 空间不足时调用 `morestack` 申请更大的栈。
3. 在安全点复制旧栈到新栈，继续执行。

要点：
- 栈增长是**按需**发生的。
- 复制发生在**安全点**，确保指针修复安全。

---

## 三、栈拷贝与指针修复

- 栈扩容时会把旧栈内容复制到新栈。
- 运行时会根据栈帧元数据修复指针（stack map），避免悬空引用。
- 该过程对用户代码透明，但会影响“栈上地址稳定性”。

---

## 四、栈收缩

- 当 goroutine 使用的栈远小于当前大小时，运行时可能在安全点收缩栈。
- 收缩常发生在 GC 过程中，目的是回收空闲栈空间。

---

## 五、g0 / system stack

- 每个线程有一个 `g0`，用于执行调度器和运行时逻辑。
- `g0` 使用系统栈（system stack），与普通 goroutine 栈分离。
- `//go:systemstack` 和 `//go:nosplit` 是运行时常用的指令（用户代码通常不需要）。

---

## 六、抢占机制（概念）

- **协作式抢占**：在安全点检查抢占标志。
- **异步抢占**：通过信号打断长时间运行的 goroutine，使其进入安全点。
- 目的：避免某个 goroutine 长时间占用 CPU，提升调度公平性。

---

## 七、代码示例

```go
// 观测栈地址变化（仅用于理解栈拷贝的可能性）
func StackGrowthDemo() {
    const depth = 32
    addrs := make([]uintptr, 0, depth+1)

    var walk func(int)
    walk = func(n int) {
        var local [256]byte
        local[0] = byte(n) // 防止优化
        addrs = append(addrs, uintptr(unsafe.Pointer(&local)))
        if n > 0 {
            walk(n - 1)
        }
    }
    walk(depth)

    fmt.Printf("采样深度: %d\n", len(addrs))
    fmt.Printf("栈地址范围: 0x%x -> 0x%x\n", addrs[0], addrs[len(addrs)-1])

    monotonic := true
    for i := 1; i < len(addrs); i++ {
        if addrs[i] >= addrs[i-1] {
            monotonic = false
            break
        }
    }
    fmt.Printf("地址是否单调递减: %v (非单调可能发生栈拷贝)\n", monotonic)
}
```
