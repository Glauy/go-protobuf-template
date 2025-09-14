package database

import (
	"fmt"
	"log"

	"go-protos/internal/domain"

	"gorm.io/gorm"
)

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// 迁移所有表
	if err := db.AutoMigrate(
		&domain.User{},
		// 在这里添加其他实体
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// MigrateWithLog 带详细日志的迁移
func MigrateWithLog(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// 检查表是否存在
	var tables []string
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to check existing tables: %w", err)
	}

	log.Printf("Existing tables: %v", tables)

	// 执行迁移
	if err := Migrate(db); err != nil {
		return err
	}

	// 再次检查表
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to check tables after migration: %w", err)
	}

	log.Printf("Tables after migration: %v", tables)
	return nil
}
