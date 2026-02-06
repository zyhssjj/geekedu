// 生产者模块：生成100条日志，发送到logChannel，支持优雅退出
package main

import "fmt"

// ProduceLogs 生产者：生成ID 1-100的日志，通过logChannel发送
// stopChan：接收停止信号，收到后立即退出

func ProduceLogs(logChannel chan<- LogEntry, stopChan <-chan struct{}) {
	defer fmt.Println("生产者：已退出")

	// 循环生成100条日志（ID 1-100）
	for id := 1; id <= 100; id++ {
		// 监听停止信号：收到信号立即退出，不再生成日志
		select {
		case <-stopChan:
			fmt.Printf("生产者：收到停止信号，终止日志生成（当前生成到ID：%d）\n", id-1)
			return
		default:
			// 生成日志（Content 模拟日志内容）
			log := LogEntry{
				ID:      id,
				Content: fmt.Sprintf("系统日志 - 操作ID：%d，状态：正常", id),
			}
			// 发送到logChannel（若通道满则阻塞，直到有过滤器接收）
			logChannel <- log
			fmt.Printf("生产者：生成日志 ID：%d\n", id)
		}
	}

	// 生成完毕（未收到停止信号时），关闭logChannel，告知过滤器无更多日志
	fmt.Println("生产者：100条日志生成完毕，关闭logChannel")
	close(logChannel)
}