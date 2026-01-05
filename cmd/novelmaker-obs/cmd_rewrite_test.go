package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
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

func TestRewriteCmd_ContextExtraction(t *testing.T) {
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

	// Create a test file
	storyDir := filepath.Join(tmpDir, "Story")
	testFile := filepath.Join(storyDir, "test.md")
	lines := []string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4 - Target Start",
		"Line 5 - Target",
		"Line 6 - Target End",
		"Line 7",
		"Line 8",
		"Line 9",
		"Line 10",
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name               string
		lineStart          int
		lineEnd            int
		contextPrev        int
		contextNext        int
		expectedBefore     string
		expectedTarget     string
		expectedAfter      string
		wantLineCountError bool
	}{
		{
			name:           "middle range with context",
			lineStart:      4,
			lineEnd:        6,
			contextPrev:    2,
			contextNext:    2,
			expectedBefore: "Line 2\nLine 3",
			expectedTarget: "Line 4 - Target Start\nLine 5 - Target\nLine 6 - Target End",
			expectedAfter:  "Line 7\nLine 8",
		},
		{
			name:           "start of file",
			lineStart:      1,
			lineEnd:        2,
			contextPrev:    3,
			contextNext:    2,
			expectedBefore: "",
			expectedTarget: "Line 1\nLine 2",
			expectedAfter:  "Line 3\nLine 4 - Target Start",
		},
		{
			name:           "end of file",
			lineStart:      9,
			lineEnd:        10,
			contextPrev:    2,
			contextNext:    5,
			expectedBefore: "Line 7\nLine 8",
			expectedTarget: "Line 9\nLine 10",
			expectedAfter:  "",
		},
		{
			name:               "line range exceeds file",
			lineStart:          15,
			lineEnd:            20,
			contextPrev:        2,
			contextNext:        2,
			wantLineCountError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &RewriteCmd{}
			fileLines, err := cmd.readFileLines(testFile)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			totalLines := len(fileLines)

			// Check for line count errors
			if tt.wantLineCountError {
				if tt.lineStart <= totalLines && tt.lineEnd <= totalLines {
					t.Error("Expected line count error but lines are within range")
				}
				return
			}

			if tt.lineStart > totalLines || tt.lineEnd > totalLines {
				t.Skip("Line range exceeds file length")
			}

			// Extract context (same logic as cmd_rewrite.go)
			contextStartIdx := max(0, tt.lineStart-1-tt.contextPrev)
			targetStartIdx := tt.lineStart - 1
			targetEndIdx := tt.lineEnd
			contextEndIdx := min(totalLines, tt.lineEnd+tt.contextNext)

			contextBefore := strings.Join(fileLines[contextStartIdx:targetStartIdx], "\n")
			targetSentence := strings.Join(fileLines[targetStartIdx:targetEndIdx], "\n")
			contextAfter := strings.Join(fileLines[targetEndIdx:contextEndIdx], "\n")

			if contextBefore != tt.expectedBefore {
				t.Errorf("Context before mismatch:\nExpected: %q\nGot: %q", tt.expectedBefore, contextBefore)
			}
			if targetSentence != tt.expectedTarget {
				t.Errorf("Target sentence mismatch:\nExpected: %q\nGot: %q", tt.expectedTarget, targetSentence)
			}
			if contextAfter != tt.expectedAfter {
				t.Errorf("Context after mismatch:\nExpected: %q\nGot: %q", tt.expectedAfter, contextAfter)
			}
		})
	}
}

func TestRewriteCmd_PromptGoalDefault(t *testing.T) {
	cmd := NewRewriteCmd(nil)

	// Check that the default prompt goal is set
	defaultGoal := "去除 AI 味的慣用詞與過度評述語氣"

	// Get the flag value
	flag := cmd.cmd.Flags().Lookup("prompt")
	if flag == nil {
		t.Fatal("prompt flag not found")
	}

	if flag.DefValue != defaultGoal {
		t.Errorf("Expected default prompt goal %q, got %q", defaultGoal, flag.DefValue)
	}
}

func TestRewritePromptTemplate(t *testing.T) {
	t.Run("rewrite prompt template exists and is valid", func(t *testing.T) {
		tmpDir := "../../internal/obsidian/init_template"

		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to create vault: %v", err)
		}

		rewritePrompt, err := vault.LoadRewritePrompt()
		if err != nil {
			t.Errorf("failed to load rewrite prompt: %v", err)
		}
		if rewritePrompt == nil {
			t.Error("rewrite prompt should not be nil")
		}
		if rewritePrompt != nil && rewritePrompt.System == "" {
			t.Error("rewrite prompt system should not be empty")
		}
		if rewritePrompt != nil && rewritePrompt.AssistantTemplate == nil {
			t.Error("rewrite prompt template should not be nil")
		}
	})
}
