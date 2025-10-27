package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// ToJSONSchema는 GORM의 datatypes.JSON을 JSON Schema로 변환합니다.
// MCP Tool 인터페이스는 InputSchema와 OutputSchema가 반드시 type="object"여야 합니다.
func ToJSONSchema(data datatypes.JSON) map[string]any {
	if len(data) == 0 {
		return defaultObjectSchema()
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return defaultObjectSchema()
	}

	if _, hasType := schema["type"]; !hasType {
		schema["type"] = "object"
	}

	return schema
}

func defaultObjectSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}
