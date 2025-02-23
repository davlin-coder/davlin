package service

import (
	"errors"

	"github.com/davlin-coder/davlin/internal/model"
	"gorm.io/gorm"
)

type ChatService interface {
	SendMessage(message *model.ChatMessage) (map[string]interface{}, error)
	GetHistory(userID uint) ([]model.ChatMessage, error)
}

type chatService struct {
	db *gorm.DB
}

func NewChatService(db *gorm.DB) ChatService {
	return &chatService{db: db}
}

func (s *chatService) SendMessage(message *model.ChatMessage) (map[string]interface{}, error) {
	if s.db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	// 保存消息到数据库
	result := s.db.Create(message)
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"status":     "success",
		"message_id": message.ID,
	}, nil
}

func (s *chatService) GetHistory(userID uint) ([]model.ChatMessage, error) {
	if s.db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var messages []model.ChatMessage
	result := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}