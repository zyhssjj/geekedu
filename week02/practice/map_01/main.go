package main
import ("fmt")
//打印成绩单
func main() {
    mark := map[string]int{
		"小明": 60,
		"小王": 70,
		"张三": 95,
		"李四": 98,
		"王五": 100,
		"张伟": 88,
	}
	for name, score := range mark {
		fmt.Printf("%s的成绩是%d分\n", name, score)
	}
	



}
