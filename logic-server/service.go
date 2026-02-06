package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	errpkg "geekedu/common/err"
	"geekedu/common/jwt"
	"geekedu/common/oss"
	geekedu "geekedu/proto" // 统一使用geekedu别名
)

// 实现proto定义的GeekEduServiceServer接口
type geekEduService struct {
	geekedu.UnimplementedGeekEduServiceServer
	db *gorm.DB // 数据库连接
}

func newGeekEduService(db *gorm.DB) *geekEduService {
	return &geekEduService{db: db}
}

// 0. 学生注册（仅允许role=0，加详细日志）
func (s *geekEduService) Register(ctx context.Context, req *geekedu.RegisterRequest) (*geekedu.RegisterResponse, error) {
	fmt.Printf("[%s] [注册] 接收请求：username=%s\n", time.Now().Format("2006-01-02 15:04:05"), req.Username)
	if req.Username == "" || req.Password == "" {
		fmt.Printf("[%s] [注册错误] 参数为空：username=%s password=%s\n", time.Now().Format("2006-01-02 15:04:05"), req.Username, req.Password)
		return &geekedu.RegisterResponse{Code: 400, Msg: errpkg.ErrMsg[errpkg.ErrInvalidParam]}, nil
	}
	if len(req.Password) < 6 {
		fmt.Printf("[%s] [注册错误] 密码过短：username=%s 长度=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.Username, len(req.Password))
		return &geekedu.RegisterResponse{Code: 400, Msg: "密码长度不能少于6位"}, nil
	}

	// 密码加密
	hashedPwd, errPwd := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if errPwd != nil {
		fmt.Printf("[%s] [注册错误] 密码加密失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), errPwd)
		return &geekedu.RegisterResponse{Code: 500, Msg: errpkg.ErrMsg[errpkg.ErrInternalServer]}, nil
	}

	// 插入数据库
	var user User
	user.Username = req.Username
	user.Password = string(hashedPwd)
	user.Role = 0 // 强制学生
	result := s.db.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			fmt.Printf("[%s] [注册错误] 用户名已存在：%s\n", time.Now().Format("2006-01-02 15:04:05"), req.Username)
			return &geekedu.RegisterResponse{Code: 400, Msg: "用户名已存在"}, nil
		}
		fmt.Printf("[%s] [注册错误] 数据库插入失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), result.Error)
		return &geekedu.RegisterResponse{Code: 500, Msg: errpkg.ErrMsg[errpkg.ErrInternalServer]}, nil
	}

	fmt.Printf("[%s] [注册成功] 学生注册完成：username=%s uid=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.Username, user.ID)
	return &geekedu.RegisterResponse{Code: 200, Msg: "注册成功，请登录"}, nil
}

// 1. 用户登录（加详细日志，保留所有校验/Token生成逻辑）
func (s *geekEduService) UserLogin(ctx context.Context, req *geekedu.UserLoginRequest) (*geekedu.UserLoginResponse, error) {
	fmt.Printf("[%s] [登录] 接收请求：username=%s\n", time.Now().Format("2006-01-02 15:04:05"), req.Username)
	if req.Username == "" || req.Password == "" {
		fmt.Printf("[%s] [登录错误] 参数为空\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrMsg[errpkg.ErrInvalidParam])
	}

	// 查询用户
	var user User
	result := s.db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		fmt.Printf("[%s] [登录错误] 用户不存在：%s\n", time.Now().Format("2006-01-02 15:04:05"), req.Username)
		return nil, status.Error(codes.NotFound, errpkg.ErrMsg[errpkg.ErrUserNotExist])
	}

	// 验证密码
	if errPwd := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); errPwd != nil {
		fmt.Printf("[%s] [登录错误] 密码错误：%s\n", time.Now().Format("2006-01-02 15:04:05"), req.Username)
		return nil, status.Error(codes.PermissionDenied, errpkg.ErrMsg[errpkg.ErrPasswordWrong])
	}

	// 生成JWT
	token, errToken := jwt.GenerateToken(user.ID, user.Role)
	if errToken != nil {
		fmt.Printf("[%s] [登录错误] Token生成失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), errToken)
		return nil, status.Error(codes.Internal, errpkg.ErrMsg[errpkg.ErrInternalServer])
	}

	fmt.Printf("[%s] [登录成功] 用户登录：username=%s uid=%d role=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.Username, user.ID, user.Role)
	return &geekedu.UserLoginResponse{
		Uid:   uint32(user.ID),
		Token: token,
		Role:  int32(user.Role),
	}, nil
}

// 2. 获取课程列表（加详细日志）
func (s *geekEduService) GetCourseList(ctx context.Context, req *geekedu.CourseListRequest) (*geekedu.CourseListResponse, error) {
	fmt.Printf("[%s] [课程列表] 开始查询所有课程\n", time.Now().Format("2006-01-02 15:04:05"))
	var courses []Course
	result := s.db.Find(&courses)
	if result.Error != nil {
		fmt.Printf("[%s] [课程列表错误] 数据库查询失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), result.Error)
		return nil, status.Error(codes.Internal, errpkg.ErrMsg[errpkg.ErrInternalServer])
	}

	// 生成封面预签名URL
	var courseInfos []*geekedu.CourseInfo
	for _, c := range courses {
		coverUrl, _ := oss.GeneratePresignedURL(c.CoverOssKey)
		courseInfos = append(courseInfos, &geekedu.CourseInfo{
			Id:             uint32(c.ID),
			Title:          c.Title,
			Price:          c.Price,
			Intro:          c.Intro,
			CoverOssKey:    c.CoverOssKey,
			CoverSignedUrl: coverUrl,
			CreatedAt:      c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	fmt.Printf("[%s] [课程列表成功] 查询到%d门课程\n", time.Now().Format("2006-01-02 15:04:05"), len(courseInfos))
	return &geekedu.CourseListResponse{Courses: courseInfos}, nil
}

// 3. 发布课程（加详细日志）
func (s *geekEduService) CreateCourse(ctx context.Context, req *geekedu.CreateCourseRequest) (*geekedu.CreateCourseResponse, error) {
	fmt.Printf("[%s] [发布课程] 接收请求：uid=%d title=%s price=%.2f\n", time.Now().Format("2006-01-02 15:04:05"), req.Uid, req.Title, req.Price)
	if req.Uid == 0 || req.Title == "" || req.Price <= 0 || req.CoverOssKey == "" {
		fmt.Printf("[%s] [发布课程错误] 参数无效\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrMsg[errpkg.ErrInvalidParam])
	}

	// 校验管理员权限
	var user User
	if err := s.db.Where("id = ? AND role = 1", req.Uid).First(&user).Error; err != nil {
		fmt.Printf("[%s] [发布课程错误] 非管理员权限：uid=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.Uid)
		return nil, status.Error(codes.PermissionDenied, errpkg.ErrMsg[errpkg.ErrForbidden])
	}

	// 插入数据库
	course := Course{
		Title:        req.Title,
		Price:        req.Price,
		Intro:        req.Intro,
		CoverOssKey:  req.CoverOssKey,
		CreateUserID: uint(req.Uid),
	}
	if err := s.db.Create(&course).Error; err != nil {
		fmt.Printf("[%s] [发布课程错误] 插入失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return nil, status.Error(codes.Internal, errpkg.ErrMsg[errpkg.ErrInternalServer])
	}

	fmt.Printf("[%s] [发布课程成功] course_id=%d title=%s\n", time.Now().Format("2006-01-02 15:04:05"), course.ID, course.Title)
	return &geekedu.CreateCourseResponse{CourseId: uint32(course.ID)}, nil
}

// 4. 上传视频（加详细日志）
func (s *geekEduService) UploadVideo(ctx context.Context, req *geekedu.UploadVideoRequest) (*geekedu.UploadVideoResponse, error) {
	fmt.Printf("[%s] [上传视频] 接收请求：uid=%d course_id=%d title=%s\n", time.Now().Format("2006-01-02 15:04:05"), req.Uid, req.CourseId, req.Title)
	if req.Uid == 0 || req.CourseId == 0 || req.Title == "" || req.VideoOssKey == "" {
		fmt.Printf("[%s] [上传视频错误] 参数无效\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrMsg[errpkg.ErrInvalidParam])
	}

	// 校验管理员
	var user User
	if err := s.db.Where("id = ? AND role = 1", req.Uid).First(&user).Error; err != nil {
		fmt.Printf("[%s] [上传视频错误] 非管理员：uid=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.Uid)
		return nil, status.Error(codes.PermissionDenied, errpkg.ErrMsg[errpkg.ErrForbidden])
	}

	// 校验课程存在
	var course Course
	if err := s.db.Where("id = ?", req.CourseId).First(&course).Error; err != nil {
		fmt.Printf("[%s] [上传视频错误] 课程不存在：course_id=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.CourseId)
		return nil, status.Error(codes.NotFound, errpkg.ErrMsg[errpkg.ErrCourseNotExist])
	}

	// 插入视频
	video := Video{
		CourseID: uint(req.CourseId),
		Title:    req.Title,
		OssKey:   req.VideoOssKey, // 确保这里赋值正确，与数据库字段映射
	}
	if err := s.db.Create(&video).Error; err != nil {
		fmt.Printf("[%s] [上传视频错误] 插入失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return nil, status.Error(codes.Internal, errpkg.ErrMsg[errpkg.ErrInternalServer])
	}

	fmt.Printf("[%s] [上传视频成功] video_id=%d course_id=%d title=%s oss_key=%s\n", time.Now().Format("2006-01-02 15:04:05"), video.ID, req.CourseId, req.Title, req.VideoOssKey)
	return &geekedu.UploadVideoResponse{VideoId: uint32(video.ID)}, nil
}

// 5. 购买课程（加详细日志）
func (s *geekEduService) CreateOrder(ctx context.Context, req *geekedu.CreateOrderRequest) (*geekedu.CreateOrderResponse, error) {
	fmt.Printf("[%s] [购买课程] 接收请求：uid=%d course_id=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.Uid, req.CourseId)
	if req.Uid == 0 || req.CourseId == 0 {
		fmt.Printf("[%s] [购买课程错误] 参数无效\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrMsg[errpkg.ErrInvalidParam])
	}

	// 校验课程存在
	var course Course
	if err := s.db.Where("id = ?", req.CourseId).First(&course).Error; err != nil {
		fmt.Printf("[%s] [购买课程错误] 课程不存在：%d\n", time.Now().Format("2006-01-02 15:04:05"), req.CourseId)
		return nil, status.Error(codes.NotFound, errpkg.ErrMsg[errpkg.ErrCourseNotExist])
	}

	// 校验是否已购买
	var order Order
	if err := s.db.Where("user_id = ? AND course_id = ? AND order_status = 1", req.Uid, req.CourseId).First(&order).Error; err == nil {
		fmt.Printf("[%s] [购买课程错误] 已购买：uid=%d course_id=%d\n", time.Now().Format("2006-01-02 15:04:05"), req.Uid, req.CourseId)
		return nil, status.Error(codes.AlreadyExists, errpkg.ErrMsg[errpkg.ErrOrderExist])
	}

	// 创建订单
	order = Order{UserID: uint(req.Uid), CourseID: uint(req.CourseId), OrderStatus: 1}
	if err := s.db.Create(&order).Error; err != nil {
		fmt.Printf("[%s] [购买课程错误] 插入失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return nil, status.Error(codes.Internal, errpkg.ErrMsg[errpkg.ErrInternalServer])
	}

	fmt.Printf("[%s] [购买课程成功] order_id=%d uid=%d course_id=%d\n", time.Now().Format("2006-01-02 15:04:05"), order.ID, req.Uid, req.CourseId)
	return &geekedu.CreateOrderResponse{OrderId: uint32(order.ID)}, nil
}

// 6. 获取课程视频列表（已修复，返回真实数据）
func (s *geekEduService) GetCourseVideos(ctx context.Context, req *geekedu.GetCourseVideosRequest) (*geekedu.GetCourseVideosResponse, error) {
	// 补充日志格式（和其他方法保持一致）
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [获取课程视频] 接收请求：uid=%d course_id=%d\n", now, req.Uid, req.CourseId)

	// 1. 参数校验
	if req.Uid == 0 || req.CourseId == 0 {
		fmt.Printf("[%s] [获取课程视频错误] 参数无效：uid=%d course_id=%d\n", now, req.Uid, req.CourseId)
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrMsg[errpkg.ErrInvalidParam])
	}

	// 2. 校验课程是否存在
	var course Course
	if err := s.db.Where("id = ?", req.CourseId).First(&course).Error; err != nil {
		fmt.Printf("[%s] [获取课程视频错误] 课程不存在：course_id=%d\n", now, req.CourseId)
		return nil, status.Error(codes.NotFound, errpkg.ErrMsg[errpkg.ErrCourseNotExist])
	}

	// 3. 权限校验：管理员直接查询，普通用户需已购买
	var user User
	if err := s.db.Where("id = ?", req.Uid).First(&user).Error; err != nil {
		fmt.Printf("[%s] [获取课程视频错误] 用户不存在：uid=%d\n", now, req.Uid)
		return nil, status.Error(codes.NotFound, errpkg.ErrMsg[errpkg.ErrUserNotExist])
	}
	if user.Role != 1 { // 非管理员
		var order Order
		if err := s.db.Where("user_id = ? AND course_id = ? AND order_status = 1", req.Uid, req.CourseId).First(&order).Error; err != nil {
			fmt.Printf("[%s] [获取课程视频错误] 未购买课程：uid=%d course_id=%d\n", now, req.Uid, req.CourseId)
			return nil, status.Error(codes.PermissionDenied, errpkg.ErrMsg[errpkg.ErrNotPaid])
		}
	}

	// 4. 查询该课程下所有视频（按创建时间正序排列）
	var videos []Video
	result := s.db.Where("course_id = ?", req.CourseId).Order("created_at ASC").Find(&videos)
	if result.Error != nil {
		fmt.Printf("[%s] [获取课程视频错误] 数据库查询失败：%v\n", now, result.Error)
		return nil, status.Error(codes.Internal, errpkg.ErrMsg[errpkg.ErrInternalServer])
	}

	// 5. 转换为proto响应格式（填充真实视频数据）
	var videoSimpleInfos []*geekedu.VideoSimpleInfo
	for _, v := range videos {
		videoSimpleInfos = append(videoSimpleInfos, &geekedu.VideoSimpleInfo{
			VideoId:   uint32(v.ID),       // 数据库Video ID转uint32
			Title:     v.Title,            // 视频标题
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"), // 时间格式化
		})
	}

	// 6. 返回真实数据（非空列表）
	fmt.Printf("[%s] [获取课程视频成功] 查询到%d个视频：course_id=%d\n", now, len(videoSimpleInfos), req.CourseId)
	return &geekedu.GetCourseVideosResponse{
		Videos: videoSimpleInfos,
	}, nil
}

// 7. 获取视频播放URL（核心修复：增加OssKey校验+完整日志+确保返回有效URL）
func (s *geekEduService) GetVideoPlayUrl(ctx context.Context, req *geekedu.GetVideoPlayUrlRequest) (*geekedu.GetVideoPlayUrlResponse, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [获取播放URL] 接收请求：uid=%d video_id=%d\n", now, req.Uid, req.VideoId)

	// 1. 严格参数校验
	if req.Uid == 0 || req.VideoId == 0 {
		fmt.Printf("[%s] [获取播放URL错误] 参数无效：uid=%d video_id=%d\n", now, req.Uid, req.VideoId)
		return nil, status.Error(codes.InvalidArgument, errpkg.ErrMsg[errpkg.ErrInvalidParam])
	}

	// 2. 查询视频信息（确保OssKey不为空）
	var video Video
	result := s.db.Where("id = ?", req.VideoId).First(&video)
	if result.Error != nil {
		fmt.Printf("[%s] [获取播放URL错误] 视频不存在：video_id=%d\n", now, req.VideoId)
		return nil, status.Error(codes.NotFound, errpkg.ErrMsg[errpkg.ErrVideoNotExist])
	}
	// 关键修复：校验OssKey非空，避免生成无效URL
	if video.OssKey == "" {
		fmt.Printf("[%s] [获取播放URL错误] 视频OSS Key为空：video_id=%d\n", now, req.VideoId)
		return nil, status.Error(codes.InvalidArgument, "视频存储地址无效，无法播放")
	}

	// 3. 查询用户信息
	var user User
	result = s.db.Where("id = ?", req.Uid).First(&user)
	if result.Error != nil {
		fmt.Printf("[%s] [获取播放URL错误] 用户不存在：uid=%d\n", now, req.Uid)
		return nil, status.Error(codes.NotFound, errpkg.ErrMsg[errpkg.ErrUserNotExist])
	}

	// 4. 权限校验：管理员直接播放，普通用户需已购买
	if user.Role != 1 {
		var order Order
		result = s.db.Where("user_id = ? AND course_id = ? AND order_status = 1", req.Uid, video.CourseID).First(&order)
		if result.Error != nil {
			fmt.Printf("[%s] [获取播放URL错误] 未购买课程：uid=%d course_id=%d\n", now, req.Uid, video.CourseID)
			return nil, status.Error(codes.PermissionDenied, errpkg.ErrMsg[errpkg.ErrNotPaid])
		}
	}

	// 5. 生成有效OSS预签名URL（确保3600秒过期，日志记录URL生成状态）
	signedUrl, err := oss.GeneratePresignedURL(video.OssKey)
	if err != nil {
		fmt.Printf("[%s] [获取播放URL错误] 生成签名URL失败：video_id=%d oss_key=%s error=%v\n", now, req.VideoId, video.OssKey, err)
		return nil, status.Error(codes.Internal, "生成视频播放地址失败，请重试")
	}
	if signedUrl == "" {
		fmt.Printf("[%s] [获取播放URL错误] 签名URL为空：video_id=%d oss_key=%s\n", now, req.VideoId, video.OssKey)
		return nil, status.Error(codes.Internal, "视频播放地址无效，无法播放")
	}

	// 6. 日志记录成功状态，返回有效URL
	fmt.Printf("[%s] [获取播放URL成功] video_id=%d oss_key=%s 签名URL生成完成（有效期1小时）\n", now, req.VideoId, video.OssKey)
	return &geekedu.GetVideoPlayUrlResponse{SignedUrl: signedUrl}, nil
}

// 辅助函数：生成JWT（保留，无额外修改）
func generateToken(uid uint, role int) (string, error) {
	return jwt.GenerateToken(uid, role)
}