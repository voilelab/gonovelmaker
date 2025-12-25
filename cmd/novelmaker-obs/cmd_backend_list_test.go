package main
package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/config"
)

func TestBackendListCmd(t *testing.T) {
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

	// Add some backends
	cfg := &config.Config{}
	cfg.AddOrUpdateBackend("backend1", config.LLMBackendConfig{
		Type:   "openai",
		APIKey: "key1",
		Model:  "gpt-4o",
	})
	cfg.AddOrUpdateBackend("backend2", config.LLMBackendConfig{
		Type:   "openrouter",
		APIKey: "key2",
		Model:  "claude-3",
	})
	cfg.UserLLMBackend = "backend1"
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// List backends (text output)
	cmd := NewBackendListCmd()
	err := cmd.run(cmd.cmd, []string{})


























































































}	}		}			t.Errorf("maskAPIKey(%q) = %q, want %q", tt.input, result, tt.expected)		if result != tt.expected {		result := maskAPIKey(tt.input)	for _, tt := range tests {	}		{"abcd", "...abcd"},		{"abc", "***"},		{"", ""},		{"key", "***"},		{"sk-1234567890abcdef", "...cdef"},	}{		expected string		input    string	tests := []struct {func TestMaskAPIKey(t *testing.T) {}	}		t.Error("Expected is_default to be true")	if result[0]["is_default"] != true {	}		t.Errorf("Expected backend name 'backend1', got '%v'", result[0]["name"])	if result[0]["name"] != "backend1" {	}		t.Errorf("Expected 1 backend in JSON output, got %d", len(result))	if len(result) != 1 {	}		t.Fatalf("Failed to parse JSON output: %v", err)	if err := json.Unmarshal([]byte(output), &result); err != nil {	var result []map[string]interface{}	// Verify JSON output	output := buf.String()	buf.ReadFrom(r)	var buf bytes.Buffer	// Read captured output	os.Stdout = oldStdout	w.Close()	// Restore stdout	}		t.Fatalf("Failed to list backends: %v", err)	if err != nil {	err := cmd.run(cmd.cmd, []string{})	cmd.jsonOutput = true	cmd := NewBackendListCmd()	// List backends (JSON output)	os.Stdout = w	r, w, _ := os.Pipe()	oldStdout := os.Stdout	// Capture stdout	}		t.Fatalf("Failed to save config: %v", err)	if err := cfg.Save(); err != nil {	cfg.UserLLMBackend = "backend1"	})		Model:  "gpt-4o",		APIKey: "key1",		Type:   "openai",	cfg.AddOrUpdateBackend("backend1", config.LLMBackendConfig{	cfg := &config.Config{}	// Add some backends	}		t.Fatalf("Failed to create config dir: %v", err)	if err := os.MkdirAll(configDir, 0755); err != nil {	configDir := filepath.Join(tmpDir, ".novelmaker")	// Create .novelmaker directory	defer os.Setenv("HOME", oldHome)	os.Setenv("HOME", tmpDir)	oldHome := os.Getenv("HOME")	tmpDir := t.TempDir()	// Create a temporary config directoryfunc TestBackendListCmdJSON(t *testing.T) {}	}		t.Fatalf("Failed to list backends: %v", err)	if err != nil {