package register

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

func handleSSE(endpoint string) (client.MCPClient, error) {
	sseClient, err := client.NewSSEMCPClient(endpoint)
	if err != nil {
		return nil, err
	}
	err = sseClient.Start(context.Background())
	return sseClient, err
}

func handleStreamable(endpoint string) (client.MCPClient, error) {
	sseClient, err := client.NewStreamableHttpClient(endpoint)
	if err != nil {
		return nil, err
	}
	err = sseClient.Start(context.Background())
	return sseClient, err
}

func InitClient(endpoint, transportType string) (client.MCPClient, error) {
	var mcpClient client.MCPClient
	var err error
	if "streamable" == transportType {
		mcpClient, err = handleStreamable(endpoint)
	} else {
		mcpClient, err = handleSSE(endpoint)
	}
	if err != nil {
		return nil, err
	}

	request := mcp.InitializeRequest{}
	request.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	request.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp register action",
		Version: "1.0",
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	initialize, err := mcpClient.Initialize(ctx, request)
	if err != nil {
		return nil, err
	}

	logrus.Infof(
		"Initialized with server: %s %s\n\n",
		initialize.ServerInfo.Name,
		initialize.ServerInfo.Version,
	)
	return mcpClient, nil
}

// RemoveAnyOf removes anyOf fields from tool schemas
func removeAnyOf(tools *mcp.ListToolsResult) {
	for _, tool := range tools.Tools {
		processAnyOf(tool.InputSchema.Properties)
	}
}

// processAnyOf processes anyOf fields in tool schemas
func processAnyOf(obj interface{}) {
	switch v := obj.(type) {
	case map[string]interface{}:
		// 如果存在 anyOf
		if anyOf, ok := v["anyOf"]; ok {
			if list, ok := anyOf.([]interface{}); ok && len(list) > 0 {
				if first, ok := list[0].(map[string]interface{}); ok {
					// 移除 anyOf
					delete(v, "anyOf")
					// 替换为第一个对象的内容
					for k, val := range first {
						v[k] = val
					}
				}
			}
		}
		// 递归处理 map 中的所有值
		for _, val := range v {
			processAnyOf(val)
		}
	case []interface{}:
		for _, item := range v {
			processAnyOf(item)
		}
	}
}
