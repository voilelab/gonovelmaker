package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

func TestGenNextCmd_Run_Success(t *testing.T) {
	t.Run("generate next chapter with dummy backend", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three - The Great Adventure"
		genNextCmd.prompt = "Write about Alice meeting a dragon."
		genNextCmd.prevChapters = 2

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify the new chapter was created
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("failed to load chapters: %v", err)
		}

		if len(chapters) != 3 {
			t.Fatalf("expected 3 chapters after gen-next, got %d", len(chapters))
		}

		// Find the newly created chapter
		var newChapter *novelmaker.Chapter
		for i, ch := range chapters {
			if ch.Title == "Chapter Three - The Great Adventure" {
				newChapter = &chapters[i]
				break
			}
		}

		if newChapter == nil {
			t.Fatal("newly created chapter not found")
		}

		// Verify chapter properties
		if newChapter.Index != 3 {
			t.Errorf("new chapter index = %d, want 3", newChapter.Index)
		}
		if newChapter.Prompt != "Write about Alice meeting a dragon." {
			t.Errorf("new chapter prompt = %s, want 'Write about Alice meeting a dragon.'", newChapter.Prompt)
		}
		if !strings.Contains(newChapter.Content, "dummy response") {
			t.Error("new chapter should contain dummy backend response")
		}
	})

	t.Run("generate first chapter in empty vault", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create a .novelmaker config directory and config.toml file
		configDir := filepath.Join(tmpDir, ".novelmaker")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("failed to create config directory: %v", err)
		}
		configContent := `user_llm_backend = "test"

[llm_backend.test]
type = "openai"
api_key = "test-key"
model = "gpt-4o"
image_model = "dall-e-3"
`
		configPath := filepath.Join(configDir, "config.toml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("failed to write config file: %v", err)
		}

		// Set HOME environment variable to tmpDir so config.Load() finds our test config
		oldHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", oldHome)

		// Create minimal project
		projectContent := `---
name: New Project
system_prompt: You are a writer.
---
A new world.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create chapter prompt template
		chapterPromptContent := `---
system: |
    You are a professional novel writing assistant.
---
`
		writeTestFile(t, tmpDir, "Config/chapter_prompt.md", chapterPromptContent)

		// Create character prompt template
		characterPromptContent := `---
system: You are a character development AI assistant.
---
Create a detailed character profile.`
		writeTestFile(t, tmpDir, "Config/character_prompt.md", characterPromptContent)

		// Create empty directories
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Prologue"
		genNextCmd.prompt = "Write a prologue."
		genNextCmd.prevChapters = 3

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify the new chapter was created
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("failed to load chapters: %v", err)
		}

		if len(chapters) != 1 {
			t.Fatalf("expected 1 chapter, got %d", len(chapters))
		}

		// Verify chapter properties
		if chapters[0].Index != 1 {
			t.Errorf("first chapter index = %d, want 1", chapters[0].Index)
		}
		if chapters[0].Title != "Prologue" {
			t.Errorf("first chapter title = %s, want 'Prologue'", chapters[0].Title)
		}
	})

	t.Run("generate with custom prompt", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prompt = "This is a custom prompt with specific instructions."
		genNextCmd.prevChapters = 1

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify the prompt was saved
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("failed to load chapters: %v", err)
		}

		// Find the newly created chapter
		var newChapter *novelmaker.Chapter
		for i, ch := range chapters {
			if ch.Title == "Chapter Three" {
				newChapter = &chapters[i]
				break
			}
		}

		if newChapter == nil {
			t.Fatal("newly created chapter not found")
		}

		if newChapter.Prompt != "This is a custom prompt with specific instructions." {
			t.Errorf("chapter prompt = %s, want custom prompt", newChapter.Prompt)
		}
	})

	t.Run("generate without prompt", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prompt = "" // Empty prompt
		genNextCmd.prevChapters = 2

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify the chapter was created
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("failed to load chapters: %v", err)
		}

		if len(chapters) != 3 {
			t.Fatalf("expected 3 chapters, got %d", len(chapters))
		}
	})
}

func TestGenNextCmd_Run_JSONOutput(t *testing.T) {
	t.Run("json output format", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prompt = "Test prompt"
		genNextCmd.prevChapters = 2
		genNextCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := genNextCmd.run(genNextCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Parse JSON output
		var buf []byte
		buf, _ = io.ReadAll(r)
		output := string(buf)

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		// Verify JSON contains filepath
		filepath, ok := data["filepath"].(string)
		if !ok {
			t.Fatal("JSON output should contain 'filepath' field")
		}

		if filepath == "" {
			t.Error("filepath should not be empty")
		}

		// Verify the file exists
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Errorf("file should exist at path: %s", filepath)
		}
	})
}

func TestGenNextCmd_Run_ErrorCases(t *testing.T) {
	t.Run("error when project not found", func(t *testing.T) {
		tmpDir := createTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter One"

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when project not found, got nil")
		}

		if !strings.Contains(err.Error(), "failed to load project") {
			t.Errorf("error should mention 'failed to load project', got: %v", err)
		}
	})

	t.Run("error when config cannot be loaded", func(t *testing.T) {
		tmpDir := createTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		// Change to a directory that doesn't have a config
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter One"

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when config cannot be loaded")
		}
	})

	t.Run("error with empty title", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "" // Empty title

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error with empty title, got nil")
		}

		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("error should mention 'cannot be empty', got: %v", err)
		}
	})
}
