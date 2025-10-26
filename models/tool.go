package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Tool은 MCP 도구를 나타내는 모델입니다
type Tool struct {
	gorm.Model
	Name         string         `gorm:"uniqueIndex;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	InputSchema  datatypes.JSON `gorm:"type:jsonb" json:"input_schema"`
	OutputSchema datatypes.JSON `gorm:"type:jsonb" json:"output_schema"`
	ToolType     string         `gorm:"not null;default:'api_call'" json:"tool_type"`
	Config       datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"config"`
	Enabled      bool           `gorm:"default:true" json:"enabled"`
}
