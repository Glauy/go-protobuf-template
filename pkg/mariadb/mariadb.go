package mariadb

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 数据库配置
type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	Params   string
	MaxOpen  int
	MaxIdle  int
	MaxLife  time.Duration
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Host:     "localhost",
	Port:     3306,
	User:     "root",
	Password: "password",
	DBName:   "testdb",
	Params:   "charset=utf8mb4&parseTime=True&loc=Local",
	MaxOpen:  50,
	MaxIdle:  10,
	MaxLife:  30 * time.Minute,
}

// New 使用配置结构体创建连接
func New(cfg Config) (*gorm.DB, error) {
	// 设置默认值
	if cfg.MaxOpen == 0 {
		cfg.MaxOpen = 50
	}
	if cfg.MaxIdle == 0 {
		cfg.MaxIdle = 10
	}
	if cfg.MaxLife == 0 {
		cfg.MaxLife = 30 * time.Minute
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Params)

	return createConnection(dsn, cfg.MaxOpen, cfg.MaxIdle, cfg.MaxLife)
}

// NewWithDefaults 使用默认配置创建连接
func NewWithDefaults() (*gorm.DB, error) {
	return New(DefaultConfig)
}

// NewFromDSN 从DSN字符串创建连接
func NewFromDSN(dsn string) (*gorm.DB, error) {
	return createConnection(dsn, 50, 10, 30*time.Minute)
}

// NewFromDSNWithPool 从DSN字符串创建连接，指定连接池参数
func NewFromDSNWithPool(dsn string, maxOpen, maxIdle int, maxLife time.Duration) (*gorm.DB, error) {
	return createConnection(dsn, maxOpen, maxIdle, maxLife)
}

// NewSimple 简单创建连接（最常用）
func NewSimple(host, user, password, dbName string) (*gorm.DB, error) {
	cfg := Config{
		Host:     host,
		User:     user,
		Password: password,
		DBName:   dbName,
		Params:   "charset=utf8mb4&parseTime=True&loc=Local",
	}
	return New(cfg)
}

// NewWithPort 指定端口创建连接
func NewWithPort(host string, port int, user, password, dbName string) (*gorm.DB, error) {
	cfg := Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		Params:   "charset=utf8mb4&parseTime=True&loc=Local",
	}
	return New(cfg)
}

// NewWithOptions 使用选项模式创建连接
func NewWithOptions(opts ...Option) (*gorm.DB, error) {
	cfg := DefaultConfig

	// 应用选项
	for _, opt := range opts {
		opt(&cfg)
	}

	return New(cfg)
}

// Option 配置选项函数
type Option func(*Config)

// WithHost 设置主机
func WithHost(host string) Option {
	return func(cfg *Config) {
		cfg.Host = host
	}
}

// WithPort 设置端口
func WithPort(port int) Option {
	return func(cfg *Config) {
		cfg.Port = port
	}
}

// WithUser 设置用户名
func WithUser(user string) Option {
	return func(cfg *Config) {
		cfg.User = user
	}
}

// WithPassword 设置密码
func WithPassword(password string) Option {
	return func(cfg *Config) {
		cfg.Password = password
	}
}

// WithDatabase 设置数据库名
func WithDatabase(dbName string) Option {
	return func(cfg *Config) {
		cfg.DBName = dbName
	}
}

// WithParams 设置连接参数
func WithParams(params string) Option {
	return func(cfg *Config) {
		cfg.Params = params
	}
}

// WithPool 设置连接池参数
func WithPool(maxOpen, maxIdle int, maxLife time.Duration) Option {
	return func(cfg *Config) {
		cfg.MaxOpen = maxOpen
		cfg.MaxIdle = maxIdle
		cfg.MaxLife = maxLife
	}
}

// WithLogger 设置日志级别
func WithLogger(level logger.LogLevel) Option {
	return func(cfg *Config) {
		// 这里可以扩展配置来支持日志级别
	}
}

// createConnection 创建数据库连接的核心函数
func createConnection(dsn string, maxOpen, maxIdle int, maxLife time.Duration) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLife)

	return db, nil
}

// MustNew 创建连接，失败时panic
func MustNew(cfg Config) *gorm.DB {
	db, err := New(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create database connection: %v", err))
	}
	return db
}

// MustNewSimple 简单创建连接，失败时panic
func MustNewSimple(host, user, password, dbName string) *gorm.DB {
	db, err := NewSimple(host, user, password, dbName)
	if err != nil {
		panic(fmt.Sprintf("failed to create database connection: %v", err))
	}
	return db
}
