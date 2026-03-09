package context

import "testing"

func TestBackgroundDemo(t *testing.T) {
	BackgroundDemo()
}

func TestWithCancelDemo(t *testing.T) {
	WithCancelDemo()
}

func TestWithTimeoutDemo(t *testing.T) {
	WithTimeoutDemo()
}

func TestWithDeadlineDemo(t *testing.T) {
	WithDeadlineDemo()
}

func TestWithValueDemo(t *testing.T) {
	WithValueDemo()
}

func TestChainContextDemo(t *testing.T) {
	ChainContextDemo()
}

func TestPropagationDemo(t *testing.T) {
	PropagationDemo()
}

func TestErrDemo(t *testing.T) {
	ErrDemo()
}
