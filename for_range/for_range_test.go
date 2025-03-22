package for_range

import (
	"testing"
)

func TestAnalysisBasic(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"整型和字符串测试"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AnalysisBasic()
		})
	}
}

func TestAnalysisArr(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"数组测试"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AnalysisArr()
		})
	}
}

func TestAnalysisSlice(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"切片测试"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AnalysisSlice()
		})
	}
}

func TestAnalysisMap(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Map测试"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AnalysisMap()
		})
	}
}

func TestAnalysisRangeSliceIndex(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"遍历切片索引测试"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AnalysisRangeSliceIndex()
		})
	}
}

func TestAnalysisRangeStruct(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"遍历结构体测试"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AnalysisRangeStruct()
		})
	}
}
