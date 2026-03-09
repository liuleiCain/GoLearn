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
