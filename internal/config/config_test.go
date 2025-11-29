package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	// Setup temporary home directory
	tempHome, err := os.MkdirTemp("", "code-reviewer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	// Mock os.UserHomeDir by setting HOME env var
	// Note: This might not work if os.UserHomeDir implementation on the OS doesn't respect HOME,
	// but on Unix/Linux/Mac it usually does.
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempHome)

	// Test case 1: Load non-existent config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed for non-existent config: %v", err)
	}
	if cfg.GoogleAIAPIKey != "" {
		t.Errorf("Expected empty API key, got %q", cfg.GoogleAIAPIKey)
	}

	// Test case 2: Save and Load
	expectedKey := "test-api-key-123"
	cfg.GoogleAIAPIKey = expectedKey
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed after save: %v", err)
	}
	if loadedCfg.GoogleAIAPIKey != expectedKey {
		t.Errorf("Expected API key %q, got %q", expectedKey, loadedCfg.GoogleAIAPIKey)
	}

	// Test case 3: Save and Load with AI Model
	expectedModel := "gemini-pro"
	cfg.AIModel = expectedModel
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed with model: %v", err)
	}

	loadedCfgWithModel, err := Load()
	if err != nil {
		t.Fatalf("Load() failed after save with model: %v", err)
	}
	if loadedCfgWithModel.AIModel != expectedModel {
		t.Errorf("Expected AI Model %q, got %q", expectedModel, loadedCfgWithModel.AIModel)
	}

	// Test case 3: Check permissions
	configPath := filepath.Join(tempHome, configDirName, configFileName)
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", mode)
	}
}
