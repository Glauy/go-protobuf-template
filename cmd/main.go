package main

import (
	"flag"
	"log"

	"go-protos/config"
	"go-protos/internal/application"
	"go-protos/internal/domain"
	"go-protos/internal/infrastructure/database"
	persistence "go-protos/internal/infrastructure/persistence/mariadb"
	"go-protos/internal/interfaces/grpc"
	"go-protos/pkg/mariadb"
)

func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	log.Printf("Starting %s v%s in %s mode", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// 初始化数据库连接
	maxLife, err := cfg.Database.GetMaxLifeDuration()
	if err != nil {
		log.Fatal("Failed to parse max life duration:", err)
	}

	db, err := mariadb.NewFromDSNWithPool(
		cfg.Database.GetDSN(),
		cfg.Database.MaxOpen,
		cfg.Database.MaxIdle,
		maxLife,
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 执行数据库迁移
	if err := database.MigrateWithLog(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化仓储
	userRepo := persistence.NewUserRepository(db)

	// 初始化领域服务
	userDomainSvc := domain.NewUserDomainService(userRepo)

	// 初始化应用服务
	userAppSvc := application.NewUserAppService(userRepo, userDomainSvc)

	// 初始化gRPC服务器
	grpcServer := grpc.NewServer(userAppSvc)

	// 启动gRPC服务器
	grpcAddress := cfg.GRPC.GetAddress()
	log.Printf("gRPC server starting on %s", grpcAddress)
	if err := grpcServer.Start(grpcAddress); err != nil {
		log.Fatal("Failed to start gRPC server:", err)
	}
}
