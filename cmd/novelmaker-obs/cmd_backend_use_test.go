package main
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/config"
)

func TestBackendUseCmd(t *testing.T) {
	// Create a temporary config directory
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Create .novelmaker directory
	configDir := filepath.Join(tmpDir, ".novelmaker")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Add a backend
	cfg := &config.Config{}
	cfg.AddOrUpdateBackend("test-backend", config.LLMBackendConfig{
		Type:   "openai",
		APIKey: "test-key",
	})










































}	}		t.Error("Expected error when setting non-existent backend as default")	if err == nil {	err := cmd.run(cmd.cmd, []string{"nonexistent"})	cmd := NewBackendUseCmd()	// Try to set a non-existent backend as default	}		t.Fatalf("Failed to create config dir: %v", err)	if err := os.MkdirAll(configDir, 0755); err != nil {	configDir := filepath.Join(tmpDir, ".novelmaker")	// Create .novelmaker directory	defer os.Setenv("HOME", oldHome)	os.Setenv("HOME", tmpDir)	oldHome := os.Getenv("HOME")	tmpDir := t.TempDir()	// Create a temporary config directoryfunc TestBackendUseCmdNotFound(t *testing.T) {}	}		t.Errorf("Expected default backend to be 'test-backend', got '%s'", cfg.UserLLMBackend)	if cfg.UserLLMBackend != "test-backend" {	}		t.Fatalf("Failed to load config: %v", err)	if err != nil {	cfg, err = config.Load()	// Verify the default was set	}		t.Fatalf("Failed to set default backend: %v", err)	if err != nil {	err := cmd.run(cmd.cmd, []string{"test-backend"})	cmd := NewBackendUseCmd()	// Set it as default	}		t.Fatalf("Failed to save initial config: %v", err)	if err := cfg.Save(); err != nil {