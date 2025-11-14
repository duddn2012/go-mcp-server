package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_mcp_server/internal/model"
	"go_mcp_server/pkg/utils"
	"io"
	"log"
	"net/http"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolManager struct {
	server *mcpsdk.Server
	tools  map[string]model.Tool
}

func NewToolManager(server *mcpsdk.Server) *ToolManager {
	return &ToolManager{
		server: server,
		tools:  make(map[string]model.Tool),
	}
}

func (tm *ToolManager) AddTool(tool model.Tool) error {
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
		tm.tools = make(map[string]model.Tool)
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

func (tm *ToolManager) convertToMCPTool(tool model.Tool) (*mcpsdk.Tool, error) {
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

	// body를 io.Reader로 변환
	var bodyReader io.Reader
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		// config와 args에서 body를 병합
		mergedBody := make(map[string]any)

		// 1. config의 body를 먼저 복사 (기본값)
		if bodyData, ok := config["body"].(map[string]any); ok {
			for k, v := range bodyData {
				mergedBody[k] = v
			}
		}

		// 2. args의 body로 덮어쓰기 (런타임 값 우선)
		if bodyData, ok := args["body"].(map[string]any); ok {
			for k, v := range bodyData {
				mergedBody[k] = v
			}
		}

		// 병합된 body가 있으면 JSON으로 직렬화
		if len(mergedBody) > 0 {
			jsonBody, err := json.Marshal(mergedBody)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal body: %w", err)
			}
			bodyReader = bytes.NewBuffer(jsonBody)
		}
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, nil, err
	}

	// Body가 있으면 Content-Type 헤더 설정
	if bodyReader != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// config와 args에서 headers를 병합
	mergedHeaders := make(map[string]string)

	// 1. config의 headers를 먼저 복사 (기본값)
	if headers, ok := config["headers"].(map[string]any); ok {
		for key, value := range headers {
			if strVal, ok := value.(string); ok {
				mergedHeaders[key] = strVal
			}
		}
	}

	// 2. args의 headers로 덮어쓰기 (런타임 값 우선)
	if headers, ok := args["headers"].(map[string]any); ok {
		for key, value := range headers {
			if strVal, ok := value.(string); ok {
				mergedHeaders[key] = strVal
			}
		}
	}

	// 병합된 headers 설정
	for key, value := range mergedHeaders {
		httpReq.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	result := map[string]any{
		"status": resp.StatusCode,
		"body":   string(respBody),
	}

	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: string(respBody)},
		},
	}, result, nil
}
