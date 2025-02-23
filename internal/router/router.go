package router

import (
	"github.com/davlin-coder/davlin/internal/controller"
	"github.com/davlin-coder/davlin/internal/middleware"
	"github.com/davlin-coder/davlin/internal/resource/tools"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userController   controller.UserController
	chatController   controller.ChatController
	healthController controller.HealthController
	jwtManager       *tools.JWTManager
}

func NewRouter(userController controller.UserController, chatController controller.ChatController, jwtManager *tools.JWTManager) *gin.Engine {
	healthController := controller.NewHealthController()
	router := &Router{
		userController:   userController,
		chatController:   chatController,
		healthController: healthController,
		jwtManager:       jwtManager,
	}
	return router.InitRouter()
}

func (r *Router) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 添加全局中间件
	router.Use(middleware.Cors())
	router.Use(middleware.Logger())

	// 健康检查路由
	router.GET("/health", r.healthController.Check)

	// API 版本分组
	v1 := router.Group("/api/v1")
	{
		// 用户相关路由
		userGroup := v1.Group("/user")
		{
			userGroup.POST("/register", r.userController.Register)
			userGroup.POST("/login", r.userController.Login)
			userGroup.POST("/verify-code", r.userController.SendVerificationCode)
		}

		// 需要认证的路由组
		authGroup := v1.Group("", middleware.Auth(r.jwtManager))
		{
			// 聊天相关路由
			chatGroup := authGroup.Group("/chat")
			{
				chatGroup.POST("/message", r.chatController.SendMessage)
				chatGroup.GET("/history", r.chatController.GetChatHistory)
			}
		}
	}

	return router
}
