package main
import ("fmt")
//切片练习
func main(){
	slice := []int{1,2,3,4,5,6,7,8,9,10}
	fmt.Println("步骤1结果：", slice)
	slice=slice[2:7]
	fmt.Println("步骤2结果：", slice)
	slice=append(slice,11,12,13)
	fmt.Println("步骤3结果：", slice)
	slice=append(slice[:4],slice[5:]...)
	fmt.Println("步骤4结果：", slice)
	for i := range slice {
		slice[i] *= 2
	}
	fmt.Println("步骤5结果：", slice)
	fmt.Println("最终结果：", slice,cap(slice))
}