package channel

import "testing"

func TestUnbufferedChannel(t *testing.T) {
	UnbufferedChannel()
}

func TestBufferedChannel(t *testing.T) {
	BufferedChannel()
}

func TestChannelDirection(t *testing.T) {
	ChannelDirection()
}

func TestCloseChannel(t *testing.T) {
	CloseChannel()
}

func TestSelectDemo(t *testing.T) {
	SelectDemo()
}

func TestSelectTimeout(t *testing.T) {
	SelectTimeout()
}

func TestSelectNonBlocking(t *testing.T) {
	SelectNonBlocking()
}

func TestNilChannel(t *testing.T) {
	NilChannel()
}

func TestWorkerPool(t *testing.T) {
	WorkerPool()
}

func TestPingPong(t *testing.T) {
	PingPong()
}
