// 过滤器模块：3个goroutine并发过滤，筛选偶数ID日志，支持停止信号
package main

import (
	"fmt"
	"sync"
)

// FilterLogs 过滤器：消费logChannel，筛选偶数ID日志到filteredChannel
// 输入参数：
// logChannel：日志通道
// filteredChannel：筛选后的日志通道
// stopChan：停止信号通道
// wg：同步等待组
// filterID：过滤器ID（用于日志区分）
func FilterLogs(
	logChannel <-chan LogEntry,
	filteredChannel chan<- LogEntry,
	stopChan chan struct{}, // 修正1：从 chan<- struct{} → chan struct{}（双向通道）
	wg *sync.WaitGroup,
	filterID int, // 过滤器ID（用于日志区分）
) {
	defer func() {
		wg.Done() // 过滤器完成，通知WaitGroup
		fmt.Printf("过滤器%d：已退出\n", filterID)
	}()

	fmt.Printf("过滤器%d：启动，开始处理日志\n", filterID)

	// 循环接收logChannel日志（通道关闭则退出循环）
	for log := range logChannel {
		// 监听停止信号：收到信号立即退出（现在stopChan是双向通道，可接收）
		select {
		case <-stopChan: // 停止信号处理
			fmt.Printf("过滤器%d：收到停止信号，终止处理\n", filterID)
			return
		default:
			// 检查日志ID：偶数→重要日志，发送到filteredChannel
			if log.ID%2 == 0 {
				fmt.Printf("过滤器%d：筛选重要日志 ID：%d（偶数）\n", filterID, log.ID)
				filteredChannel <- log

				// 特殊处理：ID=50时，向stopChan发送停止信号（所有goroutine退出）
				if log.ID == 50 {
					fmt.Printf("过滤器%d：检测到ID=50，发送停止信号\n", filterID)
					// 非阻塞发送，避免多个过滤器重复发送导致阻塞
					select {
					case stopChan <- struct{}{}: // 双向通道可正常发送
					default:
					}
					return // 发送信号后立即退出当前过滤器
				}
			} else {
				fmt.Printf("过滤器%d：忽略普通日志 ID：%d（奇数）\n", filterID, log.ID)
			}
		}
	}
}



// CloseFilteredChannel 协调goroutine：等待所有过滤器完成后，关闭filteredChannel
// 输入参数：
// filteredChannel：筛选后日志通道
// wg：过滤器WaitGroup
// stopChan：停止信号通道
// 功能说明：等待所有过滤器完成后关闭filteredChannel，或收到停止信号时立即关闭
func CloseFilteredChannel(filteredChannel chan<- LogEntry, wg *sync.WaitGroup, stopChan <-chan struct{}) {
	
	doneChan := make(chan struct{})
	go func() {
		wg.Wait()         // 等待所有过滤器完成
		close(doneChan)   // 完成后关闭临时通道
	}()

	// 监听：要么收到停止信号，要么所有过滤器完成
	select {
	case <-stopChan:
		fmt.Println("协调器：收到停止信号，关闭filteredChannel")
	case <-doneChan: // 接收临时通道信号（替代 wg.WaitChan()）
		fmt.Println("协调器：所有过滤器处理完毕，关闭filteredChannel")
	}
	close(filteredChannel)
}