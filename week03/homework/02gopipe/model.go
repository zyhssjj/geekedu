// 数据模型模块：定义日志结构体，供其他模块复用
package main

// LogEntry 日志条目结构体（按作业要求包含 ID 和 Content）
type LogEntry struct {
	ID      int    // 日志唯一ID（1-100）
	Content string // 日志内容
}