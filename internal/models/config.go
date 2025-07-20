package models

// Structure of the config file stored inside cache
type Config struct {
	MCPServers []ServerJSON `json:"mcp_servers"`
}

// ServerJSON represents the JSON structure for a server configuration. config cache contains server details in this format.
type ServerJSON struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Type       ServerType `json:"type"`
	Version      string `json:"version"`
	STDIOServerConfig 
}

type STDIOServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}