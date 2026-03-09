package _break

import "testing"

func TestBasicBreak(t *testing.T) {
	BasicBreak()
}

func TestBreakInSwitch(t *testing.T) {
	BreakInSwitch()
}

func TestBreakVsContinue(t *testing.T) {
	BreakVsContinue()
}

func TestBreakForSelect(t *testing.T) {
	BreakForSelect()
}

func TestBreakForSelectLabel(t *testing.T) {
	BreakForSelectLabel()
}

func TestBreakNestedLoop(t *testing.T) {
	BreakNestedLoop()
}

func TestBreakWithLabelNames(t *testing.T) {
	BreakWithLabelNames()
}

func TestBreakVsReturn(t *testing.T) {
	BreakVsReturn()
}

func TestInfiniteLoopBreak(t *testing.T) {
	InfiniteLoopBreak()
}

func TestBreakPitfalls(t *testing.T) {
	BreakPitfalls()
}
