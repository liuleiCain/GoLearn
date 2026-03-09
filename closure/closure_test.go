package closure

import "testing"

func TestBasicClosure(t *testing.T) {
	BasicClosure()
}

func TestClosureFactory(t *testing.T) {
	ClosureFactory()
}

func TestLoopTrap(t *testing.T) {
	LoopTrap()
}

func TestGoroutineTrap(t *testing.T) {
	GoroutineTrap()
}

func TestClosureState(t *testing.T) {
	ClosureState()
}

func TestClosureRecursion(t *testing.T) {
	ClosureRecursion()
}

func TestClosurePerformance(t *testing.T) {
	ClosurePerformance()
}

func TestClosureDefer(t *testing.T) {
	ClosureDefer()
}

func TestClosureTimer(t *testing.T) {
	ClosureTimer()
}
