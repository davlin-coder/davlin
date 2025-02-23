package service

import (
	"context"
	"errors"
	"testing"

	"github.com/davlin-coder/davlin/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 创建一个mock的LLM客户端
type MockLLMClient struct {
	mock.Mock
}

func (m *MockLLMClient) Chat(ctx context.Context, messages []*model.ChatMessage) (string, error) {
	args := m.Called(ctx, messages)
	return args.String(0), args.Error(1)
}

func TestChatService(t *testing.T) {
	// 创建聊天服务实例
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移数据库表结构
	err = db.AutoMigrate(&model.ChatMessage{})
	assert.NoError(t, err)

	chatService := NewChatService(db)

	// 测试发送消息
	message := &model.ChatMessage{
		Role:    "user",
		Content: "Hello",
	}

	// 执行测试
	response, err := chatService.SendMessage(message)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, map[string]interface{}{}, response)
}

func TestChatServiceError(t *testing.T) {
	// 创建一个无效的数据库连接来模拟错误
	_, err := gorm.Open(sqlite.Open("/invalid/path"), &gorm.Config{})
	if err != nil {
		// 如果数据库连接失败，创建一个新的服务实例
		chatService := NewChatService(nil)

		// 测试发送消息
		message := &model.ChatMessage{
			Role:    "user",
			Content: "Hello",
		}

		// 执行测试
		response, err := chatService.SendMessage(message)

		// 验证错误处理
		assert.Error(t, err)
		assert.Equal(t, errors.New("database connection is not initialized"), err)
		assert.Nil(t, response)
		return
	}

	// 如果数据库连接成功（不应该发生），标记测试失败
	t.Error("Expected database connection to fail")
}
