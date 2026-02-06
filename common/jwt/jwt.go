package jwt

import (
	"time"

	"geekedu/common/config"
	"geekedu/common/err"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义Claims（包含用户ID和角色，满足鉴权需求）
type CustomClaims struct {
	UID  uint `json:"uid"`  // 用户ID
	Role int  `json:"role"` // 0-学员，1-管理员
	jwt.RegisteredClaims
}

// 生成JWT Token（登录成功后返回给前端）
func GenerateToken(uid uint, role int) (string, error) {
	cfg := config.LoadConfig()

	// 兜底：防止JWT Secret为空（返回error类型的内置常量）
	if cfg.JWT.Secret == "" {
		return "", jwt.ErrSignatureInvalid
	}

	// 构建Claims
	claims := CustomClaims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(cfg.JWT.Expire))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "geekedu",
		},
	}

	// 生成Token（使用HS256加密）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// 验证JWT Token（仅去掉jwt.ErrTokenInvalid，解决未定义报错）
func VerifyToken(tokenString string) (*CustomClaims, *err.Response) {
	cfg := config.LoadConfig()

	// 兜底1：防止JWT Secret为空（返回统一错误响应）
	if cfg.JWT.Secret == "" {
		return nil, err.ErrorResponse(err.ErrInternalServer)
	}

	// 解析Token（适配jwt/v5，补充算法校验，不依赖ErrTokenInvalid）
	token, errParse := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 关键：校验签名算法是否为HS256（解决算法篡改问题）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWT.Secret), nil
	})

	// 解析失败（仅保留存在的常量，其他错误统一返回无效Token，解决报错）
	if errParse != nil {
		switch errParse {
		// 仅保留jwt/v5 一定存在的常量：Token已过期
		case jwt.ErrTokenExpired:
			return nil, err.ErrorResponse(err.ErrTokenExpired)
		// 签名无效（所有版本都有），返回Token无效
		case jwt.ErrSignatureInvalid:
			return nil, err.ErrorResponse(err.ErrInvalidToken)
		// 其他所有未知错误，统一返回未授权（兜底，不依赖任何不确定常量）
		default:
			return nil, err.ErrorResponse(err.ErrUnauthorized)
		}
	}

	// 验证Token有效性+类型断言（保持严谨）
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, err.ErrorResponse(err.ErrInvalidToken)
	}

	// 验证通过，返回自定义Claims
	return claims, nil
}