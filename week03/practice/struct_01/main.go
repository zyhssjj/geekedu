package main
import ("fmt"
	"encoding/json"
	)
// Person:表示一个人
type Person struct {
	Name string
	Age int
	Email string
}
// NewPerson:创建并返回一个新的Person实例
// name:姓名
// age:年龄
// email:电子邮件
// *Person:指向Person实例的指针
func NewPerson(name string, age int, email string) *Person {
	return &Person{
		Name: name,
		Age: age,
		Email: email,
	}
}
// PrintPerson:打印Person的信息
// p:*Person:指向Person实例的指针
func PrintPerson(p *Person) {
	json, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("序列化Person时出错: %v\n", err)
		return
	}
	fmt.Println("Person:", string(json))
	fmt.Printf("Name: %s, Age: %d, Email: %s\n", p.Name, p.Age, p.Email)
}

func main() {
	p1 := NewPerson("Alice", 30, "dadDda")
	PrintPerson(p1)
}