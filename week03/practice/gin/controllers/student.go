package controllers
import (
"strconv"
"fmt"
"net/http"
"zhangyuhao/week03/practice/gin/models"
"github.com/gin-gonic/gin")
// CreateStudentHandler POST /students → 创建学生
// 接收前端传递的 JSON 数据，绑定到 Student 结构体，添加到学生列表中
// 返回创建成功的响应，包含新学生的信息
func CreateStudentHandler(c *gin.Context) { 
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"message": "请求参数错误",
			"error": err.Error()})
		return
	}
	models.StudentList=append(models.StudentList, student)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "创建学生成功",
		"count": len(models.StudentList),
		"data": student})

}
// GetAllStudents GET /students → 获取所有学生
// 返回所有学生的列表

func GetAllStudents(c *gin.Context) {
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "获取所有学生成功",
		"count": len(models.StudentList),
		"data": models.StudentList})

}
// GetStudent GET /students/:id → 根据 ID 获取学生
// 返回指定 ID 的学生信息，如果未找到则返回错误信息
func GetStudent(c *gin.Context) {
	idParam := c.Param("id")
	var student *models.Student
	for _, s := range models.StudentList {
		if fmt.Sprintf("%d", s.ID) == idParam {
			student = &s
			break
		}	
	}

	if student == nil { 
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "学生未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "获取学生成功",
		"data": student})
	}
// UpdateStudent PUT /students/:id , 简化版更新学生信息
// 接收前端传递的 JSON 数据，更新对应 ID 的学生信息
// 返回更新成功的响应，包含更新后的学生信息
func UpdateStudent(c *gin.Context) {
	// 1. 从 URL 拿 id 参数（字符串转整数）
	idStr := c.Param("id") // 读取路径参数 :id，比如 /students/1 → idStr = "1"
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 错误：id 不是数字（比如 /students/abc）
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "错误：学生 ID 必须是数字",
			"success": false,
		})
		return
	}

	// 2. 遍历切片找对应 id 的学生（记录索引，方便后续更新）
	findIndex := -1 // 初始值 -1 表示没找到
	for i, stu := range models.StudentList {
		if stu.ID == id {
			findIndex = i // 找到后记录索引位置
			break
		}
	}
	if findIndex == -1 {
		// 错误：没找到该 id 的学生
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("错误：未找到 ID 为 %d 的学生", id),
			"success": false,
		})
		return
	}

	// 3. 接收前端传递的更新后 JSON 数据（直接绑定到 Student 结构体）
	var updatedStu models.Student
	if err := c.ShouldBindJSON(&updatedStu); err != nil {
		// 错误：JSON 格式错误或字段不匹配
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "错误：更新数据格式不正确",
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	// 4. 执行更新（用新数据覆盖原数据，ID 保持不变，避免被篡改）
	updatedStu.ID = id // 强制保持原 ID，防止前端传错 ID
	models.StudentList[findIndex] = updatedStu // 用索引直接替换切片中的原学生

	// 5. 控制台打印更新后的信息（按你的要求）
	fmt.Printf("更新学生成功！ID：%d，更新后信息：%+v\n", id, updatedStu)

	// 6. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "学生信息更新成功",
		"success": true,
		"data":    updatedStu, // 返回更新后的完整信息
	})
}
func DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")
	index := -1
	for i, s := range models.StudentList {
		if fmt.Sprintf("%d", s.ID) == idParam {
			index = i
			break
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "学生未找到"})
		return
	}
	models.StudentList = append(models.StudentList[:index], models.StudentList[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "删除学生成功",
		"count": len(models.StudentList),
		"data": models.StudentList})
		
}