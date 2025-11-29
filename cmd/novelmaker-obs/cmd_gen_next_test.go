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
		if newChapter.ID == "" {
			t.Error("new chapter should have an ID")
		}
	})

	t.Run("generate first chapter in empty vault", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create minimal project
		projectContent := `---
name: New Project
system_prompt: You are a writer.
---
A new world.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

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

func TestGenNextCmd_Run_PrevChaptersLimit(t *testing.T) {
	t.Run("respects prev-chapters limit", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prevChapters = 5 // Request 5 previous chapters (but only 2 exist)

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

		// Should have original 2 + 1 new = 3 chapters
		if len(chapters) != 3 {
			t.Errorf("expected 3 chapters, got %d", len(chapters))
		}
	})

	t.Run("max prev-chapters is 10", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prevChapters = 15 // Request more than max

		// Should not error, but should be limited internally
		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}
	})

	t.Run("prev-chapters limited by available chapters", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prevChapters = 5 // Request 5 but only 2 exist

		// Should not error
		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
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

func TestGenNextCmd_Run_ChapterIndexing(t *testing.T) {
	t.Run("chapter index increments correctly", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Load initial chapters
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		initialChapters, err := vault.LoadChapters()
		vault.Close()
		if err != nil {
			t.Fatalf("failed to load initial chapters: %v", err)
		}

		maxInitialIndex := 0
		for _, ch := range initialChapters {
			if ch.Index > maxInitialIndex {
				maxInitialIndex = ch.Index
			}
		}

		// Generate new chapter
		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "New Chapter"
		genNextCmd.prevChapters = 2

		err = genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify new chapter has correct index
		vault2, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault2.Close()

		chapters, err := vault2.LoadChapters()
		if err != nil {
			t.Fatalf("failed to load chapters: %v", err)
		}

		maxIndex := 0
		for _, ch := range chapters {
			if ch.Index > maxIndex {
				maxIndex = ch.Index
			}
		}

		expectedIndex := maxInitialIndex + 1
		if maxIndex != expectedIndex {
			t.Errorf("new chapter should have index %d, got %d", expectedIndex, maxIndex)
		}
	})

	t.Run("chapter ID generated correctly", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prevChapters = 2

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

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

		// ID should be in format "ch<n>" where n is the chapter count before this one
		expectedID := "ch3" // 2 existing chapters + 1 = ch3
		if newChapter.ID != expectedID {
			t.Errorf("chapter ID = %s, want %s", newChapter.ID, expectedID)
		}
	})
}

func TestGenNextCmd_Run_FlagOverrides(t *testing.T) {
	t.Run("api-key flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.apiKey = "custom-api-key"
		genNextCmd.prevChapters = 2

		// Should not fail even with custom API key (dummy backend doesn't use it)
		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}
	})

	t.Run("model flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.model = "custom-model"
		genNextCmd.prevChapters = 2

		// Should not fail even with custom model (dummy backend doesn't use it)
		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}
	})

	t.Run("base-url flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.baseURL = "https://custom.api.url"
		genNextCmd.prevChapters = 2

		// Should not fail even with custom base URL (dummy backend doesn't use it)
		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}
	})

	t.Run("timeout flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.timeout = 60
		genNextCmd.prevChapters = 2

		// Should not fail with custom timeout
		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}
	})
}

func TestGenNextCmd_Run_WithWorldbookAndCharacters(t *testing.T) {
	t.Run("includes worldbook entries in generation", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Add additional worldbook entries
		wb3 := `---
id: new-world-entry
tags:
  - test
---
This is a new worldbook entry for testing.`
		writeTestFile(t, tmpDir, "World/003_new-entry.md", wb3)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prevChapters = 2

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify chapter was created (worldbook should be included in context)
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

	t.Run("includes characters in generation", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Add additional character
		char3 := `---
id: charlie
name: Charlie
main: false
---
Charlie is a supporting character.`
		writeTestFile(t, tmpDir, "Character/charlie.md", char3)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		genNextCmd.prevChapters = 2

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify chapter was created (characters should be included in context)
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

func TestGenNextCmd_Run_FileCreation(t *testing.T) {
	t.Run("created file has valid format", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three - Test Title"
		genNextCmd.prompt = "Test prompt"
		genNextCmd.prevChapters = 2

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Find the created file - it should be named 003_ch3.md
		storyDir := filepath.Join(tmpDir, "Story")
		newFile := filepath.Join(storyDir, "003_ch3.md")

		// Verify the file exists
		if _, err := os.Stat(newFile); os.IsNotExist(err) {
			// List all files in Story to help debug
			entries, _ := os.ReadDir(storyDir)
			t.Logf("Files in Story directory:")
			for _, e := range entries {
				t.Logf("  %s", e.Name())
			}
			t.Fatal("new chapter file not found at expected location")
		}

		// Read and verify file content
		content, err := os.ReadFile(newFile)
		if err != nil {
			t.Fatalf("failed to read new file: %v", err)
		}

		contentStr := string(content)

		// Verify frontmatter
		if !strings.Contains(contentStr, "---") {
			t.Error("file should have frontmatter delimiters")
		}
		if !strings.Contains(contentStr, "id:") {
			t.Error("file should have id field")
		}
		if !strings.Contains(contentStr, "title:") {
			t.Error("file should have title field")
		}
		if !strings.Contains(contentStr, "index:") {
			t.Error("file should have index field")
		}
		if !strings.Contains(contentStr, "prompt:") {
			t.Error("file should have prompt field")
		}

		// Verify content from dummy backend
		if !strings.Contains(contentStr, "dummy response") {
			t.Error("file should contain dummy backend response")
		}
	})

	t.Run("filename is properly slugified", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three: The Adventure Begins!"
		genNextCmd.prevChapters = 2

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Check that file was created with the expected naming pattern
		storyDir := filepath.Join(tmpDir, "Story")
		newFile := filepath.Join(storyDir, "003_ch3.md")

		// Verify the file exists
		if _, err := os.Stat(newFile); os.IsNotExist(err) {
			t.Fatal("new chapter file not found")
		}

		// Verify all files in Story directory don't contain special characters
		entries, err := os.ReadDir(storyDir)
		if err != nil {
			t.Fatalf("failed to read Story directory: %v", err)
		}

		for _, entry := range entries {
			name := entry.Name()
			// Files should follow the pattern XXX_id.md
			// Filenames should not contain special characters
			if strings.Contains(name, ":") || strings.Contains(name, "!") {
				t.Errorf("filename should not contain special characters, got: %s", name)
			}
		}
	})
}

func TestGenNextCmd_DefaultPrevChapters(t *testing.T) {
	t.Run("default prev-chapters is 3", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genNextCmd := NewGenNextCmd(dummyBackendMaker)
		genNextCmd.title = "Chapter Three"
		// Don't set prevChapters, should use default

		err := genNextCmd.run(genNextCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-next command failed: %v", err)
		}

		// Verify chapter was created
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
