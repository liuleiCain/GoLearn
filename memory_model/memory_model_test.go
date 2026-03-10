package memorymodel

import (
	"testing"
)

func TestDataRace(t *testing.T) {
	DataRace()
}

func TestAtomicOperation(t *testing.T) {
	AtomicOperation()
}

func TestMutexLock(t *testing.T) {
	MutexLock()
}

func TestHappensBefore(t *testing.T) {
	HappensBefore()
}

func TestChannelSynchronization(t *testing.T) {
	ChannelSynchronization()
}

func TestOnceInitialization(t *testing.T) {
	OnceInitialization()
}

func TestWaitGroupSynchronization(t *testing.T) {
	WaitGroupSynchronization()
}

func TestMemoryOrdering(t *testing.T) {
	MemoryOrdering()
}

func TestRWMutexExample(t *testing.T) {
	RWMutexExample()
}

func TestAtomicValue(t *testing.T) {
	AtomicValue()
}
