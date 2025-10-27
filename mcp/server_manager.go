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
	mcpServer             *mcpsdk.Server
	streamableHTTPHandler *mcpsdk.StreamableHTTPHandler
	toolSet               map[string]models.Tool
}

type SayHiParams struct {
	Name string `json:"name"`
}

func NewServerManager() *ServerManager {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	streamableHTTPHandler := mcpsdk.NewStreamableHTTPHandler(func(req *http.Request) *mcpsdk.Server {
		return server
	}, nil)

	return &ServerManager{
		mcpServer:             server,
		streamableHTTPHandler: streamableHTTPHandler,
		toolSet:               make(map[string]models.Tool),
	}
}

// ServeHTTP implements http.Handler interface
func (serverManager *ServerManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serverManager.streamableHTTPHandler.ServeHTTP(w, r)
}

func echo(ctx context.Context, req *mcpsdk.CallToolRequest, args SayHiParams) (*mcpsdk.CallToolResult, any, error) {

	// typed handler로 args가 이미 들어온 경우:
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: args.Name},
		},
	}, nil, nil
}

func (serverManager *ServerManager) convertModelToolToMcpTool(modelTool models.Tool) (*mcpsdk.Tool, error) {
	log.Printf("Converting tool: %s, InputSchema raw: %s", modelTool.Name, string(modelTool.InputSchema))

	inputSchema, err := utils.ConvertGormJsonToObject(modelTool.InputSchema)
	if err != nil {
		log.Printf("Failed to convert InputSchema for %s: %v", modelTool.Name, err)
		return nil, fmt.Errorf("invalid JSON convert")
	}

	log.Printf("Converted InputSchema: %+v (type: %T)", inputSchema, inputSchema)
	log.Printf("OutputSchema raw: %s", string(modelTool.OutputSchema))

	outputSchema, err := utils.ConvertGormJsonToObject(modelTool.OutputSchema)
	if err != nil {
		log.Printf("Failed to convert OutputSchema for %s: %v", modelTool.Name, err)
		return nil, fmt.Errorf("invalid JSON convert")
	}

	log.Printf("Converted OutputSchema: %+v (type: %T)", outputSchema, outputSchema)

	return &mcpsdk.Tool{
		Description:  modelTool.Description,
		InputSchema:  inputSchema,
		Name:         modelTool.Name,
		OutputSchema: outputSchema,
	}, nil
}

// TODO: Add Tool Wrapping 함수
func (serverManager *ServerManager) DynamicAddTool(tool models.Tool) error {
	mcpTool, err := serverManager.convertModelToolToMcpTool(tool)
	if err != nil {
		log.Printf("Failed to convert tool %s: %v", tool.Name, err)
		return fmt.Errorf("failed to convert tool: %w", err)
	}

	log.Printf("Adding tool: %s, InputSchema: %+v", tool.Name, mcpTool.InputSchema)

	mcpsdk.AddTool(serverManager.mcpServer, mcpTool, echo)
	serverManager.toolSet[tool.Name] = tool
	return nil
}

// TODO: Delete Tool Wrapping 함수
func (serverManager *ServerManager) DynamicRemoveTool(toolName string) {
	delete(serverManager.toolSet, toolName)
	serverManager.mcpServer.RemoveTools(toolName)
}

func (serverManager *ServerManager) DynamicRemoveAllTool() {
	for toolName := range serverManager.toolSet {
		delete(serverManager.toolSet, toolName)
		serverManager.mcpServer.RemoveTools(toolName)
	}
}

func (serverManager *ServerManager) ToolList() []string {
	toolNames := make([]string, 0)
	i := 0
	for toolName := range serverManager.toolSet {
		toolNames[i] = toolName
		i++
	}
	return toolNames
}
