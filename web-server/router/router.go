package router

import (
	"github.com/gin-gonic/gin"

	"geekedu/web-server/handler"
	"geekedu/web-server/middleware"
)

// 初始化路由（满足Restful风格要求）
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 注册全局跨域中间件（JWT鉴权不全局使用，仅对需要权限的接口分组使用）
	r.Use(middleware.Cors())

	// API v1分组
	apiV1 := r.Group("/api/v1")
	{
		// 认证接口（无需JWT鉴权：登录、注册）
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", handler.LoginHandler) // 修正：对应handler中的方法名（原UserLogin改为LoginHandler）
			auth.POST("/register", handler.RegisterHandler) // 新增注册接口
		}

		// 需要JWT鉴权的接口分组
		authRequired := apiV1.Group("")
		authRequired.Use(middleware.JWTAuth())
		{
			// 课程接口
			courses := authRequired.Group("/courses")
			{
				courses.GET("", handler.GetCourseList)       // 获取课程列表
				courses.POST("", handler.CreateCourse)       // 发布课程（仅管理员）
				courses.POST("/upload/cover", handler.UploadCover) // 上传封面（仅管理员）
				courses.POST("/upload/video", handler.UploadVideo) // 上传视频（仅管理员）
				courses.GET("/:course_id/videos", handler.GetCourseVideos) // 新增：获取指定课程下所有视频（鉴权+验购买）
			}

			// 订单接口
			orders := authRequired.Group("/orders")
			{
				orders.POST("", handler.CreateOrder) // 购买课程（仅学生/登录用户）
			}

			// 播放接口
			player := authRequired.Group("/player")
			{
				player.GET("/:video_id", handler.GetVideoPlayUrl) // 播放视频（仅已购买用户）
			}
		}
	}

	return r
}