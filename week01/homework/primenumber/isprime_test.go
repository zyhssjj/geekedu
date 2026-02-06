package main

import (
	"testing"
)
// TestIsPrime:测试IsPrime函数
//t:测试框架提供的测试对象
func TestIsPrime(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
		wantErr  bool
	}{{name: "负数", input: -5, expected: false},
		{name: "0和1", input: 1, expected: false},
		{name: "2和3", input: 2, expected: true},
		{name: "合数", input: 4, expected: false},
		{name: "质数", input: 13, expected: true},
	}
	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPrime(tt.input)

			if tt.expected != result {
				t.Errorf("IsPrime(%d) = %v, 期望结果 = %v", tt.input, result, tt.expected)
			}
			
		})
	}

}
