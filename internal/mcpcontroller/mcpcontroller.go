package mcpcontroller

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	"github.com/ArpitSureka/MCPMan/goMCPSDK"
	"github.com/ArpitSureka/MCPMan/internal/models"
)


func GetMCPTools(server models.ServerJSON) ([]gomcpsdk.Tool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := createMCPAndInitializeClient(server, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	log.Info("Fetching tools for server...", "name", server.Name, "version", server.Version)
	listToolsRequest := gomcpsdk.ListToolsRequest{}
	response, err := client.ListTools(ctx, listToolsRequest)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to list tools for %s: %w",
			server.Name,
			err,
		)
	}

	return response.Tools, nil
}

func GetMCPPrompts(server models.ServerJSON) ([]gomcpsdk.Prompt, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := createMCPAndInitializeClient(server, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	log.Info("Fetching prompts for server...", "name", server.Name, "version", server.Version)
	listPromptsRequest := gomcpsdk.ListPromptsRequest{}
	response, err := client.ListPrompts(ctx, listPromptsRequest)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to list prompts for %s: %w",
			server.Name,
			err,
		)
	}

	return response.Prompts, nil
}

func GetMCPResources(server models.ServerJSON) ([]gomcpsdk.Resource, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := createMCPAndInitializeClient(server, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	defer client.Close()

	log.Info("Fetching resources for server...", "name", server.Name, "version", server.Version)
	listResourcesRequest := gomcpsdk.ListResourcesRequest{}
	response, err := client.ListResources(ctx, listResourcesRequest)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to list resources for %s: %w",
			server.Name,
			err,
		)
	}

	return response.Resources, nil
}

// func GetMCPRoots(server models.ServerJSON) ([]gomcpsdk.Root, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	client, err := createMCPAndInitializeClient(server, ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create MCP client: %w", err)
// 	}
// 	defer client.Close()

// 	log.Info("Fetching roots for server...", "name", server.Name, "version", server.Version)
// 	listRootsRequest := gomcpsdk.ListRootsRequest{}
// 	response, err := client.ListRoots(ctx, listRootsRequest)
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"failed to list roots for %s: %w",
// 			server.Name,
// 			err,
// 		)
// 	}

// 	return response.Roots, nil
// }


// Not sure if this function is required or not. 
// Inspired from https://github.com/mark3labs/mcphost/blob/b4750c5852a4f8591814a7a990f5deb15303a770/cmd/mcp.go
func createMCPAndInitializeClient(server models.ServerJSON, ctx context.Context) (gomcpsdk.MCPClient, error) {

	var client gomcpsdk.MCPClient
	var err error

	client, err = gomcpsdk.NewStdioMCPClient(
		server.Command,
		[]string{},
		server.Args...,
	)
	
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create MCP client for %s: %w",
			server.Name,
			err,
		)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	log.Info("Initializing server...", "name", server.Name, "version", server.Version)
	initRequest := gomcpsdk.InitializeRequest{}
	initRequest.Params.ProtocolVersion = gomcpsdk.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = gomcpsdk.Implementation{
		Name:    "mcphost",
		Version: "0.1.0",
	}
	initRequest.Params.Capabilities = gomcpsdk.ClientCapabilities{}

	_, err = client.Initialize(ctx, initRequest)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf(
			"failed to initialize MCP client for %s: %w",
			server.Name,
			err,
		)
	}

	return client, nil
}
