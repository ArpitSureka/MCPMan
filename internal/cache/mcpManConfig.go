package cache

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/ArpitSureka/MCPMan/config"
	"github.com/ArpitSureka/MCPMan/internal/models"

)

// This File contains functions to read & modify config file inside cache. 
// Path of Config for MacOS: ~/Library/Caches/mcpman/config.json


// Function to read Config File inside Cache. 
func ReadConfigJSON() models.Config {
	var configData models.Config
	configFilePath, err := config.GetConfigJSONFile()
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("File doesn't exist. Creating new one.")
		// Initialize with default values
		configData = models.Config{
			MCPServers: []models.ServerJSON{},
		}
		saveConfigJSON(configData)
	} else {
		// Read file
		data, err := os.ReadFile(configFilePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		// Unmarshal JSON
		err = json.Unmarshal(data, &configData)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			os.Exit(1)
		}
	}

	return configData
}

// Function to read servers inside Cache Config File. 
func GetServerFromConfigJSON(name string) (models.ServerJSON, error) {
	configData := ReadConfigJSON()
	for _, server := range configData.MCPServers {
		if server.Name == name {
			return server, nil
		}
	}
	return models.ServerJSON{}, fmt.Errorf("server with name %s not found", name)
}

// Function to add servers inside Cache Config File. 
func AddServerToConfigJSON(serverJSON models.ServerJSON) error {
	configData := ReadConfigJSON()
	configData.MCPServers = append(configData.MCPServers, serverJSON)
	return saveConfigJSON(configData)
}


// Function to save data inside Cache Config File. 
func saveConfigJSON(configData models.Config) error {
	configFilePath, err := config.GetConfigJSONFile()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, data, 0644)
}