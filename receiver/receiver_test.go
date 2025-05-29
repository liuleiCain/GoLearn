package learn_receiver_test

import (
	learn_receiver "go-learn/receiver"
	"testing"
)

func TestReceiverTypeCompare(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试receiver普通类型和指针类型比较"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learn_receiver.ReceiverTypeCompare()
		})
	}
}

func TestReceiverListTypeCompare(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试结构体数组和结构体数组指针对receiver方法调用的影响"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learn_receiver.ReceiverListTypeCompare()
		})
	}
}
