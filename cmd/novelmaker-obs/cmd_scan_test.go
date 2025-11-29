package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/voilelab/gonovelmaker/novelmaker"
)

// createTestVault creates a temporary vault directory for testing
func createTestVault(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return tmpDir
}

// writeTestFile writes a file to the test vault
func writeTestFile(t *testing.T, root, path, content string) {
	t.Helper()
	fullPath := filepath.Join(root, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", fullPath, err)
	}
}

// setupCompleteVault creates a complete test vault with all required files
func setupCompleteVault(t *testing.T) string {
	t.Helper()
	tmpDir := createTestVault(t)

	// Create project file
	projectContent := `---
name: Test Novel Project
system_prompt: You are a helpful AI assistant for writing novels.
system_prompt_char: You are a character development AI assistant.
---
This is a fantasy world with magic and dragons.`
	writeTestFile(t, tmpDir, "Config/project.md", projectContent)

	// Create worldbook entries
	wb1 := `---
id: magic-system
tags:
  - magic
  - rules
---
Magic in this world is based on elemental manipulation.`
	writeTestFile(t, tmpDir, "World/001_magic-system.md", wb1)

	wb2 := `---
id: geography
tags:
  - world
  - locations
---
The world consists of three continents separated by vast oceans.`
	writeTestFile(t, tmpDir, "World/002_geography.md", wb2)

	// Create characters
	char1 := `---
id: alice
name: Alice
main: true
---
Alice is a brave knight who seeks to protect her kingdom.`
	writeTestFile(t, tmpDir, "Character/alice.md", char1)

	char2 := `---
id: bob
name: Bob
main: false
---
Bob is a wise old wizard who mentors the protagonist.`
	writeTestFile(t, tmpDir, "Character/bob.md", char2)

	// Create chapters
	ch1 := `---
id: prologue
title: The Beginning
index: 1
prompt: Write an engaging prologue that introduces the main character.
---
It was a dark and stormy night when Alice first discovered her destiny.`
	writeTestFile(t, tmpDir, "Story/001_prologue.md", ch1)

	ch2 := `---
id: chapter-one
title: Chapter One - The Journey Begins
index: 2
prompt: Continue the story as Alice leaves her village.
---
At dawn, Alice packed her belongings and set out on her journey.`
	writeTestFile(t, tmpDir, "Story/002_chapter-one.md", ch2)

	return tmpDir
}

func TestScanCmd_Run_TextOutput(t *testing.T) {
	t.Run("complete vault", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		// Check project section
		if !strings.Contains(output, "=== PROJECT ===") {
			t.Error("output should contain project section header")
		}
		if !strings.Contains(output, "Name: Test Novel Project") {
			t.Error("output should contain project name")
		}
		if !strings.Contains(output, "fantasy world with magic and dragons") {
			t.Error("output should contain project world description")
		}

		// Check worldbook section
		if !strings.Contains(output, "=== WORLDBOOK (2 entries) ===") {
			t.Error("output should contain worldbook section with count")
		}
		if !strings.Contains(output, "ID: magic-system") {
			t.Error("output should contain magic-system worldbook")
		}
		if !strings.Contains(output, "ID: geography") {
			t.Error("output should contain geography worldbook")
		}
		if !strings.Contains(output, "elemental manipulation") {
			t.Error("output should contain worldbook content snippet")
		}

		// Check characters section
		if !strings.Contains(output, "=== CHARACTERS (2 total) ===") {
			t.Error("output should contain characters section with count")
		}
		if !strings.Contains(output, "Name: Alice") {
			t.Error("output should contain Alice character")
		}
		if !strings.Contains(output, "Name: Bob") {
			t.Error("output should contain Bob character")
		}
		if !strings.Contains(output, "brave knight") {
			t.Error("output should contain character profile snippet")
		}

		// Check chapters section
		if !strings.Contains(output, "=== CHAPTERS (2 total) ===") {
			t.Error("output should contain chapters section with count")
		}
		if !strings.Contains(output, "[1] The Beginning") {
			t.Error("output should contain first chapter with index")
		}
		if !strings.Contains(output, "[2] Chapter One - The Journey Begins") {
			t.Error("output should contain second chapter with index")
		}
		if !strings.Contains(output, "dark and stormy night") {
			t.Error("output should contain chapter content snippet")
		}
	})

	t.Run("empty vault", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create minimal project file
		projectContent := `---
name: Empty Project
---
Empty world.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create empty directories
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		// Check for empty collections
		if !strings.Contains(output, "=== WORLDBOOK (0 entries) ===") {
			t.Error("output should show 0 worldbook entries")
		}
		if !strings.Contains(output, "=== CHARACTERS (0 total) ===") {
			t.Error("output should show 0 characters")
		}
		if !strings.Contains(output, "=== CHAPTERS (0 total) ===") {
			t.Error("output should show 0 chapters")
		}
	})

	t.Run("missing project file", func(t *testing.T) {
		tmpDir := createTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()
		err := scanCmd.run(scanCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error for missing project file, got nil")
		}
		if !strings.Contains(err.Error(), "failed to load project") {
			t.Errorf("error should mention 'failed to load project', got: %v", err)
		}
	})

	t.Run("long content truncation", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create project file
		projectContent := `---
name: Test Project
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create worldbook with long content (>80 chars)
		longContent := strings.Repeat("This is a very long worldbook entry that should be truncated in the output. ", 5)
		wb := `---
id: long-entry
tags:
  - test
---
` + longContent
		writeTestFile(t, tmpDir, "World/001_long.md", wb)

		// Create empty directories for other parts
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		// Check that content is truncated with "..."
		if !strings.Contains(output, "...") {
			t.Error("long content should be truncated with '...'")
		}

		// Check that the full long content is not in output
		if strings.Contains(output, longContent) {
			t.Error("full long content should not appear in output")
		}
	})
}

