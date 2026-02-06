package routes

import ("github.com/gin-gonic/gin"
"zhangyuhao/week03/practice/gin/controllers")

// SetupStudentRoutes 设置学生相关的路由
func SetupStudentRoutes(router *gin.Engine) { 
	studentRouter := router.Group("/student")
	{ 
		studentRouter.POST("/", controllers.CreateStudentHandler)
		studentRouter.GET("/", controllers.GetAllStudents)
		studentRouter.GET("/:id", controllers.GetStudent)
		studentRouter.PUT("/:id", controllers.UpdateStudent)
		studentRouter.DELETE("/:id", controllers.DeleteStudent)
	}

}