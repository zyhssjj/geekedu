package main
import(
	"fmt"
	"zhangyuhao/week03/practice/gin/routes"
	"github.com/gin-gonic/gin"


)
func main() { 
	
	r := gin.Default()
	// Setup student routes
	routes.SetupStudentRoutes(r)
	fmt.Println("Server started")
	if err := r.Run(":8080"); err != nil {
		panic("failed to start server"+err.Error())
	}
}
