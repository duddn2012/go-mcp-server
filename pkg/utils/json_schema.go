package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// JSONToMap은 datatypes.JSON을 map[string]any로 변환합니다.
// type 필드가 없으면 "object"를 추가합니다.
func JSONToMap(data datatypes.JSON) map[string]any {
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
