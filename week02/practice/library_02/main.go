package main

import (
    "flag"
    "fmt"
    "os"
    "time"
)

// 1. 定义日志级别类型和常量
type LogLevel int

const (
    INFO LogLevel = iota
    WARNING
    ERROR
)

// 将 LogLevel 转换为字符串的辅助函数
func (l LogLevel) String() string {
    switch l {
    case INFO:
        return "INFO"
    case WARNING:
        return "WARNING"
    case ERROR:
        return "ERROR"
    default:
        return "UNKNOWN"
    }
}

// 2. 定义用于存储命令行参数的变量
var logLevel LogLevel
var outputPath string

// 6. 定义一个通用的日志写入函数
func log(level LogLevel, message string) {
    // 如果日志级别低于配置的级别，则不输出
    if level < logLevel {
        return
    }

    // 生成时间戳
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logMessage := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level.String(), message)

    // 7. 根据输出位置，决定写入文件还是控制台
    if outputPath == "console" {
        fmt.Print(logMessage)
    } else {
        file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Printf("打开日志文件失败: %v\n", err)
            os.Exit(1)
        }
        defer file.Close()

        _, err = file.WriteString(logMessage)
        if err != nil {
            fmt.Printf("写入日志文件失败: %v\n", err)
            os.Exit(1)
        }
    }
}

func main() {
    // 定义一个字符串变量，用于临时接收日志级别参数
    var levelStr string

    // 3. 使用 flag 包定义命令行参数
    flag.StringVar(&levelStr, "level", "info", "指定日志级别 (info, warning, error)")
    flag.StringVar(&outputPath, "output", "console", "指定输出位置 (console 或文件路径)")

    // 4. 解析命令行参数
    flag.Parse()

    // 5. 将字符串类型的日志级别转换为我们定义的 LogLevel 类型
    switch levelStr {
    case "info":
        logLevel = INFO
    case "warning":
        logLevel = WARNING
    case "error":
        logLevel = ERROR
    default:
        fmt.Printf("错误：不支持的日志级别 '%s'。\n", levelStr)
        os.Exit(1)
    }

    // 8. 模拟日志输出
    log(INFO, "程序启动成功")
    log(WARNING, "发现一个潜在问题")
    log(ERROR, "数据库连接失败")

    fmt.Println("日志已按配置输出。")
}