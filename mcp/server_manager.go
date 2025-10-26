package mcp

import (
	"context"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerManager struct {
	mcpServer             *mcp.Server
	streamableHTTPHandler *mcp.StreamableHTTPHandler
}

type SayHiParams struct {
	Name string `json:"name"`
}

func NewServerManager() *ServerManager {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	streamableHTTPHandler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	return &ServerManager{
		mcpServer:             server,
		streamableHTTPHandler: streamableHTTPHandler,
	}
}

// ServeHTTP implements http.Handler interface
func (serverManager *ServerManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serverManager.streamableHTTPHandler.ServeHTTP(w, r)
}

// TODO: 제거 예정
func echo(ctx context.Context, req *mcp.CallToolRequest, args SayHiParams) (*mcp.CallToolResult, any, error) {

	// typed handler로 args가 이미 들어온 경우:
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: args.Name},
		},
	}, nil, nil
}

// TODO: exec type 별 실행 함수 - api_call, echo
func (serverManager *ServerManager) executeTool(input string) string {
	return input
}

// TODO: Add Tool Wrapping 함수
func (serverManager *ServerManager) DynamicAddTool() {
	mcp.AddTool(serverManager.mcpServer, &mcp.Tool{Name: "greet", Description: "say hi"}, echo)
}
