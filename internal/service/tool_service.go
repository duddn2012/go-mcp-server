package service

import (
	"go_mcp_server/internal/mcp"
	"go_mcp_server/internal/model"
	"go_mcp_server/internal/repository"
	"log"
)

type ToolService interface {
	GetAllTools() ([]model.Tool, error)
	GetEnabledTools() ([]model.Tool, error)
	SyncTools(sm *mcp.ServerManager) error
}

type toolService struct {
	toolRepo repository.ToolRepository
}

func NewToolService(toolRepo repository.ToolRepository) ToolService {
	return &toolService{
		toolRepo: toolRepo,
	}
}

func (s *toolService) GetAllTools() ([]model.Tool, error) {
	return s.toolRepo.FindAll()
}

func (s *toolService) GetEnabledTools() ([]model.Tool, error) {
	return s.toolRepo.FindEnabled()
}

func (s *toolService) SyncTools(sm *mcp.ServerManager) error {
	tools, err := s.GetEnabledTools()
	if err != nil {
		return err
	}

	sm.RemoveAllTools()

	successCount := 0
	for _, tool := range tools {
		if err := sm.AddTool(tool); err != nil {
			log.Printf("[ToolService] Failed to add tool %s: %v", tool.Name, err)
			continue
		}
		successCount++
	}

	log.Printf("[ToolService] Tools synced: %d/%d", successCount, len(tools))
	return nil
}
