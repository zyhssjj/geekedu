package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"geekedu/common/jwt"
)

// Cors：跨域中间件（前端可正常请求）
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 预处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JWTAuth：JWT鉴权中间件（仅挂载在需要鉴权的路由分组上）
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "缺少token",
			})
			return
		}

		// 校验Token格式（Bearer + Token）
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token格式错误",
			})
			return
		}

		// 验证Token有效性
		claims, respErr := jwt.VerifyToken(parts[1])
		if respErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  respErr.Msg,
			})
			return
		}

		// 保存用户信息到Gin上下文
		c.Set("uid", claims.UID)
		c.Set("role", claims.Role)

		c.Next()
	}
}