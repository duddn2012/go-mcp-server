package services

import (
	"go_mcp_server/models"
	"go_mcp_server/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestGetAllTools(t *testing.T) {
	test := testutils.SetupTestDB(t)
	db := test.DB
	service := &MCPService{db: db}

	tool1 := models.Tool{
		Name:        "test_tool_1",
		Description: "Test_tool_1",
		InputSchema: datatypes.JSON([]byte(`{"type":"object", "properties":{"name":{"type":"string"}}}`)),
		ToolType:    "api_call",
		Config:      datatypes.JSON([]byte(`{}`)),
		Enabled:     true,
	}
	tool2 := models.Tool{
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

// func TestGetEnabledTools(t *testing.T) {
// 	test := testutils.SetupTestDB(t)
// 	db := test.DB
// 	service := &MCPService{db: db}

// 	db.Create(&models.Tool{
// 		Name:        "enabled_tool",
// 		Description: "Enabled_tool",
// 		InputSchema: datatypes.JSON([]byte(`{"type":"object","properties":{}}`)),
// 		ToolType:    "api_call",
// 		Config:      datatypes.JSON([]byte(`{}`)),
// 		Enabled:     true,
// 	})
// 	db.Create(&models.Tool{
// 		Name:        "disabled_tool",
// 		Description: "Disabled_tool",
// 		InputSchema: datatypes.JSON([]byte(`{"type":"object","properties":{}}`)),
// 		ToolType:    "echo",
// 		Config:      datatypes.JSON([]byte(`{}`)),
// 		Enabled:     false,
// 	})

// 	tools, err := service.GetEnabledTools()
// 	assert.NoError(t, err)
// 	assert.Len(t, tools, 1)
// 	assert.Equal(t, "enabled_tool", tools[0].Name)
// 	assert.True(t, tools[0].Enabled)
// }

// func TestRegisterTools(t *testing.T) {
// 	// Given: DB에 Tool이 저장되어 있고
// 	test := testutils.SetupTestDB(t)
// 	db := test.DB
// 	service := NewMCPService(db)
// 	tool := models.Tool{
// 		Name:        "echo_tool",
// 		Description: "Simple echo tool",
// 		InputSchema: datatypes.JSON([]byte(`{"type":"object","properties":{"message":{"type":"string"}}}`)),
// 		ToolType:    "echo",
// 		Config:      datatypes.JSON([]byte(`{}`)),
// 		Enabled:     true,
// 	}
// 	db.Create(&tool)

// 	// When: MCP 서버에 Tool을 등록하면
// 	err := service.RegisterTools()

// 	// Then: 에러 없이 등록되어야 함
// 	assert.NoError(t, err)
// }

// func TestExecuteTool(t *testing.T) {
// 	// Given: echo tool이 등록되어 있고
// 	test := testutils.SetupTestDB(t)
// 	db := test.DB
// 	service := NewMCPService(db)

// 	tool := models.Tool{
// 		Name:        "echo",
// 		Description: "Echo back the message",
// 		InputSchema: datatypes.JSON([]byte(`{"type":"object","properties":{"message":{"type":"string"}}}`)),
// 		ToolType:    "echo",
// 		Config:      datatypes.JSON([]byte(`{}`)),
// 		Enabled:     true,
// 	}
// 	db.Create(&tool)
// 	service.RegisterTools()

// 	// When: Tool을 실행하면
// 	input := map[string]interface{}{"message": "Hello, MCP!"}
// 	result, err := service.ExecuteTool("echo", input)

// 	// Then: 입력한 메시지가 그대로 반환되어야 함
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Hello, MCP!", result["message"])
// }
