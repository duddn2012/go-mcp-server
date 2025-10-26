package handlers

import (
	"go_mcp_server/config"
	"go_mcp_server/mcp"
	"go_mcp_server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MCPHandler struct {
	service          *services.MCPService
	mcpServerManager *mcp.ServerManager
}

func NewMCPHandler(service *services.MCPService, mcpServerManager *mcp.ServerManager) *MCPHandler {
	return &MCPHandler{
		service:          service,
		mcpServerManager: mcpServerManager,
	}
}

func (h *MCPHandler) Handle(c *gin.Context) {
	// Origin Header 검증
	isValid := h.isHeaderValid(c.Request.Header)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Origin header",
		})
		return
	}

	// MCP SDK의 SSEHandler에 위임
	h.mcpServerManager.ServeHTTP(c.Writer, c.Request)
}

func (h *MCPHandler) isHeaderValid(header http.Header) bool {
	// Origin 헤더가 존재하지 않으면 false
	origins, exists := header["Origin"]
	if !exists || len(origins) == 0 {
		return false
	}

	config := config.Get()
	allowOrigins := config.AllowedOrigins
	for _, origin := range allowOrigins {
		if origins[0] == origin {
			return true
		}
	}

	return false
}
