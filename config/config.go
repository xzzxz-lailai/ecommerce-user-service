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
type JWTConfig struct {
	Secret string `mapstructure:"secret"` // 密钥
	Expire int    `mapstructure:"expire"` // 过期小时数  单位：小时
}
type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

type COSConfig struct {
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	Region    string `mapstructure:"region"`
	BucketURL string `mapstructure:"bucket_url"`
}
type EmailSmtpConfig struct {
	Host     string `mapstructure:"host"`     // 邮箱服务器
	Port     int    `mapstructure:"port"`     // 邮箱服务器端口
	Username string `mapstructure:"username"` // 发送人邮箱
	Password string `mapstructure:"pass"`     // SMTP 授权码
	From     string `mapstructure:"from"`     // 发送人邮箱
	FormName string `mapstructure:"from_name"`
}
type Config struct {
	Server    ServerConfig
	MySQL     MySQLConfig
	JWT       JWTConfig
	COS       COSConfig
	EmailSmtp EmailSmtpConfig
	CORS      CORSConfig
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

	// 初始化数据库
	NewMySQL()

	return &Cfg
}
