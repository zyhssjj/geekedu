// 存储模块：接收过滤后的重要日志，打印模拟存储过程
package main

import "fmt"

// StoreLogs 存储者：消费filteredChannel中的重要日志，打印模拟存储
// stopChan：接收停止信号，收到后立即退出
// 输入参数：
// filteredChannel：筛选后日志通道
// stopChan：停止信号通道
// 功能说明：接收重要日志并打印，或收到停止信号时退出
func StoreLogs(filteredChannel <-chan LogEntry, stopChan <-chan struct{}) {
	defer fmt.Println("存储者：已退出")
	fmt.Println("存储者：启动，开始接收重要日志")

	// 循环接收重要日志（通道关闭或收到停止信号则退出）
	for {
		select {
		case <-stopChan:
			fmt.Println("存储者：收到停止信号，终止存储")
			return
		case log, ok := <-filteredChannel:
			if !ok {
				// filteredChannel关闭，说明无更多日志
				fmt.Println("存储者：filteredChannel已关闭，无更多日志")
				return
			}
			// 打印模拟存储（实际场景可写入文件/数据库）
			fmt.Printf("存储者：存储重要日志 - ID：%d，内容：%s\n", log.ID, log.Content)
		}
	}
}