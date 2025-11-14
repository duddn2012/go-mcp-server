package main

import (
	"go_mcp_server/internal/handler"
	"go_mcp_server/internal/infrastructure/config"
	"go_mcp_server/internal/infrastructure/database"
	"go_mcp_server/internal/mcp"
	"go_mcp_server/internal/repository"
	"go_mcp_server/internal/router"
	"go_mcp_server/internal/service"
	"log"
	"net/http"
)

func main() {
	// 1. Config 로드
	cfg := config.Get()

	// 2. DB 연결 생성
	db := database.GetDB(cfg)

	// 3. Repository 생성
	toolRepo := repository.NewToolRepository(db)

	// 4. Service 생성
	toolService := service.NewToolService(toolRepo)

	// 5. MCP Server Manager 생성
	mcpServerManager := mcp.NewServerManager()

	// 6. Handler 생성
	mcpHandler := handler.NewMCPHandler(toolService, mcpServerManager)

	// 7. Router 설정
	r := router.SetupRouter(mcpHandler)

	// 8. 서버 시작 (MCP 장시간 연결을 위한 타임아웃 설정)
	port := cfg.ServerPort
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    0, // MCP 연결을 위해 무제한
		WriteTimeout:   0, // MCP 연결을 위해 무제한
		IdleTimeout:    0, // MCP 연결을 위해 무제한
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Starting server on port %s with extended timeouts for MCP connections", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run server: %v", err)
	}
}
