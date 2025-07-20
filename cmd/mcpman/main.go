package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"os"
)

var CLI struct {
	install struct {
		url []string `help:"registry urls to install mcp servers"`
	} `cmd:"" help:"Install MCP servers from registry"`
	rm struct {
		path []string `arg:"" required:"" help:"Path to remove from MCP servers"`
	} `cmd:"" help:"Remove a path from MCP servers"`
	ls []string `cmd:"" help:"List all paths in MCP servers"`
}

func main() {
	kong.Parse(&CLI)
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
		case "install <url>":
			for _, url := range ctx.Args {
				fmt.Printf("Installing MCP server from URL: %s\n", url)
			}
		case "rm <path>":
			fmt.Printf("Removing path: %s\n", ctx.Args[0])
		case "ls":
			fmt.Println("Listing all paths")
		default:
			fmt.Println("No valid command provided. Use 'rm <path>' to remove a path or 'ls' to list all paths.")
	}
	os.Exit(0)
}
