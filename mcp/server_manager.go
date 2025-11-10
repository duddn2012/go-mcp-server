package mcp

import (
	"go_mcp_server/models"
	"net/http"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerManager struct {
	server      *mcpsdk.Server
	handler     *mcpsdk.StreamableHTTPHandler
	toolManager *ToolManager
}

func NewServerManager() *ServerManager {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    "mcp-server",
		Version: "v1.0.0",
	}, nil)

	handler := mcpsdk.NewStreamableHTTPHandler(func(req *http.Request) *mcpsdk.Server {
		return server
	}, nil)

	return &ServerManager{
		server:      server,
		handler:     handler,
		toolManager: NewToolManager(server),
	}
}

func (sm *ServerManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sm.handler.ServeHTTP(w, r)
}

func (sm *ServerManager) AddTool(tool models.Tool) error {
	return sm.toolManager.AddTool(tool)
}

func (sm *ServerManager) RemoveTool(name string) {
	sm.toolManager.RemoveTool(name)
}

func (sm *ServerManager) RemoveAllTools() {
	sm.toolManager.RemoveAllTools()
}

func (sm *ServerManager) GetToolNames() []string {
	return sm.toolManager.GetToolNames()
}
