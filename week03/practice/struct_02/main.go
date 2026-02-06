package main
import ("fmt"
"encoding/json")


func main() {
// 定义结构体
  type Person struct {
    Name string `json:"name"`
    Age int `json:"age"`
    Email string `json:"email"`
  }
 // 创建JSON字符串
  jsonstr := `{"name": "Jane Smith", "age": 25, "email": "janesmith@example.com"}`
  var p Person
  // 将JSON字符串反序列化到结构体变量p中
  err := json.Unmarshal([]byte(jsonstr), &p)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Printf("反序列化结果为：Name: %s, Age: %d, Email: %s\n", p.Name, p.Age, p.Email)
  //序列化结构体变量p为JSON字符串
  json1, err := json.Marshal(p)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  fmt.Println("JSON:", string(json1))



}
