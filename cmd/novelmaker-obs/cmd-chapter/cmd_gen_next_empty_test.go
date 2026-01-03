package cmdchapter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

func TestGenNextEmptyCmd(t *testing.T) {
	// Create a temporary directory for the test vault
	tmpDir, err := os.MkdirTemp("", "test-vault-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize vault structure
	vault, err := obsidian.NewVault(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create vault: %v", err)
	}
	defer vault.Close()

	if err := vault.Initialize(); err != nil {
		t.Fatalf("Failed to initialize vault: %v", err)
	}

	// Change to temp directory for test
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	tests := []struct {
		name      string
		title     string
		prompt    string
		wantError bool
	}{
		{
			name:      "basic empty chapter",
			title:     "Test Chapter",
			prompt:    "",
			wantError: false,
		},
		{
			name:      "empty chapter with prompt",
			title:     "Chapter with Prompt",
			prompt:    "Write about adventure",
			wantError: false,
		},
		{
			name:      "empty title",
			title:     "",
			prompt:    "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewGenNextEmptyCmd()
			cmd.title = tt.title
			cmd.prompt = tt.prompt
			cmd.json = true

			err := cmd.run(cmd.cmd, []string{})

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.wantError {
				// Verify file was created
				storyDir := filepath.Join(tmpDir, "Story")
				entries, err := os.ReadDir(storyDir)
				if err != nil {
					t.Fatalf("Failed to read Story directory: %v", err)
				}
				if len(entries) == 0 {
					t.Error("No chapter file was created")
				}
			}
		})
	}
}
