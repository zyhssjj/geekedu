package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// 1. 主函数：程序入口 
func main() {
	// 获取输入文件名
	inputFileName := getInputFileName()
	fmt.Printf("准备处理文件: %s\n", inputFileName)

	// 确保输出目录存在
	outputDir := "week01/homework/calc/results/"
	err := ensureDir(outputDir)
	if err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 构建输出文件名
	outputFileName := outputDir + getBaseName(inputFileName) + "_result.txt"

	// 执行核心处理逻辑
	err = processCalculations(inputFileName, outputFileName)
	if err != nil {
		fmt.Printf("处理文件时发生错误: %v\n", err)
	} else {
		fmt.Printf("处理完成！结果已保存至: %s\n", outputFileName)
	}
}

// 2.getInputFileName:获取输入文件名，默认为calculations.txt
//   string:输入的文件名


func getInputFileName() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return "calculations.txt"
}

// 3.ensureDir:确保目录存在，不存在则创建
// dirName:目录名
// error:错误信息
func ensureDir(dirName string) error {
	return os.MkdirAll(dirName, 0755)
}

// 4.getBaseName:获取文件名（不带扩展名）
// path:文件路径
// string:不带扩展名的文件名
func getBaseName(path string) string {
	ext := ""
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			ext = path[i:]
			break
		}
	}
	return path[:len(path)-len(ext)]
}

// 5. processCalculations:处理计算文件的主要逻辑
// inputFile:输入文件名
// outputFile:输出文件名
// error:错误信息
func processCalculations(inputFile, outputFile string) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %w", err)
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		resultLine, err := calculateLine(line) // 调用修改后的 calculateLine 函数
		if err != nil {
			fmt.Printf("忽略格式错误的行: '%s'。原因: %v\n", line, err)
			continue
		}

		_, err = writer.WriteString(resultLine + "\n")
		if err != nil {
			fmt.Printf("写入行 '%s' 失败: %v\n", resultLine, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时出错: %w", err)
	}

	return nil
}

// 6. calculateLine:计算单行表达式的结果
// line:输入行
// string:计算结果行，错误信息

func calculateLine(line string) (string, error) {
	cleanLine := strings.ReplaceAll(line, " ", "")

	var opIndex int
	var operator string

	
	for i, c := range cleanLine {
		if c == '+' || c == '-' || c == '*' || c == '/' {
			opIndex = i
			operator = string(c)
			break
		}
	}

	if operator == "" {
		return "", fmt.Errorf("未找到 '+'、'-'、'*' 或 '/'")
	}

	num1Str := cleanLine[:opIndex]
	num2Str := cleanLine[opIndex+1:]

	if num1Str == "" || num2Str == "" {
		return "", fmt.Errorf("运算符前后必须有数字")
	}

	num1, err := strconv.ParseFloat(num1Str, 64)
	if err != nil {
		return "", fmt.Errorf("'%s' 不是有效的数字", num1Str)
	}
	num2, err := strconv.ParseFloat(num2Str, 64)
	if err != nil {
		return "", fmt.Errorf("'%s' 不是有效的数字", num2Str)
	}

	var result float64
	switch operator {
	case "+":
		result = num1 + num2
	case "-": 
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/": 
		// 检查除数是否为零
		if num2 == 0 {
			return "", fmt.Errorf("除法运算中，除数不能为零")
		}
		result = num1 / num2
	}

	resultStr := formatNumber(result)

	return fmt.Sprintf("%s%s%s=%s", num1Str, operator, num2Str, resultStr), nil
}

// 7. formatNumber:格式化数字，去除多余的小数点和零
// num:输入数字
// string:格式化后的字符串
func formatNumber(num float64) string {
	if num == float64(int(num)) {
		return strconv.Itoa(int(num))
	}
	return strconv.FormatFloat(num, 'f', -1, 64)
}