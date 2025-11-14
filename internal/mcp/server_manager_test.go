package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerManager_ServeHTTP(t *testing.T) {

	// given
	sm := NewServerManager()

	// MCP 프로토콜의 초기화 요청
	requestBody := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	// ResponseRecorder를 사용하여 응답을 캡처
	rec := httptest.NewRecorder()

	// when
	sm.ServeHTTP(rec, req)

	// then
	// 1. ServeHTTP가 실행되어 응답 상태 코드가 설정되었는지 확인
	assert.Equal(t, http.StatusOK, rec.Code, "expected status 200")

	// 2. ServeHTTP가 실행되어 응답 본문이 작성되었는지 확인
	responseBody := rec.Body.String()
	assert.NotEmpty(t, responseBody, "expected non-empty response body")

	// 3. Content-Type 헤더가 설정되었는지 확인 (SSE 방식이므로 text/event-stream)
	contentType := rec.Header().Get("Content-Type")
	assert.Contains(t, contentType, "text/event-stream", "expected Content-Type to contain text/event-stream")

	// 4. SSE 응답에서 JSON 데이터 추출
	var jsonData string
	for _, line := range strings.Split(responseBody, "\n") {
		if data, found := strings.CutPrefix(line, "data: "); found {
			jsonData = data
			break
		}
	}
	assert.NotEmpty(t, jsonData, "SSE response should contain data field")

	// 5. JSON RPC 응답 구조 검증 - ServeHTTP가 올바른 형식으로 응답했는지 확인
	var response map[string]any
	err := json.Unmarshal([]byte(jsonData), &response)
	assert.NoError(t, err, "data field should be valid JSON")

	// 6. JSON RPC 필수 필드 확인
	assert.Contains(t, response, "jsonrpc", "response should have jsonrpc field")
	assert.Contains(t, response, "id", "response should have id field")
	assert.Equal(t, float64(1), response["id"], "id should match request id")

	// 7. ServeHTTP가 result 또는 error를 반환했는지 확인
	hasResult := response["result"] != nil
	hasError := response["error"] != nil
	assert.True(t, hasResult || hasError, "response must have either 'result' or 'error' field")
}
