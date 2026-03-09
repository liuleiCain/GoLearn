package concurrency_patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func WorkerPoolPattern() {
	fmt.Println("=== Worker Pool模式 ===")

	jobs := make(chan int, 10)
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
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for j := 1; j <= 10; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	for r := range results {
		fmt.Printf("结果: %d\n", r)
	}
}

func PipelinePattern() {
	fmt.Println("=== Pipeline模式 ===")

	generator := func(done <-chan struct{}, nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				select {
				case out <- n:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	square := func(done <-chan struct{}, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				select {
				case out <- n * n:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	double := func(done <-chan struct{}, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				select {
				case out <- n * 2:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	done := make(chan struct{})
	defer close(done)

	nums := generator(done, 1, 2, 3, 4, 5)
	squared := square(done, nums)
	doubled := double(done, squared)

	for result := range doubled {
		fmt.Printf("Pipeline结果: %d\n", result)
	}
}

func FanOutFanIn() {
	fmt.Println("=== Fan-out/Fan-in模式 ===")

	producer := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	worker := func(name string, in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for n := range in {
				time.Sleep(50 * time.Millisecond)
				out <- fmt.Sprintf("%s: %d -> %d", name, n, n*n)
			}
		}()
		return out
	}

	merger := func(channels ...<-chan string) <-chan string {
		var wg sync.WaitGroup
		out := make(chan string)

		output := func(c <-chan string) {
			defer wg.Done()
			for s := range c {
				out <- s
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go output(c)
		}

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
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id)
		}(i)
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
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("任务 %d: 执行完成\n", id)
	}

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			task(id)
		}(i)
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
