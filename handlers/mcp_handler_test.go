package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMcpHandler_Normal_Origin_Header_isHeaderValid(t *testing.T) {

	// Given: Origin 헤더 값 세팅
	handler := NewMCPHandler(nil, nil)

	header := http.Header{
		"Origin": {"http://localhost"},
	}

	// When: Origin 헤더가 정상적으로 존재
	result := handler.isHeaderValid(header)

	// Then: Handler Handle
	assert.Equal(t, true, result)
}

func TestMcpHandler_Wrong_Origin_Header_isHeaderValid(t *testing.T) {

	// Given: Origin 헤더 값 세팅
	handler := NewMCPHandler(nil, nil)

	header := http.Header{
		"Origin": {"http://wrong_origin:1234"},
	}

	// When: Origin 헤더가 비정상으로 존재
	result := handler.isHeaderValid(header)

	// Then: Handler Handle
	assert.Equal(t, false, result)
}

func TestMcpHandler_Missing_Origin_Header_isHeaderValid(t *testing.T) {

	// Given: Origin 헤더 없음
	handler := NewMCPHandler(nil, nil)

	header := http.Header{}

	// When: Origin 헤더가 존재하지 않음
	result := handler.isHeaderValid(header)

	// Then: false 반환
	assert.Equal(t, false, result)
}
