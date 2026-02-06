package main

import (
    "flag"
    "fmt"
	"os"
    
)

func main() {
    // 1. 定义三个变量，用于接收命令行参数的值
    var filePath string 
    var operation string
    var outputPath string 

    // 2. 使用 flag.StringVar 定义命令行参数
    // 参数说明：(变量指针, "参数名", "默认值", "参数说明")
    flag.StringVar(&filePath, "file", "", "指定要处理的文件路径")
    flag.StringVar(&operation, "operation", "", "指定操作类型 (count, upper, convert)")
    flag.StringVar(&outputPath, "output", "", "指定输出文件路径")

    // 3. 解析命令行参数
    flag.Parse()

	
	
    // 4. 根据 operation 的值，执行不同的函数
    switch operation {
    case "count":
        // 调用统计字符数的函数
        fmt.Printf("统计文件字符数\n")
    case "upper":
        // 调用转换为大写的函数
        // 这个操作需要输出文件路径
       
        fmt.Printf("转换为大写\n")
    case "convert":
        // 调用你自定义的数字转换函数
     
        fmt.Printf("转换数字格式\n")
    default:
        fmt.Printf("错误：不支持的操作类型 '%s'。\n", operation)
		fmt.Println("请输入正确的操作类型。")
        os.Exit(1)
    }

	fmt.Printf("filepath: %s\n", filePath)
	fmt.Printf("outputpath: %s\n", outputPath)
}
