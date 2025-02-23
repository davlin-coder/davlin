package container

import (
	"context"

	"github.com/davlin-coder/davlin/internal/config"
	"github.com/davlin-coder/davlin/internal/controller"
	"github.com/davlin-coder/davlin/internal/resource/agent"
	"github.com/davlin-coder/davlin/internal/resource/email"
	"github.com/davlin-coder/davlin/internal/resource/llm"
	"github.com/davlin-coder/davlin/internal/resource/mysql"
	"github.com/davlin-coder/davlin/internal/resource/redis"
	"github.com/davlin-coder/davlin/internal/resource/template"
	"github.com/davlin-coder/davlin/internal/resource/tools"
	"github.com/davlin-coder/davlin/internal/router"
	"github.com/davlin-coder/davlin/internal/service"
	"go.uber.org/dig"
)

// Container 依赖注入容器
type Container struct {
	container *dig.Container
}

func appContext() context.Context {
	return context.Background()
}

// NewContainer 创建一个新的依赖注入容器
func NewContainer() (*Container, error) {
	container := dig.New()

	// 注册所有依赖
	providers := []interface{}{
		config.Init,
		appContext,

		// Resource层依赖
		mysql.Init,
		llm.NewModel,
		tools.NewTools,
		agent.NewAgent,
		template.NewTemplateManager,
		redis.NewRedisClient,
		email.NewEmailSender,
		tools.NewJWTManager,

		// Service层依赖
		service.NewUserService,
		service.NewChatService,
		service.NewVerificationService,

		// Controller层依赖
		controller.NewUserController,
		controller.NewChatController,

		// Router依赖
		router.NewRouter,
	}

	// 注册所有provider
	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, err
		}
	}

	return &Container{container: container}, nil
}

// Invoke 执行依赖注入
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
