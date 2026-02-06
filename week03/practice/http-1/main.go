package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Todo 结构体定义了 JSONPlaceholder API 返回的单个任务格式
type Todo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	// 1. 定义请求的 URL
	url := "https://jsonplaceholder.typicode.com/todos"

	// 2. 创建一个 HTTP 客户端
	client := &http.Client{}

	// 3. 创建一个 GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		os.Exit(1)
	}

	// 4. 设置请求头（例如，模拟浏览器访问）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 5. 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		os.Exit(1)
	}

	// 6. 确保在函数结束时关闭响应体
	defer resp.Body.Close()

	// 7. 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求失败，状态码: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// 8. 解析 JSON 响应体
	var todos []Todo
	err = json.NewDecoder(resp.Body).Decode(&todos)
	if err != nil {
		fmt.Printf("解析 JSON 失败: %v\n", err)
		os.Exit(1)
	}

	// 9. 打印每个任务的信息
	fmt.Printf("成功获取 %d 个任务:\n", len(todos))
	for i, todo := range todos {
		fmt.Printf("\n任务 #%d:\n", i+1)
		fmt.Printf("  User ID:    %d\n", todo.UserId)
		fmt.Printf("  ID:         %d\n", todo.Id)
		fmt.Printf("  Title:      %s\n", todo.Title)
		fmt.Printf("  Completed:  %v\n", todo.Completed)
	}
}