package gomcpsdk

import (
	sdkClient "github.com/mark3labs/mcp-go/client"
)

type MCPClient = sdkClient.MCPClient
type ClientOption = sdkClient.ClientOption


func NewStdioMCPClient(command string, env []string, args ...string) (MCPClient, error) {
	return sdkClient.NewStdioMCPClient(command, env, args...)
}