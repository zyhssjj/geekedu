// 爬取模块：负责单个网页的请求、图片标识提取
package main

import (
	"io"
	"net/http"
	"regexp"
	"fmt"
)

// 预编译正则表达式：匹配 "图片1"、"图片 2" 等格式（提取图片标识）
// 正则说明：
// 图片：固定前缀
// \s?：可选空格（处理有空格和无空格的情况）
// (\d+)：匹配数字（图片编号），捕获组用于提取
var imgRegex = regexp.MustCompile(`图片\s?(\d+)`)

// CrawlPage 爬取单个网页，返回该页面的图片标识列表（不去重）
// 参数 link：网页链接
// 返回值 imgIDs：图片标识列表，error：错误信息
func CrawlPage(link string) ([]string, error) {
	// 发送HTTP GET请求
	resp, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("请求链接失败：%s，错误：%v", link, err)
	}
	defer resp.Body.Close() // 确保响应体关闭，避免资源泄露

	// 读取网页内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取网页内容失败：%s，错误：%v", link, err)
	}

	// 匹配所有图片标识
	matches := imgRegex.FindAllStringSubmatch(string(body), -1)
	var imgIDs []string
	for _, match := range matches {
		// 统一格式为 "图片X"（去掉空格，便于去重）
		imgID := fmt.Sprintf("图片%s", match[1])
		imgIDs = append(imgIDs, imgID)
	}

	return imgIDs, nil
}