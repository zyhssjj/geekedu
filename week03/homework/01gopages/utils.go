// 工具模块：生成第8到第48个网页的链接（共41个）
package main
import ("fmt")

// GenerateLinks 生成目标链接列表：https://study-test.sixue.work/html/8.html 到 48.html
func GenerateLinks() []string {
	var links []string
	// 循环生成第8~48个链接（含首尾）
	for i := 8; i <= 48; i++ {
		link := fmt.Sprintf("https://study-test.sixue.work/html/%d.html", i)
		links = append(links, link)
	}
	return links
}