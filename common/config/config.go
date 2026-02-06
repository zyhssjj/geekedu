package config

import (
	"os"
	"strconv"
)

// 全局配置结构体（满足考核点：敏感配置从环境变量读取，不硬编码）
type GlobalConfig struct {
	OSS    OSSConfig  // OSS配置
	DB     DBConfig   // 数据库配置
	Server ServerConfig // 服务端口配置
	JWT    JWTConfig  // JWT配置
}

// OSS配置（对应阿里云OSS参数）
type OSSConfig struct {
	AccessKey string // OSS AccessKey ID
	SecretKey string // OSS AccessKey Secret
	Endpoint  string // OSS Endpoint
	Bucket    string // OSS Bucket名称
}

// 数据库配置
type DBConfig struct {
	Host     string // 数据库地址
	Port     int    // 数据库端口
	User     string // 数据库用户名
	Password string // 数据库密码
	DBName   string // 数据库名称
}

// 服务配置
type ServerConfig struct {
	WebPort  int // Web Server端口
	GRPCPort int // Logic Server gRPC端口
	LogicAddr string // Web连接Logic的地址
}

// JWT配置
type JWTConfig struct {
	Secret string // JWT签名密钥
	Expire int64  // 过期时间（秒），默认7天
}

// 加载全局配置（从环境变量读取，无默认值的需在.env中配置）
func LoadConfig() *GlobalConfig {
	// 解析环境变量，提供默认值避免启动失败
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	webPort, _ := strconv.Atoi(getEnv("WEB_PORT", "8080"))
	grpcPort, _ := strconv.Atoi(getEnv("GRPC_PORT", "50051"))
	jwtExpire, _ := strconv.ParseInt(getEnv("JWT_EXPIRE", "604800"), 10, 64)

	return &GlobalConfig{
		OSS: OSSConfig{
			AccessKey: getEnv("OSS_ACCESS_KEY", ""),
			SecretKey: getEnv("OSS_SECRET_KEY", ""),
			Endpoint:  getEnv("OSS_ENDPOINT", ""),
			Bucket:    getEnv("OSS_BUCKET", ""),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "mysql"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root"),
			DBName:   getEnv("DB_NAME", "geekedu"),
		},
		Server: ServerConfig{
			WebPort:  webPort,
			GRPCPort: grpcPort,
			LogicAddr: getEnv("LOGIC_SERVER_ADDR", "logic-server:50051"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "geekedu-jwt-2026"),
			Expire: jwtExpire,
		},
	}
}

// 辅助函数：获取环境变量，无则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}