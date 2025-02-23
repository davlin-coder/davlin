package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type LLMConfig struct {
	Model   string `mapstructure:"model"`
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type APPConfig struct {
	Port int `mapstructure:"port"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// SMTPConfig 定义SMTP服务器配置
type SMTPConfig struct {
	Host     string `yaml:"host"`     // SMTP服务器地址
	Port     int    `yaml:"port"`     // SMTP服务器端口
	Username string `yaml:"username"` // SMTP认证用户名
	Password string `yaml:"password"` // SMTP认证密码
	From     string `yaml:"from"`     // 发件人邮箱地址
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey string        `mapstructure:"secret_key"` // JWT密钥
	Expire    time.Duration `mapstructure:"expire"`     // 过期时间（小时）
}

type Config struct {
	LLM   LLMConfig   `mapstructure:"llm"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	APP   APPConfig   `mapstructure:"app"`
	Redis RedisConfig `mapstructure:"redis"`
	Email SMTPConfig  `mapstructure:"email"`
	JWT   JWTConfig   `mapstructure:"jwt"`
}

var cfg *Config

// Init 初始化配置
func Init() (*Config, error) {
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")      // 查找配置文件所在的路径

	// 设置环境变量前缀，用于区分不同应用的配置
	viper.SetEnvPrefix("DAVLIN")

	viper.SetDefault("app.port", 8080)

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件，如果文件不存在则跳过
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 如果错误不是因为找不到配置文件，则返回错误
			fmt.Printf("读取配置文件失败: %s\n", err)
			return nil, err
		}
		// 如果是因为找不到配置文件，则继续使用环境变量
		fmt.Println("未找到配置文件，将使用环境变量作为配置来源")
	}

	fmt.Println(viper.AllKeys())

	// 将配置解析到struct中
	cfg = &Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		fmt.Printf("配置解析失败: %s\n", err)
		return nil, err
	}

	// 监控配置文件变化并热加载程序
	viper.WatchConfig()

	return cfg, nil
}
