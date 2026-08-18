package main

import (
	"user_service/config"
	"user_service/etcd"
	"user_service/router"
)

func main() {
	// 初始化配置（包含数据库初始化）
	cfg := config.InitConfig()

	// 初始化etcd
	if err := etcd.InitEtcd(); err != nil {
		panic(err)
	}
	defer etcd.CloseEtcd()          // 关闭etcd
	if err := etcd.RegisterService( // user-service服务注册
		config.Cfg.Etcd.ServerName, config.Cfg.Etcd.ServeAddress); err != nil {
		panic(err)
	}

	// 初始化路由
	r := router.Router()

	// 启动服务
	r.Run(cfg.Server.Port)
}
