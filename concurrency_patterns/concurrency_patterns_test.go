package concurrency_patterns

import "testing"

func TestWorkerPoolPattern(t *testing.T) {
	WorkerPoolPattern()
}

func TestPipelinePattern(t *testing.T) {
	PipelinePattern()
}

func TestFanOutFanIn(t *testing.T) {
	FanOutFanIn()
}

func TestTimeoutPattern(t *testing.T) {
	TimeoutPattern()
}

func TestCancellationPattern(t *testing.T) {
	CancellationPattern()
}

func TestSemaphorePattern(t *testing.T) {
	SemaphorePattern()
}

func TestRateLimitingPattern(t *testing.T) {
	RateLimitingPattern()
}

func TestGracefulShutdown(t *testing.T) {
	GracefulShutdown()
}
