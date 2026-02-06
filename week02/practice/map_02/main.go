package main

import "fmt"

// findMostFrequentChar:统计字符串中每个字符的出现次数，并返回出现次数最多的字符
//s:输入字符串
func findMostFrequentChar(s string) byte {
    // 用map统计每个字符的出现次数，键为字符，值为出现次数
    charCount := make(map[rune]int)
    for _, c := range s {
        charCount[c]++
    }

    // 遍历map，找到出现次数最多的字符
    maxCount := 0
    var mostFrequentChar byte
    for c, count := range charCount {
        if count > maxCount {
            maxCount = count
            mostFrequentChar = byte(c)
        }
    }
    return mostFrequentChar
}

func main() {
    // 测试findMostFrequentChar函数
    testStr := "aabbbcc111122"
    result := findMostFrequentChar(testStr)
    fmt.Printf("出现次数最多的字符是：%c\n", result)
}