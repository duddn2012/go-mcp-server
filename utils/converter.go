package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func ConvertGormJsonToObject(inputJson datatypes.JSON) (any, error) {
	// 비어있으면 기본 object schema 반환
	if len(inputJson) == 0 {
		return map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}, nil
	}

	var result map[string]any
	if err := json.Unmarshal(inputJson, &result); err != nil {
		// 파싱 실패해도 기본 object schema 반환
		return map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}, nil
	}

	// type 필드가 없으면 "object" 추가
	if _, hasType := result["type"]; !hasType {
		result["type"] = "object"
	}

	return result, nil
}
