// 程序入口：初始化通道、启动所有goroutine、等待退出
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 1. 初始化通道（缓冲通道避免阻塞，提高并发效率）
	const (
		logChanBufSize    = 10  // 日志通道缓冲大小
		filteredChanBufSize = 5 // 过滤后通道缓冲大小
	)
	logChannel := make(chan LogEntry, logChanBufSize)
	filteredChannel := make(chan LogEntry, filteredChanBufSize)
	stopChan := make(chan struct{}, 1) // 停止信号通道（缓冲1，避免发送阻塞）

	// 2. 初始化WaitGroup（用于等待3个过滤器完成）
	var filterWG sync.WaitGroup

	// 3. 启动生产者goroutine（1个）
	go ProduceLogs(logChannel, stopChan)

	// 4. 启动过滤器goroutine（3个）
	filterCount := 3
	filterWG.Add(filterCount)
	for i := 1; i <= filterCount; i++ {
		go FilterLogs(logChannel, filteredChannel, stopChan, &filterWG, i)
	}

	// 5. 启动协调goroutine：等待过滤器完成后关闭filteredChannel
	go CloseFilteredChannel(filteredChannel, &filterWG, stopChan)

	// 6. 启动存储者goroutine（1个）
	go StoreLogs(filteredChannel, stopChan)

	// 7. 等待停止信号（主线程阻塞，直到收到信号后退出）
	<-stopChan
	fmt.Println("\n主线程：收到停止信号，等待所有goroutine退出...")

	// 8. 等待所有goroutine优雅退出（给短暂时间清理）
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n程序：所有goroutine已退出，程序结束")
}