package config

import (
    "os"
	"path/filepath"
    "gopkg.in/yaml.v3"
)

type Config struct {
	CacheDir string `yaml:"cache_dir"`
}

const (
    MCPServersDirName = "mcp_servers"
    ToolsCacheFileName = "tools.cache"
    ConfigFileName = "config.json"
)

// This package provides functions to manage configuration and cache directories for MCPMan.
// Current Folder Structure: For MacOS inside : ~/Library/Caches/mcpman
//     |- mcp_servers /
//     |         |- mcp_server_executables.exe
//     |         |- mcp_server_executable2.exe
//     |- tools.cache
//     |- config.json


// GetMCPServersDir returns the directory where MCP servers executables are stored. Ex return for MacOS: ~/Library/Caches/mcpman/mcp_servers
func GetMCPServersDir() (string, error) {
    destFolder, err := GetCacheDir()
	if err != nil {
        return "", err
	}
	destFolder = filepath.Join(destFolder, MCPServersDirName)
    err = os.MkdirAll(destFolder, 0o700) // Only readable/writeable by the user
	if err != nil {
		return "", err
	}
    return destFolder, nil
}

// GetToolsCacheFile returns the path to the tools cache file. Ex return for MacOS: ~/Library/Caches/mcpman/tools.cache
func GetToolsCacheFile() (string, error) {
    destFolder, err := GetCacheDir()
	if err != nil {
        return "", err
	}
    cacheFile := filepath.Join(destFolder, ToolsCacheFileName)
	err = os.WriteFile(cacheFile, []byte(""), 0o600) // Private to user
    if err != nil {
        return "", err
    }
    return cacheFile, nil
}

// GetConfigJSONFile returns the path to config file. Ex return for MacOS: ~/Library/Caches/mcpman/config.json
func GetConfigJSONFile() (string, error) {
    destFolder, err := GetCacheDir()
	if err != nil {
        return "", err
	}
    cacheFile := filepath.Join(destFolder, ConfigFileName)
	err = os.WriteFile(cacheFile, []byte(""), 0o600) // Private to user
    if err != nil {
        return "", err
    }
    return cacheFile, nil
}


// GetCacheDir returns the cache directory for MCPMan. Ex return for MacOS: ~/Library/Caches/mcpman
func GetCacheDir() (string, error) {
	cacheBase, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir := LoadConfigFromYAML().CacheDir
	mcpmanCache := filepath.Join(cacheBase, cacheDir)
	err = os.MkdirAll(mcpmanCache, 0o700)
	if err != nil {
		return "", err
	}

	return mcpmanCache, nil
}


// This function is used to load the configuration from a YAML file. (Loads config.yaml)
func LoadConfigFromYAML() (*Config) {
	filename := "config.yaml"

    file, err := os.Open(filename)
    if err != nil {
        return nil
    }
    defer file.Close()

    config := &Config{}
    decoder := yaml.NewDecoder(file)
    err = decoder.Decode(config)
    if err != nil {
        return nil
    }
    return config
}



