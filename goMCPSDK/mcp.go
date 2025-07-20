package gomcpsdk

import (
	sdkMCP "github.com/mark3labs/mcp-go/mcp"
)

type InitializeRequest = sdkMCP.InitializeRequest
type ClientCapabilities = sdkMCP.ClientCapabilities
type Implementation = sdkMCP.Implementation
type Tool = sdkMCP.Tool
type Prompt = sdkMCP.Prompt
type Resource = sdkMCP.Resource
type Root = sdkMCP.Root

const (
	LATEST_PROTOCOL_VERSION = sdkMCP.LATEST_PROTOCOL_VERSION
)