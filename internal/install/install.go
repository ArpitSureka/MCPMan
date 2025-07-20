package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ArpitSureka/MCPMan/config"
	"github.com/ArpitSureka/MCPMan/internal/cache"
	"github.com/ArpitSureka/MCPMan/internal/models"
)

func Install(urls []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("no URLs provided for installation")
	}

	for _, url := range urls {
		fmt.Printf("Downloading %s...\n", url)
		err := downloadExecutableFile(url)
		if err != nil {
			return fmt.Errorf("error downloading %s: %w", url, err)
		}
		fmt.Printf("Successfully downloaded %s\n", url)
	}
	return nil
}

func downloadExecutableFile(url string) error {
	// Check if the URL is a local file path (file:// or no scheme)
	if strings.HasPrefix(url, "file://") || (!strings.Contains(url, "://") && (filepath.IsAbs(url) || strings.HasPrefix(url, ".") || strings.HasPrefix(url, "/"))) {
		return downloadExecutableFileFromLocalPath(url)
	}

	return fmt.Errorf("unsupported URL format: %s. Please provide a local file path or a valid URL with a scheme", url)

	// Otherwise, treat it as a remote URL
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// destFolder, err := getCacheDir()
	// if err != nil {
	// 	return err
	// }

	// if resp.StatusCode != http.StatusOK {
	// 	return fmt.Errorf("failed to download %s: %s", url, resp.Status)
	// }

	// parts := strings.Split(url, "/")
	// filename := parts[len(parts)-1]
	// if filename == "" {
	// 	filename = "downloaded_file"
	// }
	// destPath := filepath.Join(destFolder, filename)

	// out, err := os.Create(destPath)
	// if err != nil {
	// 	return err
	// }
	// defer out.Close()

	// _, err = io.Copy(out, resp.Body)
	// return err
}

func downloadExecutableFileFromLocalPath(url string) error {
	var localPath string
	if strings.HasPrefix(url, "file://") {
		localPath = strings.TrimPrefix(url, "file://")
	} else {
		localPath = url
	}
	destFolder, err := config.GetMCPServersDir()
	if err != nil {
		return err
	}
	destFolder = filepath.Join(destFolder, "mcp_servers")

	parts := strings.Split(localPath, string(os.PathSeparator))
	filename := parts[len(parts)-1]
	if filename == "" {
		return fmt.Errorf("filename cannot be empty. Please provide a valid file path")
	}

	// Check if the file is an executable for the current OS
	isExecutable := false
	if runtime.GOOS == "windows" {
		isExecutable = strings.HasSuffix(strings.ToLower(filename), ".exe")
	} else {
		// On Unix-like systems, check for no extension or known executable extensions
		ext := filepath.Ext(filename)
		if ext == "" || ext == ".out" || ext == ".bin" {
			isExecutable = true
		} else {
			// Optionally, check if the file is actually executable
			fileInfo, err := os.Stat(localPath)
			if err == nil && fileInfo.Mode()&0111 != 0 {
				isExecutable = true
			}
		}
	}

	if !isExecutable {
		return fmt.Errorf("file %s is not recognized as an executable for %s", filename, runtime.GOOS)
	}

	serverName := strings.TrimSuffix(filename, filepath.Ext(filename))

	_, err = cache.GetServerFromConfigJSON(serverName)
	if err == nil {
		return fmt.Errorf("server with name %s already exists in the configuration", serverName)
	}

	destPath := filepath.Join(destFolder, filename)
	in, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	cache.AddServerToConfigJSON(models.ServerJSON{
		Name:    serverName,
		Path:    destPath,
		Type:    models.SSE,
		Version: "1.0.0", // Always set to 1.0.0 for local servers
	})

	return err

}


