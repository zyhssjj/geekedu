package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"geekedu/common/err"
	"geekedu/common/oss"
	geekedu "geekedu/proto" // 修复：使用模块路径并设置别名
)

// 注意：这里修正 proto 包的引用，确保和你的项目一致（原代码是 geekedu，统一对齐）
var geekEduClient geekedu.GeekEduServiceClient

// 初始化gRPC客户端（连接Logic Server）
func InitGRPCClient(addr string) error {
	conn, errConn := grpc.Dial(addr, grpc.WithInsecure()) // 开发环境无TLS，生产环境需配置
	if errConn != nil {
		return errConn
	}

	geekEduClient = geekedu.NewGeekEduServiceClient(conn)
	return nil
}

// 新增：0. 学生注册接口（无需鉴权）- 核心修复：统一返回数据结构，data返回空对象而非nil
func RegisterHandler(c *gin.Context) {
	var req geekedu.RegisterRequest
	// 绑定前端JSON参数
	if errBind := c.ShouldBindJSON(&req); errBind != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// 调用gRPC注册服务
	resp, errGrpc := geekEduClient.Register(c, &req)
	if errGrpc != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	// 核心修复：统一数据结构，data返回空对象（{}）而非nil，避免前端取值undefined
	c.JSON(200, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
		"data": gin.H{}, // 修复：返回空对象，而非nil，保持和登录接口数据结构一致
	})
}

// 1. 用户登录 - 核心修复：补充返回role字段，规范数据结构
func LoginHandler(c *gin.Context) {
	var req geekedu.UserLoginRequest
	if errBind := c.ShouldBindJSON(&req); errBind != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	resp, errGrpc := geekEduClient.UserLogin(c, &req)
	if errGrpc != nil {
		// 解析gRPC错误码，返回对应前端提示
		st, ok := status.FromError(errGrpc)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(200, err.ErrorResponse(err.ErrUserNotExist))
			case codes.PermissionDenied:
				c.JSON(200, err.ErrorResponse(err.ErrPasswordWrong))
			default:
				c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
			}
			return
		}
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	// 核心修复：返回完整数据（uid、token、role），保持数据结构统一
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"uid":   resp.Uid,
			"token": resp.Token,
			"role":  resp.Role,
		},
	})
}

// 2. 获取课程列表 - 保持不变，保留封面预签名URL返回
func GetCourseList(c *gin.Context) {
	resp, errGrpc := geekEduClient.GetCourseList(c, &geekedu.CourseListRequest{})
	if errGrpc != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	c.JSON(200, err.SuccessResponse(resp))
}

// 3. 上传封面（先上传到OSS，返回OSS Key）- 保持不变，保留分片上传逻辑
func UploadCover(c *gin.Context) {
	file, fileHeader, errFile := c.Request.FormFile("cover")
	if errFile != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// 读取文件字节流
	fileData, errRead := io.ReadAll(file)
	if errRead != nil {
		c.JSON(200, err.ErrorResponse(err.ErrOSSUploadFailed))
		return
	}

	// 上传到OSS
	ossKey, respErr := oss.SimpleUpload(fileHeader.Filename, fileData)
	if respErr != nil {
		c.JSON(200, respErr)
		return
	}

	c.JSON(200, err.SuccessResponse(map[string]string{"cover_oss_key": ossKey}))
}

// 4. 发布课程（核心修复：uid nil 校验、安全类型转换，避免 panic）- 保持不变
func CreateCourse(c *gin.Context) {
	var req struct {
		Title       string  `json:"title"`
		Price       float64 `json:"price"`
		Intro       string  `json:"intro"`
		CoverOssKey string  `json:"cover_oss_key"` // 对齐前端传递的字段名，确保能正确绑定
	}
	if errBind := c.ShouldBindJSON(&req); errBind != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// ------------- 修复核心：uid 安全获取与转换 -------------
	// 1. 先获取 uid，判断是否存在且非 nil
	uidVal, exists := c.Get("uid")
	if !exists || uidVal == nil {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}

	// 2. 安全转换为 uint 类型，避免强制转换 panic
	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	// --------------------------------------------------------

	// 构建gRPC请求（转换为 uint32 适配 proto 定义）
	grpcReq := &geekedu.CreateCourseRequest{
		Uid:         uint32(uid),
		Title:       req.Title,
		Price:       req.Price,
		Intro:       req.Intro,
		CoverOssKey: req.CoverOssKey,
	}

	resp, errGrpc := geekEduClient.CreateCourse(c, grpcReq)
	if errGrpc != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	c.JSON(200, err.SuccessResponse(resp))
}

