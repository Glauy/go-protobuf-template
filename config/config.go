package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Log      LogConfig      `mapstructure:"log"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	Environment  string `mapstructure:"environment"`
	Debug        bool   `mapstructure:"debug"`
	Port         int    `mapstructure:"port"` // HTTP端口
	Host         string `mapstructure:"host"` // HTTP主机
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
	MaxOpen  int    `mapstructure:"max_open"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxLife  string `mapstructure:"max_life"`
}

// GRPCConfig gRPC配置
type GRPCConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Reflection   bool   `mapstructure:"reflection"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	Filename string `mapstructure:"filename"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// App默认值
	viper.SetDefault("app.name", "go-protos")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)
	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.read_timeout", "30s")
	viper.SetDefault("app.write_timeout", "30s")

	// Database默认值
	viper.SetDefault("database.type", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.database", "testdb")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.max_open", 50)
	viper.SetDefault("database.max_idle", 10)
	viper.SetDefault("database.max_life", "30m")

	// GRPC默认值
	viper.SetDefault("grpc.host", "0.0.0.0")
	viper.SetDefault("grpc.port", 9090) // gRPC使用不同端口
	viper.SetDefault("grpc.reflection", true)
	viper.SetDefault("grpc.read_timeout", "30s")
	viper.SetDefault("grpc.write_timeout", "30s")

	// Log默认值
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.filename", "app.log")
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}

	if c.App.Port <= 0 || c.App.Port > 65535 {
		return fmt.Errorf("app port must be between 1 and 65535")
	}

	if c.GRPC.Port <= 0 || c.GRPC.Port > 65535 {
		return fmt.Errorf("grpc port must be between 1 and 65535")
	}

	// 确保HTTP和gRPC端口不同
	if c.App.Port == c.GRPC.Port {
		return fmt.Errorf("app port and grpc port cannot be the same")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
}

// GetAddress 获取应用服务地址
func (c *AppConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetAddress 获取gRPC服务地址
func (c *GRPCConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetMaxLifeDuration 获取最大连接生存时间
func (c *DatabaseConfig) GetMaxLifeDuration() (time.Duration, error) {
	return time.ParseDuration(c.MaxLife)
}

// GetReadTimeout 获取读取超时时间
func (c *AppConfig) GetReadTimeout() (time.Duration, error) {
	return time.ParseDuration(c.ReadTimeout)
}

// GetWriteTimeout 获取写入超时时间
func (c *AppConfig) GetWriteTimeout() (time.Duration, error) {
	return time.ParseDuration(c.WriteTimeout)
}
