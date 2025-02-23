package tools

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	"github.com/cloudwego/eino/components/tool"
)

func NewTools(ctx context.Context) (tools []tool.BaseTool, err error) {
	config := &duckduckgo.Config{
		MaxResults: 5, // Limit to return 3 results
		Region:     ddgsearch.RegionCN,
		DDGConfig: &ddgsearch.Config{
			Timeout:    15 * time.Second,
			Cache:      true,
			MaxRetries: 5,
		},
	}
	duckduckTools, err := duckduckgo.NewTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{
		duckduckTools,
	}, nil
}