package mcp

import (
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestToolManager_Execute_Echo(t *testing.T) {
	// Given: echo 함수 테스트
	input := map[string]any{"message": "Hello"}

	// When: echo 함수를 직접 호출
	result, _, err := echo(nil, nil, input)

	// Then: 에러 없이 입력이 텍스트로 반환되어야 함
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Content, 1)

	// Content[0]이 TextContent인지 확인
	textContent, ok := result.Content[0].(*mcpsdk.TextContent)
	assert.True(t, ok)
	assert.Contains(t, textContent.Text, "message: Hello")
}
