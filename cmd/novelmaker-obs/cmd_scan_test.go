package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/cmd/novelmaker-obs/testutil"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

func TestScanCmd_Run_TextOutput(t *testing.T) {
	t.Run("complete vault", func(t *testing.T) {
		tmpDir := testutil.SetupCompleteVault(t)
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
		tmpDir := testutil.CreateTestVault(t)

		// Create minimal project file
		projectContent := `---
name: Empty Project
---
Empty world.`
		testutil.WriteTestFile(t, tmpDir, "Config/project.md", projectContent)

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
		tmpDir := testutil.CreateTestVault(t)
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
		tmpDir := testutil.CreateTestVault(t)

		// Create project file
		projectContent := `---
name: Test Project
---
World description.`
		testutil.WriteTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create worldbook with long content (>80 chars)
		longContent := strings.Repeat("This is a very long worldbook entry that should be truncated in the output. ", 5)
		wb := `---
tags:
  - test
---
` + longContent
		testutil.WriteTestFile(t, tmpDir, "World/001_long.md", wb)

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
		tmpDir := testutil.SetupCompleteVault(t)
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
		var data map[string]any
		if err := json.Unmarshal(output, &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, string(output))
		}

		// Check project
		project, ok := data["project"].(map[string]any)
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
		worldbooks, ok := data["worldbooks"].([]any)
		if !ok {
			t.Fatal("JSON should contain 'worldbooks' array")
		}
		if len(worldbooks) != 2 {
			t.Errorf("worldbooks count = %d, want 2", len(worldbooks))
		}

		// Check characters
		characters, ok := data["characters"].([]any)
		if !ok {
			t.Fatal("JSON should contain 'characters' array")
		}
		if len(characters) != 2 {
			t.Errorf("characters count = %d, want 2", len(characters))
		}

		// Check chapters
		chapters, ok := data["chapters"].([]any)
		if !ok {
			t.Fatal("JSON should contain 'chapters' array")
		}
		if len(chapters) != 2 {
			t.Errorf("chapters count = %d, want 2", len(chapters))
		}

		// Verify chapter data structure
		chapter1 := chapters[0].(map[string]any)
		if chapter1["title"] != "The Beginning" {
			t.Errorf("first chapter title = %v, want 'The Beginning'", chapter1["title"])
		}
		if int(chapter1["index"].(float64)) != 1 {
			t.Errorf("first chapter index = %v, want 1", chapter1["index"])
		}
	})

	t.Run("json format with empty collections", func(t *testing.T) {
		tmpDir := testutil.CreateTestVault(t)

		// Create minimal project file
		projectContent := `---
name: Minimal Project
---
Minimal world.`
		testutil.WriteTestFile(t, tmpDir, "Config/project.md", projectContent)

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
		var data map[string]any
		if err := json.Unmarshal(output, &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, string(output))
		}

		// Check empty collections (they may be null or empty arrays)
		if worldbooks, ok := data["worldbooks"].([]any); ok {
			if len(worldbooks) != 0 {
				t.Errorf("expected empty worldbooks array, got %d items", len(worldbooks))
			}
		}

		if characters, ok := data["characters"].([]any); ok {
			if len(characters) != 0 {
				t.Errorf("expected empty characters array, got %d items", len(characters))
			}
		}

		if chapters, ok := data["chapters"].([]any); ok {
			if len(chapters) != 0 {
				t.Errorf("expected empty chapters array, got %d items", len(chapters))
			}
		}
	})

	t.Run("json output structure matches novelmaker types", func(t *testing.T) {
		tmpDir := testutil.SetupCompleteVault(t)
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
			if wb.Content == "" {
				t.Errorf("worldbook %d missing content", i)
			}
		}

		// Verify characters have required fields
		if len(data.Characters) != 2 {
			t.Errorf("expected 2 characters, got %d", len(data.Characters))
		}
		for i, char := range data.Characters {
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
