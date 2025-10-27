package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMCPHandler_ValidOrigin(t *testing.T) {
	handler := NewMCPHandler(nil, nil)

	header := http.Header{
		"Origin": {"http://localhost"},
	}

	result := handler.isOriginValid(header)

	assert.True(t, result)
}

func TestMCPHandler_InvalidOrigin(t *testing.T) {
	handler := NewMCPHandler(nil, nil)

	header := http.Header{
		"Origin": {"http://wrong_origin:1234"},
	}

	result := handler.isOriginValid(header)

	assert.False(t, result)
}

func TestMCPHandler_MissingOrigin(t *testing.T) {
	handler := NewMCPHandler(nil, nil)

	header := http.Header{}

	result := handler.isOriginValid(header)

	assert.False(t, result)
}
