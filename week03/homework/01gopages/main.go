// 程序入口：任务调度、并发控制、结果统计、耗时计算
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 1. 生成目标链接列表（第8~48个网页）
	links := GenerateLinks()

	// 2. 初始化统计变量（需并发安全，用互斥锁保护）
	var (
		wg         sync.WaitGroup       // 等待所有goroutine完成
		mutex      sync.Mutex           // 保护共享变量的互斥锁
		totalCount int                  // 图片总数（不去重）
		uniqueMap  = make(map[string]struct{}) // 存储唯一图片标识（去重）
		sem        = make(chan struct{}, 5) // 信号量，控制最大并发数为5
	)

	// 3. 记录任务开始时间
	startTime := time.Now()

	// 4. 并发爬取所有链接
	for _, link := range links {
		// 信号量P操作：获取并发许可（最大5个）
		sem <- struct{}{}

		wg.Add(1)
		go func(url string) { // 传入当前链接，避免循环变量引用问题
			defer func() {
				// 信号量V操作：释放并发许可
				<-sem
				wg.Done()
			}()

			// 爬取当前网页的图片标识列表
			imgIDs, err := CrawlPage(url)
			if err != nil {
				fmt.Printf("警告：%v\n", err)
				return
			}

			// 统计结果（加锁保护共享变量）
			mutex.Lock()
			defer mutex.Unlock()

			totalCount += len(imgIDs) // 累加不去重总数
			for _, id := range imgIDs {
				uniqueMap[id] = struct{}{} // 存入map实现去重
			}
		}(link) // 传入当前循环的link
	}

	// 5. 等待所有爬取任务完成
	wg.Wait()

	// 6. 计算总耗时（毫秒）
	costTime := time.Since(startTime).Milliseconds()

	// 7. 提取去重后的图片数量
	uniqueCount := len(uniqueMap)

	// 8. 按要求格式打印结果（main函数最后一行必须是此句）
	fmt.Printf("执行任务总耗时（毫秒）：%v，图片总数（不去重）：%d，图片总数（去重）：%d\n", costTime, totalCount, uniqueCount)
}