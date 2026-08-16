package config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func GetNacosConfig() (string, error) {

	// 配置nacos服务信息
	serverConfigs := []constant.ServerConfig{
		{
			//
			IpAddr: "127.0.0.1", // 服务器
			Port:   8848,
		},
	}

	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         "",   //当namespace是public时，此处填空字符串
		TimeoutMs:           5000, // 请求 Nacos 的超时时间（毫秒）
		NotLoadCacheAtStart: true, // 启动时是否不读取本地缓存
		//CacheDir:            "/opt/go-todo-backend/tmp/nacos/cache", // 服务器
		//LogDir:              "/opt/go-todo-backend/tmp/nacos/log",
		LogDir:   "./tmp/nacos/log", // 本地
		CacheDir: "./tmp/nacos/cache",
		LogLevel: "debug", // 日志级别
	}

	// 创建动态配置客户端的另一种方式
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err) // 如果失败,直接退出
	}

	// 从nacos拉取配置
	nacosConfig, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-service.yaml",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		panic(err) // 如果失败,直接退出
	}

	return nacosConfig, nil
}
