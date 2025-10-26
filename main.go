package main

import (
	"go_mcp_server/config"
	"go_mcp_server/database"
	"go_mcp_server/handlers"
	"go_mcp_server/mcp"
	"go_mcp_server/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Config 로드
	cfg := config.Get()

	// 2. DB 연결 생성
	db := database.GetDBInstance(cfg)

	// 3. MCP Server 생성
	mcpServerManager := mcp.NewServerManager()

	// 4. Service 생성 (명시적 의존성 주입)
	mcpSvc := services.NewMCPService(db)

	// 5. Handler 생성 (명시적 의존성 주입)
	mcpHandler := handlers.NewMCPHandler(mcpSvc, mcpServerManager)

	// 6. Router 설정
	router := gin.Default()
	router.POST("/mcp/tools/sync", mcpHandler.HandleSyncTools)
	router.GET("/mcp", mcpHandler.HandleMcpServer)
	router.POST("/mcp", mcpHandler.HandleMcpServer)

	// 7. 서버 시작
	port := cfg.ServerPort
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run server: %v", err)
	}
}
