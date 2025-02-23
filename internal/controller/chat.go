package controller

import (
	"net/http"

	"github.com/davlin-coder/davlin/internal/model"
	"github.com/davlin-coder/davlin/internal/service"
	"github.com/gin-gonic/gin"
)

// ChatController 定义聊天控制器接口
type ChatController interface {
	SendMessage(c *gin.Context)
	GetChatHistory(c *gin.Context)
}

// chatController 实现ChatController接口的结构体
type chatController struct {
	chatService service.ChatService
}

// NewChatController 创建聊天控制器实例
func NewChatController(chatService service.ChatService) ChatController {
	return &chatController{
		chatService: chatService,
	}
}

// SendMessage 发送聊天消息
func (ctrl *chatController) SendMessage(c *gin.Context) {
	var message model.ChatMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息格式"})
		return
	}

	response, err := ctrl.chatService.SendMessage(&message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetChatHistory 获取聊天历史
func (ctrl *chatController) GetChatHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	history, err := ctrl.chatService.GetHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}