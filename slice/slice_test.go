package slice

import "testing"

func TestSliceStructure(t *testing.T) {
	SliceStructure()
}

func TestSliceCreate(t *testing.T) {
	SliceCreate()
}

func TestAppendDemo(t *testing.T) {
	AppendDemo()
}

func TestGrowthRule(t *testing.T) {
	GrowthRule()
}

func TestCopyDemo(t *testing.T) {
	CopyDemo()
}

func TestSliceExpression(t *testing.T) {
	SliceExpression()
}

func TestSliceTricks(t *testing.T) {
	SliceTricks()
}

func TestSliceMemory(t *testing.T) {
	SliceMemory()
}

func TestSlicePointer(t *testing.T) {
	SlicePointer()
}

func BenchmarkSliceNoPreAlloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s := make([]int, 0)
		for i := 0; i < 1000000; i++ {
			s = append(s, i)
		}
	}
}

func BenchmarkSlicePreAlloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s := make([]int, 0, 1000000)
		for i := 0; i < 1000000; i++ {
			s = append(s, i)
		}
	}
}
