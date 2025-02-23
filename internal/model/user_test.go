package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用SQLite内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移表结构
	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)

	// 创建测试用户
	user := &User{
		Username:  "testuser",
		Password:  "testpass",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 测试创建用户
	result := db.Create(user)
	assert.NoError(t, result.Error)
	assert.NotZero(t, user.ID)

	// 验证用户是否被正确创建
	var found User
	result = db.First(&found, user.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
}

func TestFindUserByUsername(t *testing.T) {
	db := setupTestDB(t)

	// 创建测试用户
	user := &User{
		Username:  "finduser",
		Password:  "testpass",
		Email:     "find@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(user)

	// 测试通过用户名查找
	var found User
	result := db.Where("username = ?", user.Username).First(&found)
	assert.NoError(t, result.Error)
	assert.Equal(t, user.Username, found.Username)

	// 测试查找不存在的用户
	result = db.Where("username = ?", "nonexistent").First(&found)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)

	// 创建测试用户
	user := &User{
		Username:  "updateuser",
		Password:  "oldpass",
		Email:     "old@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(user)

	// 更新用户信息
	newEmail := "new@example.com"
	result := db.Model(user).Update("email", newEmail)
	assert.NoError(t, result.Error)

	// 验证更新是否成功
	var updated User
	db.First(&updated, user.ID)
	assert.Equal(t, newEmail, updated.Email)
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)

	// 创建测试用户
	user := &User{
		Username:  "deleteuser",
		Password:  "testpass",
		Email:     "delete@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(user)

	// 删除用户
	result := db.Delete(user)
	assert.NoError(t, result.Error)

	// 验证用户是否被删除
	var found User
	result = db.First(&found, user.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}