package main

import (
	"fmt"

    "github.com/ArpitSureka/MCPMan/internal/mcpcontroller"
	"github.com/ArpitSureka/MCPMan/internal/models"
)

func main() {
	path := "/Users/a0s16ic/PersonalFiles/projects4/MCPMan/scripts/PythonMCPExecutables/mcp-server-fetch"
	server := models.ServerJSON{
		Name:    "TestServer",
		Path:   path,
		Type:    1,
		Version: "1.0.0",
		STDIOServerConfig: models.STDIOServerConfig{
			Command: path,
			Args:    []string{},
			Env:     map[string]string{},
		},
	}
	fmt.Printf("%+v\n", server)
	fmt.Printf("\n");
	fmt.Println("Fetching tools for server:", server.Name)
	fmt.Printf("\n");

	tools, err := mcpcontroller.GetMCPResources(server);
	if err != nil {
		fmt.Println("Error fetching tools:", err)
		return
	}
	fmt.Printf("\n");
	fmt.Println("Tools fetched successfully:")
	fmt.Printf("\n");
	fmt.Printf("%+v\n", tools);
	fmt.Printf("\n");
	fmt.Println("Total tools fetched:", len(tools))

}