package main

import (
	"net"
	"strconv" // 新增：导入正确的数字转字符串包
	_"geekedu/common/oss"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"geekedu/common/config"
	geekedu "geekedu/proto" // 修正：给proto包起别名geekedu，对应后续注册函数
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()
	dbCfg := cfg.DB
	serverCfg := cfg.Server

	// 2. 构建MySQL DSN（修正：使用strconv.Itoa()替换错误的itoa()）
	dsn := dbCfg.User + ":" + dbCfg.Password + "@tcp(" + dbCfg.Host + ":" + strconv.Itoa(dbCfg.Port) + ")/" + dbCfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// 3. 连接MySQL
	db, errDB := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 显示SQL日志，方便调试
	})
	if errDB != nil {
		panic("数据库连接失败：" + errDB.Error())
	}

	// 4. 自动迁移表结构（不存在则创建，满足快速启动需求）
	errMigrate := db.AutoMigrate(&User{}, &Course{}, &Video{}, &Order{})
	if errMigrate != nil {
		panic("表结构迁移失败：" + errMigrate.Error())
	}

	// 5. 初始化gRPC服务
	grpcServer := grpc.NewServer()

	// 6. 注册服务
	service := newGeekEduService(db)
	geekedu.RegisterGeekEduServiceServer(grpcServer, service) // 现在别名正常，可正确调用

	// 7. 监听端口（修正：使用strconv.Itoa()替换错误的itoa()）
	listen, errListen := net.Listen("tcp", ":"+strconv.Itoa(serverCfg.GRPCPort))
	if errListen != nil {
		panic("gRPC端口监听失败：" + errListen.Error())
	}

	// 8. 启动服务
	println("Logic Server启动成功，监听端口：", serverCfg.GRPCPort)
	if errServe := grpcServer.Serve(listen); errServe != nil {
		panic("gRPC服务启动失败：" + errServe.Error())
	}
	
}

