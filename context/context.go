package context

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func BackgroundDemo() {
	fmt.Println("=== Background和TODO示例 ===")

	ctx := context.Background()
	fmt.Printf("Background context: %v\n", ctx)

	todoCtx := context.TODO()
	fmt.Printf("TODO context: %v\n", todoCtx)

	fmt.Println("Background用于main函数、初始化和测试")
	fmt.Println("TODO用于不确定使用哪个context时")
}

func WithCancelDemo() {
	fmt.Println("=== WithCancel取消示例 ===")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Go(func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Worker 1: 收到取消信号，退出")
				return
			default:
				fmt.Printf("Worker 1: 工作 %d\n", i)
				time.Sleep(100 * time.Millisecond)
			}
		}
	})

	wg.Go(func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Worker 2: 收到取消信号，退出")
				return
			default:
				fmt.Printf("Worker 2: 工作 %d\n", i)
				time.Sleep(150 * time.Millisecond)
			}
		}
	})

	time.Sleep(500 * time.Millisecond)
	fmt.Println("主函数: 发送取消信号")
	cancel()

	wg.Wait()
	fmt.Println("主函数: 所有worker已退出")
}

func WithTimeoutDemo() {
	fmt.Println("=== WithTimeout超时示例 ===")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("操作完成")
	case <-ctx.Done():
		fmt.Println("操作超时:", ctx.Err())
	}
}

func WithDeadlineDemo() {
	fmt.Println("=== WithDeadline截止时间示例 ===")

	deadline := time.Now().Add(200 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("截止时间: %v\n", deadline.Format("15:04:05.000"))

	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("操作完成")
	case <-ctx.Done():
		fmt.Println("已超过截止时间:", ctx.Err())
	}
}

func WithValueDemo() {
	fmt.Println("=== WithValue值传递示例 ===")

	type key string

	ctx := context.Background()
	ctx = context.WithValue(ctx, key("userID"), 12345)
	ctx = context.WithValue(ctx, key("requestID"), "abc-123")

	userID, ok := ctx.Value(key("userID")).(int)
	if ok {
		fmt.Printf("UserID: %d\n", userID)
	}

	requestID, ok := ctx.Value(key("requestID")).(string)
	if ok {
		fmt.Printf("RequestID: %s\n", requestID)
	}

	missing := ctx.Value(key("missing"))
	if missing == nil {
		fmt.Println("missing key: 值不存在")
	}
}

func ChainContextDemo() {
	fmt.Println("=== 链式Context示例 ===")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "level", 1)

	ctx2, cancel2 := context.WithCancel(ctx)
	defer cancel2()
	ctx2 = context.WithValue(ctx2, "level", 2)

	ctx3, cancel3 := context.WithTimeout(ctx2, 500*time.Millisecond)
	defer cancel3()
	ctx3 = context.WithValue(ctx3, "level", 3)

	fmt.Printf("ctx level: %v\n", ctx.Value("level"))
	fmt.Printf("ctx2 level: %v\n", ctx2.Value("level"))
	fmt.Printf("ctx3 level: %v\n", ctx3.Value("level"))

	select {
	case <-ctx3.Done():
		fmt.Println("ctx3 完成:", ctx3.Err())
	}
}

func processRequestRecursive(ctx context.Context, depth int) {
	fmt.Printf("深度 %d: 开始处理\n", depth)

	select {
	case <-ctx.Done():
		fmt.Printf("深度 %d: 被取消\n", depth)
		return
	default:
	}

	if depth < 3 {
		processRequestRecursive(ctx, depth+1)
	}

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("深度 %d: 处理完成\n", depth)
}

func PropagationDemo() {
	fmt.Println("=== Context传播示例 ===")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	processRequestRecursive(ctx, 1)
}

func ErrDemo() {
	fmt.Println("=== Context错误处理示例 ===")

	ctx := context.Background()
	fmt.Printf("初始错误: %v\n", ctx.Err())

	ctx1, cancel1 := context.WithCancel(ctx)
	fmt.Printf("WithCancel后错误: %v\n", ctx1.Err())

	cancel1()
	fmt.Printf("取消后错误: %v\n", ctx1.Err())

	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	time.Sleep(150 * time.Millisecond)
	fmt.Printf("超时后错误: %v\n", ctx2.Err())

	fmt.Println("\n错误类型:")
	fmt.Println("- context.Canceled: 主动取消")
	fmt.Println("- context.DeadlineExceeded: 超时")
}
