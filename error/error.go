package learn_error

import (
	"errors"
	"fmt"
	"reflect"
)

type MyError struct {
	error
}

func returnsError() error {
	var e2 *MyError = nil
	return e2
}

func MyErrorEqualsNil() {
	// 内置error的nil
	var e error = nil
	fmt.Printf("isEqual: %v type: %s error: %v \n", e == nil, reflect.TypeOf(e), e)
	// 自定义error的nil
	var e1 *MyError = nil
	fmt.Printf("isEqual: %v type: %s error: %v \n", e1 == nil, reflect.TypeOf(e1), e1)
	// 自定义error的类型存在，值为nil
	e2 := returnsError()
	fmt.Printf("isEqual: %v type: %s error2: %v \n", e2 == nil, reflect.TypeOf(e2), e2)
	// 自定义error的类型存在，值也存在，内部变量为nil
	e3 := &MyError{error: nil}
	fmt.Printf("isEqual: %v type: %s error: %v \n", e3 == nil, reflect.TypeOf(e3), e3)
}

// DivisionError 自定义错误类型，包含详细信息
type DivisionError struct {
	Dividend int
	Divisor  int
	Message  string
}

func (e *DivisionError) Error() string {
	return fmt.Sprintf("division error: %s (dividend=%d, divisor=%d)", e.Message, e.Dividend, e.Divisor)
}

// CustomError 演示自定义错误类型
func CustomError() {
	fmt.Println("=== 自定义错误类型演示 ===")

	divide := func(a, b int) (int, error) {
		if b == 0 {
			return 0, &DivisionError{
				Dividend: a,
				Divisor:  b,
				Message:  "division by zero",
			}
		}
		return a / b, nil
	}

	// 正常情况
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %d\n", result)
	}

	// 错误情况
	result, err = divide(10, 0)
	if err != nil {
		// 类型断言获取详细信息
		var divErr *DivisionError
		if errors.As(err, &divErr) {
			fmt.Printf("除法错误: %s\n", divErr.Message)
			fmt.Printf("被除数: %d, 除数: %d\n", divErr.Dividend, divErr.Divisor)
		}
	}
}

// ErrorWrapping 演示错误包装 (Go 1.13+)
func ErrorWrapping() {
	fmt.Println("=== 错误包装演示 ===")

	// 模拟底层错误
	openFile := func(name string) error {
		return fmt.Errorf("open %s: no such file", name)
	}

	// 中间层包装错误
	readConfig := func() error {
		err := openFile("config.yaml")
		if err != nil {
			return fmt.Errorf("read config: %w", err) // 使用 %w 包装错误
		}
		return nil
	}

	// 上层再包装
	loadApp := func() error {
		err := readConfig()
		if err != nil {
			return fmt.Errorf("load application: %w", err)
		}
		return nil
	}

	err := loadApp()
	if err != nil {
		fmt.Printf("错误链: %v\n", err)
	}
}

// ErrorUnwrap: 演示错误解包
func ErrorUnwrap() {
	fmt.Println("=== 错误解包演示 ===")

	// 创建错误链
	baseErr := errors.New("base error")
	wrapped1 := fmt.Errorf("layer1: %w", baseErr)
	wrapped2 := fmt.Errorf("layer2: %w", wrapped1)

	fmt.Printf("最外层错误: %v\n", wrapped2)

	// 逐层解包
	unwrapped1 := errors.Unwrap(wrapped2)
	fmt.Printf("第一层解包: %v\n", unwrapped1)

	unwrapped2 := errors.Unwrap(unwrapped1)
	fmt.Printf("第二层解包: %v\n", unwrapped2)

	unwrapped3 := errors.Unwrap(unwrapped2)
	fmt.Printf("第三层解包: %v\n", unwrapped3) // nil
}

// SentinelError 哨兵错误定义
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrTimeout      = errors.New("operation timeout")
)

// ErrorIs 演示使用 errors.Is 判断错误
func ErrorIs() {
	fmt.Println("=== errors.Is 演示 ===")

	// 模拟返回错误的函数
	findUser := func(id int) error {
		if id <= 0 {
			return ErrNotFound
		}
		return nil
	}

	// 模拟包装后的错误
	findUserWrapped := func(id int) error {
		err := findUser(id)
		if err != nil {
			return fmt.Errorf("find user failed: %w", err)
		}
		return nil
	}

	// 测试
	err := findUserWrapped(-1)
	if err != nil {
		fmt.Printf("错误: %v\n", err)

		// 使用 errors.Is 判断错误类型（支持错误链）
		if errors.Is(err, ErrNotFound) {
			fmt.Println("检测到: 资源未找到错误")
		}
		if errors.Is(err, ErrUnauthorized) {
			fmt.Println("检测到: 未授权错误")
		}
	}
}

// NetworkError 网络错误类型
type NetworkError struct {
	Code    int
	Message string
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: code=%d, message=%s", e.Code, e.Message)
}

