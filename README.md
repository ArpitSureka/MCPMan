# MCPMan (MCP Manager)

An MCP Server to download & manage MCP Servers. 
Idea is to enable non coders to use mcp servers without them installing pip / nodejs to install mcp servers.  
<br>
Currently supporting only Mac

## Installation
MacOS
```
brew install mcpman
```

## Installing MCP Servers 

PIP MCP : `mcpman pip mcp-server-fetch`

NPM MCP : `mcpman npm @playwright/mcp`

Local MCP : `mcpman install <path_to_executable>`

Github MCP : `mcpman github <github_url_of_exe>`


## Upcoming Features
- Automatic dowloading of required MCP Servers
- Supporting SSE based MCP servers. currently supporting only stdio servers
- Support For A2A Framework