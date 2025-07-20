package mcpcontroller

import (
	"github.com/ArpitSureka/MCPMan/internal/cache"
	"github.com/ArpitSureka/MCPMan/internal/models"
)

var serverDB = models.GetServerDB()

func LoadAndRunMCPServers() {
	go loadMCPServers()

}

func loadMCPServers() {
	config := cache.ReadConfigJSON()
	for _, server := range config.MCPServers {
		serverDB.Servers <- models.Server{
			ServerJSON: server,
			Status:    1,
		}
	}
}