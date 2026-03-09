# Init 函数模块

## 学习目标
理解Go语言init函数的执行顺序和初始化机制。

## 核心概念

### 1. init函数特点
- 每个包可以有多个init函数
- 每个源文件可以有多个init函数
- init函数在main函数之前自动执行
- init函数不能被调用，不能有参数和返回值

### 2. 执行顺序
1. 导入的包先初始化
2. 包级别变量初始化
3. init函数执行
4. main函数执行

### 3. init函数规则
- 同一文件: 从上到下
- 同一包不同文件: 按文件名字母顺序
- 不同包: 按导入依赖关系

## 示例说明

本模块通过多个包演示init函数的执行顺序：

| 包 | 说明 |
|------|------|
| init_pkg0 | 不依赖其他包 |
| init_pkg1 | 不依赖其他包 |
| init_pkg2 | 依赖init_pkg1 |
| init_pkg3 | 依赖init_pkg0, init_pkg2, init_pkg4 |
| init_pkg4 | 不依赖其他包 |

## 执行顺序图解

```
导入关系:
init_pkg3 → init_pkg0
           → init_pkg2 → init_pkg1
           → init_pkg4

执行顺序:
1. init_pkg1 (被init_pkg2依赖)
2. init_pkg2 (依赖init_pkg1)
3. init_pkg0 (被init_pkg3直接依赖)
4. init_pkg4 (被init_pkg3直接依赖)
5. init_pkg3 (主包)
```

## init函数规则详解

### 1. 同一文件内的init
```go
// 按定义顺序执行
func init() {
    fmt.Println("init 1")
}

func init() {
    fmt.Println("init 2")
}
```

### 2. 同一包不同文件
```
mypackage/
├── a.go  // init先执行 (a < b)
└── b.go  // init后执行
```

### 3. 包依赖顺序
```go
// main.go
import (
    "packageA"  // 先初始化
    "packageB"  // 后初始化
)

// packageB依赖packageA时:
// 1. packageA先初始化
// 2. packageB后初始化
// 3. main最后初始化
```

## 常见用途

### 1. 初始化配置
```go
var config Config

func init() {
    config = loadConfig()
}
```

### 2. 注册驱动
```go
func init() {
    sql.Register("mysql", &MySQLDriver{})
}
```

### 3. 初始化连接池
```go
var pool *Pool

func init() {
    pool = NewPool(10)
}
```

## 常见陷阱

### 1. 循环依赖
```go
// 包A导入包B，包B导入包A → 编译错误
```

### 2. init中执行耗时操作
```go
func init() {
    time.Sleep(10 * time.Second)  // 阻塞程序启动
}
```

### 3. 依赖全局变量顺序
```go
var a = b + 1  // b还未初始化
var b = 1
```

## 最佳实践

1. **保持简单**: init函数应简单快速
2. **避免依赖**: 不要依赖其他包的init执行顺序
3. **错误处理**: init中panic会导致程序退出
4. **显式初始化**: 复杂初始化考虑使用显式初始化函数
5. **文档说明**: 说明init的作用和执行时机
