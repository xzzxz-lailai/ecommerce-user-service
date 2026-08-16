package main

import (
	"user_service/config"
	"user_service/router"
)

func main() {
	// 初始化配置（包含数据库初始化）
	cfg := config.InitConfig()

	// 初始化路由
	r := router.Router()

	// 启动服务
	r.Run(cfg.Server.Port)
}
