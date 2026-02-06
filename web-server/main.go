package main

import (
	"geekedu/common/config"
	"geekedu/web-server/handler"
	"geekedu/web-server/router"
	"strconv" // 已导入，无需额外添加
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()
	serverCfg := cfg.Server

	// 2. 初始化gRPC客户端（连接Logic Server）
	errGRPC := handler.InitGRPCClient(serverCfg.LogicAddr)
	if errGRPC != nil {
		panic("gRPC客户端初始化失败：" + errGRPC.Error())
	}

	// 3. 初始化路由
	r := router.InitRouter()

	// 4. 启动HTTP服务（修正：用strconv.Itoa()替换错误的itoa()）
	println("Web Server启动成功，监听端口：", serverCfg.WebPort)
	errRun := r.Run(":" + strconv.Itoa(serverCfg.WebPort))
	if errRun != nil {
		panic("Web Server启动失败：" + errRun.Error())
	}
}

// 关键：删除这个错误的itoa函数，整段移除！