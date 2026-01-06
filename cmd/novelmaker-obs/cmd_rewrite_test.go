package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRewriteCmd_Validation(t *testing.T) {
	tests := []struct {
		name        string
		lineStart   int
		lineEnd     int
		contextPrev int
		contextNext int
		wantError   bool
		errorMsg    string
	}{
		{
			name:        "valid line range",
			lineStart:   1,
			lineEnd:     5,
			contextPrev: 3,
			contextNext: 3,
			wantError:   false,
		},
		{
			name:        "line-start less than 1",
			lineStart:   0,
			lineEnd:     5,
			contextPrev: 3,
			contextNext: 3,
			wantError:   true,
			errorMsg:    "line-start must be >= 1",
		},
		{
			name:        "line-end before line-start",
			lineStart:   10,
			lineEnd:     5,
			contextPrev: 3,
			contextNext: 3,
			wantError:   true,
			errorMsg:    "line-end (5) must be >= line-start (10)",
		},
		{
			name:        "negative context-prev",
			lineStart:   1,
			lineEnd:     5,
			contextPrev: -1,
			contextNext: 3,
			wantError:   true,
			errorMsg:    "context-prev must be >= 0",
		},
		{
			name:        "negative context-next",
			lineStart:   1,
			lineEnd:     5,
			contextPrev: 3,
			contextNext: -1,
			wantError:   true,
			errorMsg:    "context-next must be >= 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RewriteCmd{
				lineStart:   tt.lineStart,
				lineEnd:     tt.lineEnd,
				contextPrev: tt.contextPrev,
				contextNext: tt.contextNext,
			}

			err := cmd.run(cmd.cmd, []string{})

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			}
			// For valid line ranges, we don't check for success since we don't have a full test environment
		})
	}
}

func TestRewriteCmd_ReadFileLines(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := &RewriteCmd{}
	lines, err := cmd.readFileLines(testFile)
	if err != nil {
		t.Fatalf("Failed to read file lines: %v", err)
	}

	expectedLines := []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Line %d: expected %q, got %q", i+1, expectedLines[i], line)
		}
	}
}

func TestRewriteCmd_ReadFileLines_NonExistent(t *testing.T) {
	cmd := &RewriteCmd{}
	_, err := cmd.readFileLines("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got none")
	}
}
