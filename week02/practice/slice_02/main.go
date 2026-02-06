package main

import (
	"fmt"

)

//切片练习
func main(){
	//得到合并后的切片
	slice_1 := []int{1,2,3,4}
	slice_2 := []int{3,4,5,6}
	combinedSlice := append(slice_1,slice_2...)

	uniqueMap := make(map[int]bool)
	uniqueSlice := make([]int, 0)
	//切片去重
	for _, num := range combinedSlice {
		if !uniqueMap[num] { 
			uniqueMap[num] = true 
			uniqueSlice = append(uniqueSlice, num) 
		}
	}
	fmt.Println("去重后的切片 uniqueSlice：", uniqueSlice)
	
}