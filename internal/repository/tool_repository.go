package repository

import (
	"go_mcp_server/internal/model"

	"gorm.io/gorm"
)

type ToolRepository interface {
	FindAll() ([]model.Tool, error)
	FindEnabled() ([]model.Tool, error)
}

type toolRepository struct {
	db *gorm.DB
}

func NewToolRepository(db *gorm.DB) ToolRepository {
	return &toolRepository{db: db}
}

func (r *toolRepository) FindAll() ([]model.Tool, error) {
	var tools []model.Tool
	err := r.db.Find(&tools).Error
	return tools, err
}

func (r *toolRepository) FindEnabled() ([]model.Tool, error) {
	var tools []model.Tool
	err := r.db.Where("enabled = ?", true).Find(&tools).Error
	return tools, err
}
