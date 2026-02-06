package main
import ("fmt";"strconv")
//Len:计算字符串长度
//s:输入字符串
//count:返回字符串长度
func Len(s string) int {
	count := 0
	for range s {
		count++
	}
	return count
}
//Huiwen:计算是否为回文数值
//s:数值
func Huiwen(s int) bool{
	rev := strconv.Itoa(s)
	for i := 0; i < Len(rev)/2; i++ {
		if rev[i] != rev[Len(rev)-i-1] {
			return false
		}
	}
	return true

}





func main(){
	// 打印小明信息
	var name string = "小明"
	var age int = 18
	var gender string = "男" 
	fmt.Printf("姓名：%s，年龄：%d，性别：%s\n", name, age, gender)
	fmt.Println(Len("hello,world"))


}
