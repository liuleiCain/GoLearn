package pointer

import "testing"

func TestBasicPointer(t *testing.T) {
	BasicPointer()
}

func TestPointerArithmetic(t *testing.T) {
	PointerArithmetic()
}

func TestReceiverCompare(t *testing.T) {
	ReceiverCompare()
}

func TestNilPointer(t *testing.T) {
	NilPointer()
}

func TestPointerToPointer(t *testing.T) {
	PointerToPointer()
}

func TestNewAndMake(t *testing.T) {
	newAndMake()
}

func TestEscapeAnalysis(t *testing.T) {
	EscapeAnalysis()
}

func TestPointerPerformanceTime(t *testing.T) {
	PointerPerformanceTime()
}

func TestPointerPerformanceMemory(t *testing.T) {
	PointerPerformanceMemory()
}

func TestPointerSafety(t *testing.T) {
	PointerSafety()
}
