package services

import (
	"fmt"
	"go_mcp_server/mcp"
	"go_mcp_server/models"

	"gorm.io/gorm"
)

type MCPService struct {
	db *gorm.DB
}

func NewMCPService(db *gorm.DB) *MCPService {
	return &MCPService{
		db: db,
	}
}

func (s *MCPService) GetAllTools() ([]models.Tool, error) {
	var tools []models.Tool
	if err := s.db.Find(&tools).Error; err != nil {
		return nil, err
	}
	return tools, nil
}

func (s *MCPService) GetEnabledTools() ([]models.Tool, error) {
	var tools []models.Tool

	if err := s.db.Where(&models.Tool{Enabled: true}).Find(&tools).Error; err != nil {
		return nil, err
	}
	return tools, nil
}

func (s *MCPService) SyncTools(mcpServerManager *mcp.ServerManager) error {
	// DB에서 활성화된 Tool들 가져오기
	dbTools, err := s.GetEnabledTools()
	if err != nil {
		return fmt.Errorf("failed to get tools: %w", err)
	}

	// 모든 Tool 제거
	mcpServerManager.DynamicRemoveAllTool()

	// 각 Tool을 MCP Server에 등록
	for _, tool := range dbTools {
		fmt.Printf("Registering tool: %s\n", tool.Name)
		if err := mcpServerManager.DynamicAddTool(tool); err != nil {
			fmt.Printf("Failed to register tool %s: %v\n", tool.Name, err)
			// 계속 진행 (하나 실패해도 나머지는 등록)
		}
	}

	return nil
}

func (s *MCPService) ExecuteTool(toolName string, input map[string]interface{}) (map[string]interface{}, error) {
	// TODO: 실제 MCP SDK를 사용해서 tool 실행

	if toolName == "echo" {
		return input, nil
	}

	return nil, fmt.Errorf("tool not found: %s", toolName)

}
