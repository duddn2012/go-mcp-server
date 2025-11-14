package router

import (
	"go_mcp_server/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mcpHandler *handler.MCPHandler) *gin.Engine {
	router := gin.Default()

	// MCP endpoints
	router.POST("/mcp/tools/sync", mcpHandler.HandleSyncTools)
	router.GET("/mcp", mcpHandler.HandleMcpServer)
	router.POST("/mcp", mcpHandler.HandleMcpServer)

	return router
}
