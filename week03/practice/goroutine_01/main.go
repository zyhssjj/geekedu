package main

import (
    "fmt"
    "sync"
    "time"
    "math/rand"
)

// 计算一个整数切片的平方和
func calculateSquareSum(slice []int, resultChan chan<- int, wg *sync.WaitGroup) {
    defer wg.Done() // 通知 WaitGroup 该 goroutine 已完成

    sum := 0
    for _, num := range slice {
        sum += num * num
        // 模拟计算耗时（每个元素耗时 1 毫秒）
        time.Sleep(1 * time.Millisecond)
    }

    resultChan <- sum // 将结果发送到通道
}

func main() {
    // 1. 生成一个长度为 1000 的随机整数切片（元素范围：0-99）
    rand.Seed(time.Now().UnixNano()) // 初始化随机数种子
    slice := make([]int, 1000)
    for i := range slice {
        slice[i] = rand.Intn(100)
    }

    // 2. 定义分割参数
    chunkSize := 200 // 每块的大小
    numChunks := len(slice) / chunkSize
    if len(slice)%chunkSize != 0 {
        numChunks++ // 如果不能整除，最后一块会小一些
    }

    // 3. 初始化通道和 WaitGroup
    resultChan := make(chan int, numChunks) // 带缓冲的通道，避免阻塞
    var wg sync.WaitGroup

    // 4. 分割切片并启动 goroutine 计算
    start := 0
    for i := 0; i < numChunks; i++ {
        end := start + chunkSize
        if end > len(slice) {
            end = len(slice) // 最后一块的结束索引不超过切片长度
        }

        wg.Add(1) // 增加 WaitGroup 计数
        // 启动 goroutine 计算当前块的平方和
        go calculateSquareSum(slice[start:end], resultChan, &wg)

        start = end // 更新下一块的起始索引
    }

    // 5. 启动一个 goroutine 等待所有计算完成后关闭通道
    go func() {
        wg.Wait()         // 等待所有 goroutine 完成
        close(resultChan) // 关闭通道，避免死锁
    }()

    // 6. 汇总所有结果
    totalSum := 0
    for sum := range resultChan {
        totalSum += sum
    }

    // 7. 输出最终结果
    fmt.Printf("整个切片的元素平方和为：%d\n", totalSum)
}