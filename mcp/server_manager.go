package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"go_mcp_server/models"
	"go_mcp_server/utils"
	"io"
	"log"
	"net/http"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerManager struct {
	server  *mcpsdk.Server
	handler *mcpsdk.StreamableHTTPHandler
	tools   map[string]models.Tool
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

	if tool.ToolType == "echo" {
		mcpsdk.AddTool(sm.server, mcpTool, echo)
	} else {
		// Config를 캡처하는 클로저 생성
		config := utils.JSONToMap(tool.Config)
		handler := func(ctx context.Context, req *mcpsdk.CallToolRequest, args map[string]any) (*mcpsdk.CallToolResult, any, error) {
			return apiCall(ctx, req, args, config)
		}
		mcpsdk.AddTool(sm.server, mcpTool, handler)
	}

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
	mcpTool := &mcpsdk.Tool{
		Name:        tool.Name,
		Description: tool.Description,
		InputSchema: utils.JSONToMap(tool.InputSchema),
	}

	return mcpTool, nil
}

func echo(ctx context.Context, req *mcpsdk.CallToolRequest, args map[string]any) (*mcpsdk.CallToolResult, any, error) {

	var str string

	for key, value := range args {
		switch v := value.(type) {
		case string:
			str += fmt.Sprintf("%s: %s\n", key, v)
		case map[string]any, []any:
			// 복잡한 객체나 배열은 JSON으로 직렬화
			jsonBytes, _ := json.Marshal(v)
			str += fmt.Sprintf("%s: %s\n", key, string(jsonBytes))
		default:
			// 기타 타입 (숫자, bool 등)
			str += fmt.Sprintf("%s: %v\n", key, v)
		}
	}

	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: str},
		},
	}, nil, nil
}

func apiCall(ctx context.Context, req *mcpsdk.CallToolRequest, args map[string]any, config map[string]any) (*mcpsdk.CallToolResult, any, error) {
	// Config에서 안전하게 url, method 추출
	url, ok := config["url"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("url not found or invalid type in config")
	}

	method, ok := config["method"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("method not found or invalid type in config")
	}

	log.Printf("[apiCall] %s %s, args: %+v", method, url, args)

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, nil, err
	}

	// Headers 설정 (config에 있으면)
	if headers, ok := config["headers"].(map[string]any); ok {
		for key, value := range headers {
			if strVal, ok := value.(string); ok {
				httpReq.Header.Set(key, strVal)
			}
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	result := map[string]any{
		"status": resp.StatusCode,
		"body":   string(body),
	}

	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: string(body)},
		},
	}, result, nil
}
