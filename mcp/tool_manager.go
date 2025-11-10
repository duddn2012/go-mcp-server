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

type ToolManager struct {
	server *mcpsdk.Server
	tools  map[string]models.Tool
}

func NewToolManager(server *mcpsdk.Server) *ToolManager {
	return &ToolManager{
		server: server,
		tools:  make(map[string]models.Tool),
	}
}

func (tm *ToolManager) AddTool(tool models.Tool) error {
	mcpTool, err := tm.convertToMCPTool(tool)
	if err != nil {
		return fmt.Errorf("convert tool failed: %w", err)
	}

	if tool.ToolType == "echo" {
		mcpsdk.AddTool(tm.server, mcpTool, echo)
	} else {
		// Config를 캡처하는 클로저 생성
		config := utils.JSONToMap(tool.Config)
		handler := func(ctx context.Context, req *mcpsdk.CallToolRequest, args map[string]any) (*mcpsdk.CallToolResult, any, error) {
			return apiCall(ctx, req, args, config)
		}
		mcpsdk.AddTool(tm.server, mcpTool, handler)
	}

	tm.tools[tool.Name] = tool

	log.Printf("[ToolManager] Tool added: %s", tool.Name)
	return nil
}

func (tm *ToolManager) RemoveTool(name string) {
	tm.server.RemoveTools(name)
	delete(tm.tools, name)
	log.Printf("[ToolManager] Tool removed: %s", name)
}

func (tm *ToolManager) RemoveAllTools() {
	names := make([]string, 0, len(tm.tools))
	for name := range tm.tools {
		names = append(names, name)
	}

	if len(names) > 0 {
		tm.server.RemoveTools(names...)
		tm.tools = make(map[string]models.Tool)
		log.Printf("[ToolManager] All tools removed: %d", len(names))
	}
}

func (tm *ToolManager) GetToolNames() []string {
	names := make([]string, 0, len(tm.tools))
	for name := range tm.tools {
		names = append(names, name)
	}
	return names
}

func (tm *ToolManager) convertToMCPTool(tool models.Tool) (*mcpsdk.Tool, error) {
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

func apiCall(ctx context.Context, _ *mcpsdk.CallToolRequest, args map[string]any, config map[string]any) (*mcpsdk.CallToolResult, any, error) {
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