// TimeoutError 超时错误类型
type TimeoutError struct {
	Duration int
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout after %dms", e.Duration)
}

// ErrorAs 演示使用 errors.As 转换错误类型
func ErrorAs() {
	fmt.Println("=== errors.As 演示 ===")

	// 模拟返回不同类型错误的函数
	doRequest := func(code int) error {
		if code == 500 {
			return &NetworkError{Code: code, Message: "internal server error"}
		}
		if code == 408 {
			return &TimeoutError{Duration: 30000}
		}
		return nil
	}

	// 测试网络错误
	err := doRequest(500)
	if err != nil {
		var netErr *NetworkError
		if errors.As(err, &netErr) {
			fmt.Printf("网络错误 - Code: %d, Message: %s\n", netErr.Code, netErr.Message)
		}

		var timeoutErr *TimeoutError
		if errors.As(err, &timeoutErr) {
			fmt.Printf("超时错误 - Duration: %dms\n", timeoutErr.Duration)
		}
	}

	// 测试超时错误
	err = doRequest(408)
	if err != nil {
		var netErr *NetworkError
		if errors.As(err, &netErr) {
			fmt.Printf("网络错误 - Code: %d, Message: %s\n", netErr.Code, netErr.Message)
		}

		var timeoutErr *TimeoutError
		if errors.As(err, &timeoutErr) {
			fmt.Printf("超时错误 - Duration: %dms\n", timeoutErr.Duration)
		}
	}
}

// PanicRecover 演示 panic 和 recover
func PanicRecover() {
	fmt.Println("=== Panic 和 Recover 演示 ===")

	// 安全执行函数
	safeExecute := func(fn func()) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v", r)
			}
		}()
		fn()
		return nil
	}

	// 正常函数
	normalFn := func() {
		fmt.Println("正常执行...")
	}

	// 会panic的函数
	panicFn := func() {
		fmt.Println("即将panic...")
		panic("something went wrong")
	}

	// 执行正常函数
	err := safeExecute(normalFn)
	fmt.Printf("正常函数执行结果: err=%v\n", err)

	// 执行会panic的函数
	err = safeExecute(panicFn)
	fmt.Printf("Panic函数执行结果: err=%v\n", err)
}

// MultiError 多错误处理
type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}
	return fmt.Sprintf("%d errors occurred", len(e.Errors))
}

func (e *MultiError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// MultiErrorDemo 演示多错误处理
func MultiErrorDemo() {
	fmt.Println("=== 多错误处理演示 ===")

	// 模拟批量操作
	items := []string{"item1", "item2", "", "item4", ""}
	processItems := func(items []string) error {
		var multiErr MultiError
		for i, item := range items {
			if item == "" {
				multiErr.Add(fmt.Errorf("item %d is empty", i))
			} else {
				fmt.Printf("处理: %s\n", item)
			}
		}
		if multiErr.HasErrors() {
			return &multiErr
		}
		return nil
	}

	err := processItems(items)
	if err != nil {
		if multiErr, ok := err.(*MultiError); ok {
			fmt.Printf("发现 %d 个错误:\n", len(multiErr.Errors))
			for _, e := range multiErr.Errors {
				fmt.Printf("  - %v\n", e)
			}
		}
	}
}

// ErrorWithContext 带上下文的错误
type ErrorWithContext struct {
	Op   string
	Path string
	Err  error
	Time string
}

func (e *ErrorWithContext) Error() string {
	return fmt.Sprintf("%s %s: %v [%s]", e.Op, e.Path, e.Err, e.Time)
}

func (e *ErrorWithContext) Unwrap() error {
	return e.Err
}

// ErrorWithContextDemo 演示带上下文的错误
func ErrorWithContextDemo() {
	fmt.Println("=== 带上下文的错误演示 ===")

	// 模拟文件操作
	deleteFile := func(path string) error {
		return &ErrorWithContext{
			Op:   "delete",
			Path: path,
			Err:  errors.New("permission denied"),
			Time: "2024-01-01 10:00:00",
		}
	}

	err := deleteFile("/tmp/important.txt")
	if err != nil {
		fmt.Printf("错误: %v\n", err)

		// 类型断言获取详细信息
		var ctxErr *ErrorWithContext
		if errors.As(err, &ctxErr) {
			fmt.Printf("操作: %s\n", ctxErr.Op)
			fmt.Printf("路径: %s\n", ctxErr.Path)
			fmt.Printf("时间: %s\n", ctxErr.Time)
			fmt.Printf("原因: %v\n", ctxErr.Err)
		}
	}
}

// DeferPanicOrder 演示 defer、panic、recover 的执行顺序
func DeferPanicOrder() {
	fmt.Println("=== Defer、Panic、Recover 执行顺序 ===")

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recover 捕获: %v\n", err)
		}
	}()

	defer fmt.Println("Defer 1")
	defer fmt.Println("Defer 2")

	fmt.Println("正常执行")
	panic("主动触发 panic")

	defer fmt.Println("Defer 3 (不会执行)")
}
