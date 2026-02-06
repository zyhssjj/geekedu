package models
//定义学生结构体和学生列表
type Student struct {
	ID int  `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Grade string `json:"grade"`
}
var StudentList []Student