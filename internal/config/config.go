package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const configDirName = ".code-reviewer"
const configFileName = "config.json"

// Config holds the application configuration.
type Config struct {
	GoogleAIAPIKey string `json:"google_ai_api_key"`
}

// getConfigPath returns the full path to the configuration file.
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, configDirName, configFileName), nil
}

// Load reads the configuration from the local file.
// It returns the config and any error encountered.
// If the file does not exist, it returns an empty config and no error (or a specific error if preferred, but here we'll handle it by checking the key).
func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to the local file.
func Save(cfg *Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with 0600 permissions (read/write by owner only)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// PromptForAPIKey asks the user to input their Google AI API Key via stdin.
func PromptForAPIKey() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Google AI API Key not found. Please enter your Google AI API Key: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(text), nil
}
