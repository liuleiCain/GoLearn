package goroutine

import "testing"

func TestBasicGoroutine(t *testing.T) {
	BasicGoroutine()
}

func TestWaitGroupDemo(t *testing.T) {
	WaitGroupDemo()
}

func TestGoroutineLeakDemo(t *testing.T) {
	GoroutineLeakDemo()
}

func TestGoroutineNumberDemo(t *testing.T) {
	GoroutineNumberDemo()
}

func TestGoschedDemo(t *testing.T) {
	GoschedDemo()
}

func TestGOMAXPROCSDemo(t *testing.T) {
	GOMAXPROCSDemo()
}

func TestOnceDemo(t *testing.T) {
	OnceDemo()
}

func TestMutexDemo(t *testing.T) {
	MutexDemo()
}

func TestRWMutexDemo(t *testing.T) {
	RWMutexDemo()
}

func TestAtomicDemo(t *testing.T) {
	AtomicDemo()
}

func TestContextCancel(t *testing.T) {
	ContextCancel()
}

func TestContextTimeout(t *testing.T) {
	ContextTimeout()
}

func TestWorkerPool(t *testing.T) {
	WorkerPool()
}

func TestPipeline(t *testing.T) {
	Pipeline()
}

func TestFanOutFanIn(t *testing.T) {
	FanOutFanIn()
}

func TestTimeoutPattern(t *testing.T) {
	TimeoutPattern()
}

func TestGracefulShutdown(t *testing.T) {
	GracefulShutdown()
}

func TestGoroutineLocal(t *testing.T) {
	GoroutineLocal()
}

func TestSelectDemo(t *testing.T) {
	SelectDemo()
}

func TestTimerDemo(t *testing.T) {
	TimerDemo()
}

func TestTickerDemo(t *testing.T) {
	TickerDemo()
}