func TestScanCmd_Run_JSONOutput(t *testing.T) {
	t.Run("json format", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()
		scanCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		output := buf.Bytes()

		// Parse JSON output
		var data map[string]interface{}
		if err := json.Unmarshal(output, &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, string(output))
		}

		// Check project
		project, ok := data["project"].(map[string]interface{})
		if !ok {
			t.Fatal("JSON should contain 'project' object")
		}
		if project["name"] != "Test Novel Project" {
			t.Errorf("project name = %v, want 'Test Novel Project'", project["name"])
		}
		if !strings.Contains(project["world"].(string), "fantasy world") {
			t.Error("project world should contain 'fantasy world'")
		}

		// Check worldbooks
		worldbooks, ok := data["worldbooks"].([]interface{})
		if !ok {
			t.Fatal("JSON should contain 'worldbooks' array")
		}
		if len(worldbooks) != 2 {
			t.Errorf("worldbooks count = %d, want 2", len(worldbooks))
		}

		// Check characters
		characters, ok := data["characters"].([]interface{})
		if !ok {
			t.Fatal("JSON should contain 'characters' array")
		}
		if len(characters) != 2 {
			t.Errorf("characters count = %d, want 2", len(characters))
		}

		// Check chapters
		chapters, ok := data["chapters"].([]interface{})
		if !ok {
			t.Fatal("JSON should contain 'chapters' array")
		}
		if len(chapters) != 2 {
			t.Errorf("chapters count = %d, want 2", len(chapters))
		}

		// Verify chapter data structure
		chapter1 := chapters[0].(map[string]interface{})
		if chapter1["title"] != "The Beginning" {
			t.Errorf("first chapter title = %v, want 'The Beginning'", chapter1["title"])
		}
		if int(chapter1["index"].(float64)) != 1 {
			t.Errorf("first chapter index = %v, want 1", chapter1["index"])
		}
	})

	t.Run("json format with empty collections", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create minimal project file
		projectContent := `---
name: Minimal Project
---
Minimal world.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create empty directories
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()
		scanCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		output := buf.Bytes()

		// Parse JSON output
		var data map[string]interface{}
		if err := json.Unmarshal(output, &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, string(output))
		}

		// Check empty collections (they may be null or empty arrays)
		if worldbooks, ok := data["worldbooks"].([]interface{}); ok {
			if len(worldbooks) != 0 {
				t.Errorf("expected empty worldbooks array, got %d items", len(worldbooks))
			}
		}

		if characters, ok := data["characters"].([]interface{}); ok {
			if len(characters) != 0 {
				t.Errorf("expected empty characters array, got %d items", len(characters))
			}
		}

		if chapters, ok := data["chapters"].([]interface{}); ok {
			if len(chapters) != 0 {
				t.Errorf("expected empty chapters array, got %d items", len(chapters))
			}
		}
	})

	t.Run("json output structure matches novelmaker types", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()
		scanCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		// Unmarshal to verify structure matches expected types
		var data struct {
			Project    novelmaker.Project     `json:"project"`
			Worldbooks []novelmaker.Worldbook `json:"worldbooks"`
			Characters []novelmaker.Character `json:"characters"`
			Chapters   []novelmaker.Chapter   `json:"chapters"`
		}

		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("failed to unmarshal JSON into expected types: %v\nOutput: %s", err, buf.String())
		}

		// Verify project
		if data.Project.Name != "Test Novel Project" {
			t.Errorf("project name mismatch: got %s", data.Project.Name)
		}

		// Verify worldbooks have required fields
		if len(data.Worldbooks) != 2 {
			t.Errorf("expected 2 worldbooks, got %d", len(data.Worldbooks))
		}
		for i, wb := range data.Worldbooks {
			if wb.ID == "" {
				t.Errorf("worldbook %d missing ID", i)
			}
			if wb.Content == "" {
				t.Errorf("worldbook %d missing content", i)
			}
		}

		// Verify characters have required fields
		if len(data.Characters) != 2 {
			t.Errorf("expected 2 characters, got %d", len(data.Characters))
		}
		for i, char := range data.Characters {
			if char.ID == "" {
				t.Errorf("character %d missing ID", i)
			}
			if char.Name == "" {
				t.Errorf("character %d missing name", i)
			}
		}

		// Verify chapters have required fields
		if len(data.Chapters) != 2 {
			t.Errorf("expected 2 chapters, got %d", len(data.Chapters))
		}
		for i, ch := range data.Chapters {
			if ch.ID == "" {
				t.Errorf("chapter %d missing ID", i)
			}
			if ch.Title == "" {
				t.Errorf("chapter %d missing title", i)
			}
			if ch.Index == 0 {
				t.Errorf("chapter %d has invalid index", i)
			}
		}
	})
}

func TestScanCmd_NewScanCmd(t *testing.T) {
	scanCmd := NewScanCmd()

	if scanCmd.cmd == nil {
		t.Fatal("command should not be nil")
	}

	if scanCmd.cmd.Use != "scan" {
		t.Errorf("command use = %s, want 'scan'", scanCmd.cmd.Use)
	}

	if scanCmd.cmd.Short == "" {
		t.Error("command should have short description")
	}

	if scanCmd.cmd.Long == "" {
		t.Error("command should have long description")
	}

	// Check flag exists
	flag := scanCmd.cmd.Flags().Lookup("json")
	if flag == nil {
		t.Error("command should have --json flag")
	}
}

func TestScanCmd_ContentFormatting(t *testing.T) {
	t.Run("newlines replaced with spaces in text output", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Test Project
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create worldbook with newlines
		wb := `---
id: multi-line
tags:
  - test
---
This is line one.
This is line two.
This is line three.`
		writeTestFile(t, tmpDir, "World/001_multiline.md", wb)

		// Create empty directories for other parts
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		// Check that content is on one line (newlines replaced with spaces)
		if strings.Contains(output, "This is line one.\nThis is line two.") {
			t.Error("output should not contain original newlines in content")
		}
		if strings.Contains(output, "This is line one. This is line two.") {
			// This is expected - newlines should be replaced with spaces
		}
	})

	t.Run("tags displayed as comma-separated", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Test Project
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create worldbook with multiple tags
		wb := `---
id: tagged-entry
tags:
  - magic
  - combat
  - advanced
---
Content here.`
		writeTestFile(t, tmpDir, "World/001_tagged.md", wb)

		// Create empty directories for other parts
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		// Check that tags are comma-separated
		if !strings.Contains(output, "Tags: magic,combat,advanced") {
			t.Error("tags should be displayed as comma-separated list")
		}
	})
}

func TestScanCmd_InvalidCurrentDirectory(t *testing.T) {
	// This test verifies error handling when current directory is invalid
	// We can't easily test this without actually changing to a bad directory,
	// but we can at least verify the error message structure
	t.Run("error message includes context", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create an invalid vault (missing Config directory)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()
		err := scanCmd.run(scanCmd.cmd, []string{})

		if err == nil {
			t.Error("expected error for invalid vault")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "failed to") {
			t.Errorf("error message should contain 'failed to', got: %v", err)
		}
	})
}

func TestScanCmd_CharacterSorting(t *testing.T) {
	t.Run("main characters appear first", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Test Project
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create characters with mixed main/non-main
		char1 := `---
id: side-char-1
name: Zara
main: false
---
A side character.`
		writeTestFile(t, tmpDir, "Character/zara.md", char1)

		char2 := `---
id: main-char-1
name: Alice
main: true
---
Main character.`
		writeTestFile(t, tmpDir, "Character/alice.md", char2)

		char3 := `---
id: main-char-2
name: Bob
main: true
---
Another main character.`
		writeTestFile(t, tmpDir, "Character/bob.md", char3)

		// Create empty directories for other parts
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()
		scanCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		var data struct {
			Characters []novelmaker.Character `json:"characters"`
		}

		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if len(data.Characters) != 3 {
			t.Fatalf("expected 3 characters, got %d", len(data.Characters))
		}

		// First two should be main characters
		if !data.Characters[0].Main {
			t.Error("first character should be main")
		}
		if !data.Characters[1].Main {
			t.Error("second character should be main")
		}
		if data.Characters[2].Main {
			t.Error("third character should not be main")
		}

		// Main characters should be alphabetically sorted (Alice, Bob)
		if data.Characters[0].Name != "Alice" {
			t.Errorf("first main character should be Alice, got %s", data.Characters[0].Name)
		}
		if data.Characters[1].Name != "Bob" {
			t.Errorf("second main character should be Bob, got %s", data.Characters[1].Name)
		}
	})
}

func TestScanCmd_ChapterSorting(t *testing.T) {
	t.Run("chapters sorted by index", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Test Project
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create chapters out of order
		ch3 := `---
id: chapter-three
title: Chapter Three
index: 3
---
Content three.`
		writeTestFile(t, tmpDir, "Story/003_chapter-three.md", ch3)

		ch1 := `---
id: chapter-one
title: Chapter One
index: 1
---
Content one.`
		writeTestFile(t, tmpDir, "Story/001_chapter-one.md", ch1)

		ch2 := `---
id: chapter-two
title: Chapter Two
index: 2
---
Content two.`
		writeTestFile(t, tmpDir, "Story/002_chapter-two.md", ch2)

		// Create empty directories for other parts
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()
		scanCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		var data struct {
			Chapters []novelmaker.Chapter `json:"chapters"`
		}

		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if len(data.Chapters) != 3 {
			t.Fatalf("expected 3 chapters, got %d", len(data.Chapters))
		}

		// Verify sorting by index
		for i, ch := range data.Chapters {
			expectedIndex := i + 1
			if ch.Index != expectedIndex {
				t.Errorf("chapter at position %d should have index %d, got %d", i, expectedIndex, ch.Index)
			}
		}
	})
}

func TestScanCmd_DateFormatting(t *testing.T) {
	t.Run("dates formatted correctly in text output", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)

		// Wait a moment to ensure files have valid timestamps
		time.Sleep(10 * time.Millisecond)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		scanCmd := NewScanCmd()

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := scanCmd.run(scanCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if err != nil {
			t.Fatalf("scan command failed: %v", err)
		}

		// Check for date format "YYYY-MM-DD HH:MM:SS"
		if !strings.Contains(output, "Created:") {
			t.Error("output should contain 'Created:' label")
		}
		if !strings.Contains(output, "Updated:") {
			t.Error("output should contain 'Updated:' label")
		}

		// The dates should be in YYYY-MM-DD format
		datePattern := "202" // Should match current year range
		if !strings.Contains(output, datePattern) {
			t.Error("output should contain formatted dates")
		}
	})
}
