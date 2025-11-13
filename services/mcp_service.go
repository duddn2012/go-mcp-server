package services

import (
	"fmt"
	"go_mcp_server/mcp"
	"go_mcp_server/models"
	"log"

	"gorm.io/gorm"
)

type MCPService struct {
	db *gorm.DB
}

func NewMCPService(db *gorm.DB) *MCPService {
	return &MCPService{db: db}
}

func (s *MCPService) GetAllTools() ([]models.Tool, error) {
	var tools []models.Tool
	if err := s.db.Find(&tools).Error; err != nil {
		return nil, fmt.Errorf("get all tools failed: %w", err)
	}
	return tools, nil
}

func (s *MCPService) GetEnabledTools() ([]models.Tool, error) {
	var tools []models.Tool
	if err := s.db.Where(&models.Tool{Enabled: true}).Find(&tools).Error; err != nil {
		return nil, fmt.Errorf("get enabled tools failed: %w", err)
	}
	return tools, nil
}

func (s *MCPService) SyncTools(sm *mcp.ServerManager) error {
	tools, err := s.GetEnabledTools()
	if err != nil {
		return err
	}

	sm.RemoveAllTools()

	successCount := 0
	for _, tool := range tools {
		if err := sm.AddTool(tool); err != nil {
			log.Printf("[MCPService] Failed to add tool %s: %v", tool.Name, err)
			continue
		}
		successCount++
	}

	log.Printf("[MCPService] Tools synced: %d/%d", successCount, len(tools))
	return nil
}
