package mcp

import (
	"context"
	"fmt"
	"go_mcp_server/models"
	"go_mcp_server/utils"
	"log"
	"net/http"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerManager struct {
	server  *mcpsdk.Server
	handler *mcpsdk.StreamableHTTPHandler
	tools   map[string]models.Tool
}

type SayHiParams struct {
	Name string `json:"name"`
}

func NewServerManager() *ServerManager {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    "mcp-server",
		Version: "v1.0.0",
	}, nil)

	handler := mcpsdk.NewStreamableHTTPHandler(func(req *http.Request) *mcpsdk.Server {
		return server
	}, nil)

	return &ServerManager{
		server:  server,
		handler: handler,
		tools:   make(map[string]models.Tool),
	}
}

func (sm *ServerManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sm.handler.ServeHTTP(w, r)
}

func (sm *ServerManager) AddTool(tool models.Tool) error {
	mcpTool, err := sm.convertToMCPTool(tool)
	if err != nil {
		return fmt.Errorf("convert tool failed: %w", err)
	}

	mcpsdk.AddTool(sm.server, mcpTool, echo)
	sm.tools[tool.Name] = tool

	log.Printf("[ServerManager] Tool added: %s", tool.Name)
	return nil
}

func (sm *ServerManager) RemoveTool(name string) {
	sm.server.RemoveTools(name)
	delete(sm.tools, name)
	log.Printf("[ServerManager] Tool removed: %s", name)
}

func (sm *ServerManager) RemoveAllTools() {
	names := make([]string, 0, len(sm.tools))
	for name := range sm.tools {
		names = append(names, name)
	}

	if len(names) > 0 {
		sm.server.RemoveTools(names...)
		sm.tools = make(map[string]models.Tool)
		log.Printf("[ServerManager] All tools removed: %d", len(names))
	}
}

func (sm *ServerManager) GetToolNames() []string {
	names := make([]string, 0, len(sm.tools))
	for name := range sm.tools {
		names = append(names, name)
	}
	return names
}

func (sm *ServerManager) convertToMCPTool(tool models.Tool) (*mcpsdk.Tool, error) {
	return &mcpsdk.Tool{
		Name:         tool.Name,
		Description:  tool.Description,
		InputSchema:  utils.ToJSONSchema(tool.InputSchema),
		OutputSchema: utils.ToJSONSchema(tool.OutputSchema),
	}, nil
}

func echo(ctx context.Context, req *mcpsdk.CallToolRequest, args SayHiParams) (*mcpsdk.CallToolResult, any, error) {
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: args.Name},
		},
	}, nil, nil
}
