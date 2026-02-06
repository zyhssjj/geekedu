package main

import (
	_"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"strconv"
	
)


func main() {
	// 从命令行参数获取起始数字和结束数字并验证
	if len(os.Args) < 3 {
		fmt.Println("请提供起始数字和结束数字作为命令行参数")
		return
	}
	n1, n2 := os.Args[1], os.Args[2]
	n1Int, err1 := strconv.Atoi(n1)
	if err1 != nil {
		fmt.Println("无效的起始数字")
		return
	}
	n2Int, err2 := strconv.Atoi(n2)
	if err2 != nil {
		fmt.Println("无效的结束数字")
		return
	}
	
	if n1Int > n2Int {
		fmt.Println("起始数字不能大于结束数字")
		return
	}
	if n1Int <= 0 || n2Int <= 0 {
		fmt.Println("数字必须大于0")
		return
	}
	//执行函数并计算时间
	startTime := time.Now()
	countLen, primes, _ :=  CountPrimes(n1Int, n2Int)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	// 输出结果
	fmt.Printf("在范围 [%d, %d] 内共有 %d 个素数。\n", n1Int, n2Int, countLen)
	fmt.Printf("计算耗时: %s\n", elapsedTime)
	// 替换固定文件名，用n1Int和n2Int动态生成
filename := fmt.Sprintf("zhangyuhao_primes_%d_%d.txt", n1Int, n2Int)
file, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
    //处理文件打开错误
	if err != nil{
		fmt.Printf("无法打开文件: %v\n", err)
		os.Exit(1)
	}

	defer file.Close()
	// 将素数写入文件，空格分隔
	var buf strings.Builder
	for k, p := range primes {
		if k > 0 {
			buf.WriteString(" ")

	}
	buf.WriteString(fmt.Sprintf("%d", p))
}
	_, err = file.WriteString(buf.String())
	if err != nil{
		fmt.Printf("无法写入文件: %v\n", err)
		os.Exit(1)
	}
}
