package service

import (
	"go_mcp_server/internal/model"
	"go_mcp_server/internal/repository"
	"go_mcp_server/test/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestGetAllTools(t *testing.T) {
	test := testutils.SetupTestDB(t)
	db := test.DB

	toolRepo := repository.NewToolRepository(db)
	service := NewToolService(toolRepo)

	tool1 := model.Tool{
		Name:        "test_tool_1",
		Description: "Test_tool_1",
		InputSchema: datatypes.JSON([]byte(`{"type":"object", "properties":{"name":{"type":"string"}}}`)),
		ToolType:    "api_call",
		Config:      datatypes.JSON([]byte(`{}`)),
		Enabled:     true,
	}
	tool2 := model.Tool{
		Name:        "test_tool_2",
		Description: "Test_tool_2",
		InputSchema: datatypes.JSON([]byte(`{"type":"object","properties":{"message":{"type":"string"}}}`)),
		ToolType:    "echo",
		Config:      datatypes.JSON([]byte(`{}`)),
		Enabled:     false,
	}
	db.Create(&tool1)
	db.Create(&tool2)

	tools, err := service.GetAllTools()
	assert.NoError(t, err)
	assert.Len(t, tools, 2)
}

func TestGetEnabledTools(t *testing.T) {
	test := testutils.SetupTestDB(t)
	db := test.DB

	toolRepo := repository.NewToolRepository(db)
	service := NewToolService(toolRepo)

	enabledTool := model.Tool{
		Name:        "enabled_tool",
		Description: "Enabled_tool",
		InputSchema: datatypes.JSON([]byte(`{"type":"object","properties":{}}`)),
		ToolType:    "api_call",
		Config:      datatypes.JSON([]byte(`{}`)),
		Enabled:     true,
	}
	db.Create(&enabledTool)

	// Disabled tool - false 값을 명시적으로 저장하기 위해 Update 사용
	disabledTool := model.Tool{
		Name:        "disabled_tool",
		Description: "Disabled_tool",
		InputSchema: datatypes.JSON([]byte(`{"type":"object","properties":{}}`)),
		ToolType:    "echo",
		Config:      datatypes.JSON([]byte(`{}`)),
	}
	db.Create(&disabledTool)
	db.Model(&disabledTool).Update("enabled", false)

	tools, err := service.GetEnabledTools()
	assert.NoError(t, err)
	assert.Len(t, tools, 1, "Should only return enabled tools")
	assert.Equal(t, "enabled_tool", tools[0].Name)
	assert.True(t, tools[0].Enabled)
}
