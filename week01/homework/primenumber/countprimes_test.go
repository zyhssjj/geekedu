package main
import ("testing"
)
// TestCountPrimes:测试CountPrimes函数
//t:测试框架提供的测试对象

func TestCountPrimes(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name        string
		input1       int
		input2       int
		expected    int
		expectedErr bool
	}{{name: "起始数字大于结束数字", input1: 10,input2: 5, expected: 0, expectedErr: true},
		{name: "包含负数", input1: -5, input2: 5,  expected: 0, expectedErr: true},
		{name: "正常范围", input1: 10, input2: 20, expected: 4, expectedErr: false},
		{name: "范围内无质数", input1: 14, input2:15, expected: 0, expectedErr: false},
	}
	// 遍历测试用例
	for _, tt := range tests { 

		t.Run(tt.name, func(t *testing.T) {
			count, _, err := CountPrimes(tt.input1, tt.input2)
			if (err != nil) != tt.expectedErr {
				t.Errorf("CountPrimes() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if count != tt.expected {
				t.Errorf("CountPrimes() = %v, want %v", count, tt.expected)
			}

		})
	}

}