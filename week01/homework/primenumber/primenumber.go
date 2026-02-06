package main
import "errors"
// IsPrime:判断一个整数是否为质数
//n:输入整数
//返回值:如果是质数返回true，否则返回false
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i:=2; i<=n/2; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
// CountPrimes:计算在给定范围内的质数个数
//n1:范围起始数字
//n2:范围结束数字
//返回值:质数个数，质数切片，错误信息
func CountPrimes(n1 int, n2 int) (int, []int,error) {
	if n1 > n2 {
		return 0, nil, errors.New("起始数字不能大于结束数字")

	}
	if n1 <= 0 || n2 <= 0 {
		return 0, nil, errors.New("数字必须大于0")
	}
	 primes := make([]int,0,1000)

	for n := n1 ; n <= n2; n++ {
		if IsPrime(n) {
			primes = append (primes, n)
		}



	}
	return len(primes),primes,nil
}