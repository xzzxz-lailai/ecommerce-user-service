package config

import (
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}
type MySQLConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"password"`
	DBName string `mapstructure:"dbname"`
}

type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

type EtcdConfig struct {
	Host         string `mapstructure:"host"` // etcd 服务地址
	ServerName   string `mapstructure:"server_name"`
	ServeAddress string `mapstructure:"address"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"` // JWT 密钥
	Expire int    `mapstructure:"expire"` // JWT 过期时间，单位小时
}

type Config struct {
	Server ServerConfig
	MySQL  MySQLConfig
	Etcd   EtcdConfig
	JWT    JWTConfig
}

var Cfg Config

func InitConfig() *Config {
	nacosConfig, err := GetNacosConfig() // 从 Nacos 拉配置（返回 yaml 字符串）
	if err != nil {
		panic(err)
	}

	// 配置文件类型（扩展名）
	viper.SetConfigType("yaml")

	// 读取配置内容
	err = viper.ReadConfig(strings.NewReader(nacosConfig))
	if err != nil {
		panic(err)
	}
	// 将读取到的配置解码到 Config 结构体
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(err)
	}

	if _, err := NewMySQL(); err != nil { // 初始化数据库
		panic(err)
	}

	return &Cfg
}
