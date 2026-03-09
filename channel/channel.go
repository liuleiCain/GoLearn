package channel

import (
	"fmt"
	"sync"
	"time"
)

func UnbufferedChannel() {
	fmt.Println("=== 无缓冲通道示例 ===")

	ch := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "Hello from goroutine"
	}()

	fmt.Println("等待接收消息...")
	msg := <-ch
	fmt.Println("接收到:", msg)
}

func BufferedChannel() {
	fmt.Println("=== 有缓冲通道示例 ===")

	ch := make(chan int, 3)

	fmt.Println("发送3个值到缓冲通道...")
	for i := 0; i < 3; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
	}
	//ch <- 1
	//ch <- 2
	//ch <- 3

	fmt.Printf("通道长度: %d, 容量: %d\n", len(ch), cap(ch))

	fmt.Println("接收值:", <-ch)
	fmt.Println("接收值:", <-ch)
	fmt.Println("接收值:", <-ch)
}

func ChannelDirection() {
	fmt.Println("=== 通道方向示例 ===")

	sendOnly := func(ch chan<- int) {
		for i := 0; i < 3; i++ {
			ch <- i
			fmt.Printf("发送: %d\n", i)
		}
		close(ch)
	}

	receiveOnly := func(ch <-chan int) {
		for v := range ch {
			fmt.Printf("接收: %d\n", v)
		}
	}

	ch := make(chan int, 3)
	go sendOnly(ch)
	receiveOnly(ch)
}

func CloseChannel() {
	fmt.Println("=== 关闭通道示例 ===")

	ch := make(chan int, 5)

	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch)
		fmt.Println("通道已关闭")
	}()

	for {
		val, ok := <-ch
		if !ok {
			fmt.Println("通道已关闭，退出循环")
			break
		}
		fmt.Println("接收到:", val)
	}

	fmt.Println("\n使用range遍历通道:")
	ch2 := make(chan int, 3)
	go func() {
		ch2 <- 1
		ch2 <- 2
		ch2 <- 3
		close(ch2)
	}()

	for v := range ch2 {
		fmt.Println("range接收到:", v)
	}
	fmt.Println("通道已关闭，退出循环")
}

func SelectDemo() {
	fmt.Println("=== Select多路复用示例 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "来自ch1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "来自ch2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("ch1收到:", msg1)
		case msg2 := <-ch2:
			fmt.Println("ch2收到:", msg2)
		}
	}
}

func SelectTimeout() {
	fmt.Println("=== Select超时示例 ===")

	ch := make(chan string)

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "延迟消息"
	}()

	select {
	case msg := <-ch:
		fmt.Println("收到:", msg)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("超时! 200ms内未收到消息")
	}
}

func SelectNonBlocking() {
	fmt.Println("=== 非阻塞Select示例 ===")

	ch := make(chan int, 1)

	select {
	case ch <- 42:
		fmt.Println("发送成功")
	default:
		fmt.Println("发送失败，通道已满")
	}

	select {
	case val := <-ch:
		fmt.Println("接收成功:", val)
	default:
		fmt.Println("接收失败，通道为空")
	}

	select {
	case val := <-ch:
		fmt.Println("再次接收:", val)
	default:
		fmt.Println("通道已空，无数据可接收")
	}
}

func NilChannel() {
	fmt.Println("=== Nil通道示例 ===")

	var nilCh chan int
	fmt.Println("nil通道会永久阻塞:")
	go func() {
		fmt.Println("从nil通道接收: <-nilCh 会永久阻塞")
		val := <-nilCh
		fmt.Println("val:", val)
		fmt.Println("从nil通道接收: 程序执行结束")
	}()
	go func() {
		fmt.Println("向nil通道发送: nilCh <- x 会永久阻塞")
		nilCh <- 1
		fmt.Println("向nil通道发送: 程序执行结束")
	}()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("关闭nil通道: close(nilCh) 会panic")

	ch := make(chan int, 1)
	ch <- 1

	ch = nil
	fmt.Println("\n将通道设为nil后，select会忽略该case:")

	select {
	case v := <-ch:
		fmt.Println("收到:", v)
	default:
		fmt.Println("nil通道被忽略，执行default")
	}
}

func WorkerPool() {
	fmt.Println("=== Worker Pool模式示例 ===")

	jobs := make(chan int, 5)
	results := make(chan int, 5)

	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range jobs {
				fmt.Printf("Worker %d 处理任务 %d\n", id, j)
				results <- j * 2
			}
		}(w)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for r := range results {
		fmt.Printf("结果: %d\n", r)
	}
}

func PingPong() {
	fmt.Println("=== Ping-Pong模式示例 ===")

	ping := make(chan string)
	pong := make(chan string)

	go func() {
		for i := 0; i < 3; i++ {
			msg := <-ping
			fmt.Println("Player A收到:", msg, "发出 Pong")
			pong <- "Pong"
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			msg := <-pong
			fmt.Println("Player B收到:", msg, "发出 Ping")
			ping <- "Ping"
		}
	}()

	ping <- "Ping"
	time.Sleep(100 * time.Millisecond)
}