// 5. 上传视频（分片上传到OSS，创建视频记录：同步修复 uid 安全校验）- 保持不变
func UploadVideo(c *gin.Context) {
	file, fileHeader, errFile := c.Request.FormFile("video")
	if errFile != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// 获取表单参数
	courseId, _ := strconv.Atoi(c.PostForm("course_id"))
	title := c.PostForm("title")
	if courseId <= 0 || title == "" {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// ------------- 同步修复：uid 安全获取与转换 -------------
	uidVal, exists := c.Get("uid")
	if !exists || uidVal == nil {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	// --------------------------------------------------------

	// 分片上传到OSS
	ossKey, respErr := oss.MultipartUpload(fileHeader.Filename, file)
	if respErr != nil {
		c.JSON(200, respErr)
		return
	}

	// 构建gRPC请求
	grpcReq := &geekedu.UploadVideoRequest{
		Uid:         uint32(uid),
		CourseId:    uint32(courseId),
		Title:       title,
		VideoOssKey: ossKey,
	}

	resp, errGrpc := geekEduClient.UploadVideo(c, grpcReq)
	if errGrpc != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	c.JSON(200, err.SuccessResponse(resp))
}

// 6. 购买课程（同步修复：uid 安全校验）- 保持不变
func CreateOrder(c *gin.Context) {
	var req struct {
		CourseId uint32 `json:"course_id"`
	}
	if errBind := c.ShouldBindJSON(&req); errBind != nil || req.CourseId == 0 {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// ------------- 同步修复：uid 安全获取与转换 -------------
	uidVal, exists := c.Get("uid")
	if !exists || uidVal == nil {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	// --------------------------------------------------------

	// 构建gRPC请求
	grpcReq := &geekedu.CreateOrderRequest{
		Uid:      uint32(uid),
		CourseId: req.CourseId,
	}

	resp, errGrpc := geekEduClient.CreateOrder(c, grpcReq)
	if errGrpc != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	c.JSON(200, err.SuccessResponse(resp))
}

// 7. 获取视频播放URL（核心修复：解析详细gRPC错误，返回明确提示，确保前端能接收）
func GetVideoPlayUrl(c *gin.Context) {
	// 获取URL参数
	videoId, errParse := strconv.Atoi(c.Param("video_id"))
	if errParse != nil || videoId <= 0 {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// ------------- 同步修复：uid 安全获取与转换 -------------
	uidVal, exists := c.Get("uid")
	if !exists || uidVal == nil {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	// --------------------------------------------------------

	// 构建gRPC请求
	grpcReq := &geekedu.GetVideoPlayUrlRequest{
		Uid:     uint32(uid),
		VideoId: uint32(videoId),
	}

	// 调用gRPC服务，解析详细错误码（关键：让前端收到明确提示，不再无反应）
	resp, errGrpc := geekEduClient.GetVideoPlayUrl(c, grpcReq)
	if errGrpc != nil {
		st, ok := status.FromError(errGrpc)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.JSON(200, gin.H{"code": 400, "msg": st.Message(), "data": gin.H{}})
			case codes.NotFound:
				c.JSON(200, gin.H{"code": 404, "msg": st.Message(), "data": gin.H{}})
			case codes.PermissionDenied:
				c.JSON(200, gin.H{"code": 403, "msg": st.Message(), "data": gin.H{}})
			default:
				c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
			}
			return
		}
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	// 核心：返回规范数据结构，前端可直接获取 signed_url 进行播放
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取播放地址成功",
		"data": gin.H{
			"signed_url": resp.SignedUrl, // 明确返回播放URL，前端可直接使用
		},
	})
}

// 补充：获取课程下所有视频接口（对应gRPC新增方法，保持uid安全校验）
func GetCourseVideos(c *gin.Context) {
	// 获取URL参数
	courseId, errParse := strconv.Atoi(c.Param("course_id"))
	if errParse != nil || courseId <= 0 {
		c.JSON(200, err.ErrorResponse(err.ErrInvalidParam))
		return
	}

	// uid 安全获取与转换
	uidVal, exists := c.Get("uid")
	if !exists || uidVal == nil {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}
	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, err.ErrorResponse(err.ErrUnauthorized))
		return
	}

	// 构建gRPC请求
	grpcReq := &geekedu.GetCourseVideosRequest{
		Uid:      uint32(uid),
		CourseId: uint32(courseId),
	}

	resp, errGrpc := geekEduClient.GetCourseVideos(c, grpcReq)
	if errGrpc != nil {
		c.JSON(200, err.ErrorResponse(err.ErrInternalServer))
		return
	}

	c.JSON(200, err.SuccessResponse(resp))
}