# Go 内存分配器详解

本文档聚焦 Go runtime 的内存分配器设计、分配路径，以及与 GC/逃逸分析的关系。

---

## 一、整体架构

Go 运行时分配器采用“本地缓存 + 中央缓存 + 全局堆”的分层设计，核心目标是减少锁竞争并提高分配速度。

```
P 的本地缓存 (mcache)
        |
        v
中央缓存 (mcentral)
        |
        v
全局堆 (mheap)
        |
        v
OS (mmap / sysAlloc)
```

### 1.1 关键结构（概念级）

- `mcache`：每个 P 私有，分配小对象的最快路径，基本无锁。
- `mcentral`：按 size class 维护 span 列表，向 mcache 供给可用 span。
- `mheap`：全局堆，负责向 OS 申请/回收内存，并管理 span。
- `mspan`：一段连续内存页，内部按固定大小切分对象。

---

## 二、分配路径

### 2.1 小对象分配（典型路径）

1. 根据大小选择 size class。
2. 从当前 P 的 `mcache` 中获取可用 span。
3. 若 span 耗尽，则向 `mcentral` 申请新的 span。
4. `mcentral` 不足时再向 `mheap` 申请，必要时向 OS 申请新内存。

### 2.2 大对象分配

- 大对象通常绕过 size class，直接向 `mheap` 申请专用 span。
- 代价更高，但避免在小对象 span 中造成碎片。

---

## 三、size class 与 span

### 3.1 size class

- 将小对象划分为离散的尺寸等级，减少内部碎片和元数据开销。
- 同一 size class 的对象共用一个 span，分配/回收成本低。

### 3.2 span 与 GC 位图

- span 除了对象内存，还包含 GC 标记/分配位图（bitmap）。
- 对于不含指针的对象，会进入 `noscan` 路径，减少 GC 扫描成本。

---

## 四、tiny allocator（概念）

- 对非常小且不含指针的对象，运行时会尝试“打包分配”以减少元数据开销。
- 该机制属于优化路径，具体细节可能随版本变化。

---

## 五、逃逸分析与分配位置

编译器会通过逃逸分析决定对象在栈还是堆上分配：

- 返回局部变量指针
- 闭包捕获外部变量
- 接口装箱（interface conversion）
- 大对象或生命周期不确定

查看逃逸分析：

```bash
go build -gcflags="-m" ./...
```

---

## 六、观测与调优

### 6.1 runtime.MemStats

- `Mallocs / Frees`：分配与释放次数
- `HeapObjects`：当前堆对象数量
- `HeapAlloc / HeapInuse`：堆使用情况

### 6.2 常见工具

- `pprof`：定位分配热点
- `GODEBUG=allocfreetrace=1`：跟踪分配与释放（调试用途）

---

## 七、代码示例

```go
// 观察分配统计信息
func AllocatorStatsDemo() {
    var before, after runtime.MemStats
    runtime.ReadMemStats(&before)

    // 小对象分配
    small := make([][]byte, 10000)
    for i := 0; i < len(small); i++ {
        small[i] = make([]byte, 32+(i%128))
    }

    // 大对象分配
    big := make([][]byte, 64)
    for i := 0; i < len(big); i++ {
        big[i] = make([]byte, 256*1024)
    }

    runtime.ReadMemStats(&after)

    fmt.Printf("Mallocs: %d -> %d\n", before.Mallocs, after.Mallocs)
    fmt.Printf("HeapObjects: %d -> %d\n", before.HeapObjects, after.HeapObjects)

    runtime.KeepAlive(small)
    runtime.KeepAlive(big)
}
```
