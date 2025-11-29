package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
tags:
  - magic
  - rules
---
Magic in this world is based on elemental manipulation.`
	writeTestFile(t, tmpDir, "World/001_magic-system.md", wb1)

	wb2 := `---
tags:
  - world
  - locations
---
The world consists of three continents separated by vast oceans.`
	writeTestFile(t, tmpDir, "World/002_geography.md", wb2)

	// Create characters
	char1 := `---
name: Alice
main: true
---
Alice is a brave knight who seeks to protect her kingdom.`
	writeTestFile(t, tmpDir, "Character/alice.md", char1)

	char2 := `---
name: Bob
main: false
---
Bob is a wise old wizard who mentors the protagonist.`
	writeTestFile(t, tmpDir, "Character/bob.md", char2)

	// Create chapters
	ch1 := `---
title: The Beginning
index: 1
prompt: Write an engaging prologue that introduces the main character.
---
It was a dark and stormy night when Alice first discovered her destiny.`
	writeTestFile(t, tmpDir, "Story/001_prologue.md", ch1)

	ch2 := `---
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
		if !strings.Contains(output, "ID: World/001_magic-system.md") {
			t.Error("output should contain magic-system worldbook")
		}
		if !strings.Contains(output, "ID: World/002_geography.md") {
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
			if ch.Title == "" {
				t.Errorf("chapter %d missing title", i)
			}
			if ch.Index == 0 {
				t.Errorf("chapter %d has invalid index", i)
			}
		}
	})
}
