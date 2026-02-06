package err

// 全局错误码定义（满足考核点：统一错误处理，不直接panic）
const (
	Success = 200 // 操作成功

	ErrInvalidParam   = 400 // 参数无效
	ErrUnauthorized   = 401 // 未授权（JWT无效/未登录）
	ErrForbidden      = 403 // 禁止访问（无权限）
	ErrNotFound       = 404 // 资源不存在
	ErrInternalServer = 500 // 服务器内部错误
	ErrInvalidToken  = 40101 // Token格式无效
	ErrTokenExpired  = 40102 // Token已过期
	ErrUserNotExist    = 1001 // 用户不存在
	ErrPasswordWrong   = 1002 // 密码错误
	ErrCourseNotExist  = 2001 // 课程不存在
	ErrVideoNotExist   = 2002 // 视频不存在
	ErrNotPaid         = 2003 // 未购买该课程
	ErrOrderExist      = 2004 // 已购买该课程
	ErrOSSUploadFailed = 3001 // OSS上传失败
	ErrOSSGenerateURL  = 3002 // OSS预签名URL生成失败
)

// 错误信息映射
var ErrMsg = map[int]string{
	Success:           "操作成功",
	ErrInvalidParam:   "参数无效",
	ErrUnauthorized:   "未授权，请先登录",
	ErrForbidden:      "无权限访问该资源",
	ErrNotFound:       "请求的资源不存在",
	ErrInternalServer: "服务器内部错误",

	ErrUserNotExist:    "用户不存在",
	ErrPasswordWrong:   "密码错误",
	ErrCourseNotExist:  "课程不存在",
	ErrVideoNotExist:   "视频不存在",
	ErrNotPaid:         "请先购买该课程再播放",
	ErrOrderExist:      "已购买该课程，无需重复购买",
	ErrOSSUploadFailed: "文件上传失败",
	ErrOSSGenerateURL:  "播放地址生成失败",
	ErrInvalidToken:  "Token无效，请重新登录",
	ErrTokenExpired:  "Token已过期，请重新登录",
}

// 统一响应结构体（前端可直接解析）
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 构建成功响应
func SuccessResponse(data interface{}) *Response {
	return &Response{
		Code: Success,
		Msg:  ErrMsg[Success],
		Data: data,
	}
}

// 构建错误响应
func ErrorResponse(code int) *Response {
	return &Response{
		Code: code,
		Msg:  ErrMsg[code],
		Data: nil,
	}
}