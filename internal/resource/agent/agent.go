package agent

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

func NewAgent(ctx context.Context, chatModel model.ChatModel, tools []tool.BaseTool) (agent *react.Agent, err error) {
	return react.NewAgent(
		ctx,
		&react.AgentConfig{
			Model:           chatModel,
			ToolsConfig:     compose.ToolsNodeConfig{Tools: tools},
			MessageModifier: react.NewPersonaModifier("You are a helpful assistant"),
		},
	)

}
