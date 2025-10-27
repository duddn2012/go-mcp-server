package handlers

import (
	"go_mcp_server/config"
	"go_mcp_server/mcp"
	"go_mcp_server/services"
	"log"
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

func (h *MCPHandler) HandleSyncTools(c *gin.Context) {
	if !h.isOriginValid(c.Request.Header) {
		log.Printf("[MCPHandler] Invalid Origin from %s", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Origin header"})
		return
	}

	if err := h.service.SyncTools(h.mcpServerManager); err != nil {
		log.Printf("[MCPHandler] Failed to sync tools: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync tools"})
		return
	}

	log.Printf("[MCPHandler] Tools synced successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Tools synced"})
}

func (h *MCPHandler) HandleMcpServer(c *gin.Context) {
	if !h.isOriginValid(c.Request.Header) {
		log.Printf("[MCPHandler] Invalid Origin from %s", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Origin header"})
		return
	}

	h.mcpServerManager.ServeHTTP(c.Writer, c.Request)
}

func (h *MCPHandler) isOriginValid(header http.Header) bool {
	origins, exists := header["Origin"]
	if !exists || len(origins) == 0 {
		return false
	}

	cfg := config.Get()
	for _, allowed := range cfg.AllowedOrigins {
		if origins[0] == allowed {
			return true
		}
	}

	return false
}
