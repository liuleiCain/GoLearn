package concurrency_patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func WorkerPoolPattern() {
	fmt.Println("=== Worker Pool模式 ===")

	jobs := make(chan int, 5)
	results := make(chan int, 10)

	worker := func(id int, jobs <-chan int, results chan<- int) {
		for j := range jobs {
			fmt.Printf("Worker %d: 处理任务 %d\n", id, j)
			time.Sleep(100 * time.Millisecond)
			results <- j * 2
		}
	}

	numWorkers := 3
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Go(func() {
			worker(w, jobs, results)
		})
	}

	go func() {
		for j := 1; j <= 10; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("结果: %d\n", r)
	}
}

func PipelinePattern() {
	fmt.Println("=== Pipeline模式 ===")

	// generator 生成器函数，将一组整数发送到通道
	// done 取消信号通道，用于优雅关闭
	// nums 可变参数，要发送的整数列表
	// 返回 只读整数通道，用于接收生成的数字
	generator := func(done <-chan struct{}, nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out) // 确保goroutine退出时关闭通道
			for _, n := range nums {
				select {
				case out <- n: // 将数字发送到输出通道
				case <-done: // 收到取消信号，立即退出
					return
				}
			}
		}()
		return out
	}

	// square 平方函数，计算输入数字的平方
	// done 取消信号通道
	// in 输入通道，接收要处理的数字
	// 返回 只读整数通道，输出平方后的结果
	square := func(done <-chan struct{}, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in { // 从输入通道读取数据
				select {
				case out <- n * n: // 发送平方结果
				case <-done: // 收到取消信号，立即退出
					return
				}
			}
		}()
		return out
	}

	// double 双倍函数，将输入数字乘以2
	// done 取消信号通道
	// in 输入通道，接收要处理的数字
	// 返回 只读整数通道，输出双倍后的结果
	double := func(done <-chan struct{}, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				select {
				case out <- n * 2: // 发送双倍结果
				case <-done: // 收到取消信号，立即退出
					return
				}
			}
		}()
		return out
	}

	// 创建done通道，用于发送取消信号
	done := make(chan struct{})
	defer close(done) // 函数退出时关闭done，通知所有goroutine停止

	// 构建Pipeline流水线 generator -> square -> double
	// 数据流向 1,2,3,4,5 -> 1,4,9,16,25 -> 2,8,18,32,50
	nums := generator(done, 1, 2, 3, 4, 5)
	squared := square(done, nums)
	doubled := double(done, squared)

	// 从最终输出通道读取结果
	for result := range doubled {
		fmt.Printf("Pipeline结果: %d\n", result)
	}
}

func FanOutFanIn() {
	fmt.Println("=== Fan-out/Fan-in模式 ===")

	// producer 生产者函数，将一组整数发送到通道
	// nums 可变参数，要发送的整数列表
	// 返回 只读整数通道，用于接收生成的数字
	producer := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out) // 确保goroutine退出时关闭通道
			for _, n := range nums {
				out <- n // 将数字发送到输出通道
			}
		}()
		return out
	}

	// worker 工作者函数，处理输入的数字并返回格式化结果
	// name 工作者名称，用于标识输出
	// in 输入通道，接收要处理的数字
	// 返回 只读字符串通道，输出处理结果
	worker := func(name string, in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for n := range in { // 从输入通道读取数据
				time.Sleep(50 * time.Millisecond)                // 模拟处理耗时
				out <- fmt.Sprintf("%s: %d -> %d", name, n, n*n) // 发送格式化结果
			}
		}()
		return out
	}

	// merger 合并器函数，将多个通道的数据合并到一个通道
	// channels 可变参数，多个只读字符串通道
	// 返回 只读字符串通道，输出合并后的所有数据
	merger := func(channels ...<-chan string) <-chan string {
		var wg sync.WaitGroup
		out := make(chan string)

		// output: 从单个通道读取数据并转发到输出通道
		output := func(c <-chan string) {
			defer wg.Done()
			for s := range c {
				out <- s // 将数据转发到合并输出通道
			}
		}

		wg.Add(len(channels)) // 设置等待计数为通道数量
		for _, c := range channels {
			go output(c) // 为每个输入通道启动一个goroutine
		}

		// 等待所有goroutine完成后关闭输出通道
		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}

	input := producer(1, 2, 3, 4, 5, 6, 7, 8, 9)

	c1 := worker("Worker1", input)
	c2 := worker("Worker2", input)
	c3 := worker("Worker3", input)

	for result := range merger(c1, c2, c3) {
		fmt.Printf("Fan-out/Fan-in结果: %s\n", result)
	}
}

func TimeoutPattern() {
	fmt.Println("=== 超时模式 ===")

	slowOperation := func() <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			time.Sleep(200 * time.Millisecond)
			out <- "操作完成"
		}()
		return out
	}

	fmt.Println("短超时(100ms):")
	select {
	case result := <-slowOperation():
		fmt.Println(result)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("超时!")
	}

	fmt.Println("\n长超时(500ms):")
	select {
	case result := <-slowOperation():
		fmt.Println(result)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("超时!")
	}
}

func CancellationPattern() {
	fmt.Println("=== 取消模式 ===")

	worker := func(ctx context.Context, id int) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d: 收到取消信号，退出\n", id)
				return
			default:
				fmt.Printf("Worker %d: 工作中...\n", id)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Go(func() {
			worker(ctx, i)
		})
	}

	time.Sleep(300 * time.Millisecond)
	fmt.Println("主函数: 发送取消信号")
	cancel()

	wg.Wait()
	fmt.Println("主函数: 所有worker已退出")
}

func SemaphorePattern() {
	fmt.Println("=== 信号量模式 ===")

	sem := make(chan struct{}, 3)

	task := func(id int) {
		sem <- struct{}{}
		defer func() { <-sem }()

		fmt.Printf("任务 %d: 开始执行\n", id)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("任务 %d: 执行完成\n", id)
	}

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Go(func() {
			task(i)
		})
	}

	wg.Wait()
	fmt.Println("所有任务完成")
}

func RateLimitingPattern() {
	fmt.Println("=== 速率限制模式 ===")

	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(200 * time.Millisecond)

	for req := range requests {
		<-limiter
		fmt.Printf("请求 %d: %v\n", req, time.Now().Format("15:04:05.000"))
	}
}

func GracefulShutdown() {
	fmt.Println("=== 优雅关闭模式 ===")

	server := func(ctx context.Context, id int) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("服务器 %d: 正在关闭...\n", id)
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("服务器 %d: 已关闭\n", id)
				return
			default:
				fmt.Printf("服务器 %d: 处理请求\n", id)
				time.Sleep(50 * time.Millisecond)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			server(ctx, id)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有服务器已优雅关闭")
}
