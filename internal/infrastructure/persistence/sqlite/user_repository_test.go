package sqlite

import (
	"os/user"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}
	// 自动迁移
	if err := db.AutoMigrate(&user.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

// func TestUserRepository_CRUD(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := NewUserRepository(db)

// 	// 1. 新增用户
// 	u := &user.User{Username: "alice", Email: "alice@test.com", Phone: "123456789"}
// 	err := repo.Save(u)
// 	assert.NoError(t, err)
// 	assert.NotZero(t, u.ID)

// 	// 2. 按ID查找
// 	got, err := repo.FindByID(u.ID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "alice", got.Username)

// 	// 3. 按用户名查找
// 	got2, err := repo.FindByUsername("alice")
// 	assert.NoError(t, err)
// 	assert.Equal(t, got.ID, got2.ID)

// 	// 4. 更新状态
// 	err = repo.UpdateStatus(u.ID, 3)
// 	assert.NoError(t, err)

// 	updated, err := repo.FindByID(u.ID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int8(3), updated.Status)
// }
