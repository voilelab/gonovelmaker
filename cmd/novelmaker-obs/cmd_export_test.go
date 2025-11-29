package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

func TestExportCmd_Run_Success(t *testing.T) {
	t.Run("export to txt format", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Verify the output file was created
		if _, err := os.Stat(exportCmd.outFile); os.IsNotExist(err) {
			t.Error("output file should be created")
		}

		// Read and verify the output file content
		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Verify project name appears in output
		if !strings.Contains(contentStr, "Test Novel Project") {
			t.Error("output should contain project name")
		}

		// Verify chapters appear in output
		if !strings.Contains(contentStr, "Chapter 1: The Beginning") {
			t.Error("output should contain first chapter title")
		}
		if !strings.Contains(contentStr, "Chapter 2: Chapter One - The Journey Begins") {
			t.Error("output should contain second chapter title")
		}

		// Verify chapter content appears
		if !strings.Contains(contentStr, "dark and stormy night") {
			t.Error("output should contain first chapter content")
		}
		if !strings.Contains(contentStr, "Alice packed her belongings") {
			t.Error("output should contain second chapter content")
		}
	})

	t.Run("export with absolute path", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		outputFile := filepath.Join(tmpDir, "novel-export.txt")

		exportCmd := NewExportCmd()
		exportCmd.outFile = outputFile
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Verify the output file exists at the absolute path
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Error("output file should be created at specified path")
		}
	})

	t.Run("export with relative path", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = "novel.txt"
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Verify the output file was created in current directory
		outputPath := filepath.Join(tmpDir, "novel.txt")
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Error("output file should be created in current directory")
		}
	})

	t.Run("export with nested directory path", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Create subdirectory for output
		exportDir := filepath.Join(tmpDir, "exports")
		os.MkdirAll(exportDir, 0755)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(exportDir, "novel.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Verify the output file was created
		if _, err := os.Stat(exportCmd.outFile); os.IsNotExist(err) {
			t.Error("output file should be created in subdirectory")
		}
	})
}

func TestExportCmd_Run_ChapterOrdering(t *testing.T) {
	t.Run("chapters are exported in correct order", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Ordering Test
system_prompt: You are a writer.
system_prompt_char: You are a character development AI.
---
Test world.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create chapters out of order in filesystem
		ch3 := `---
id: chapter-three
title: Third Chapter
index: 3
---
Content of chapter three.`
		writeTestFile(t, tmpDir, "Story/003_chapter-three.md", ch3)

		ch1 := `---
id: chapter-one
title: First Chapter
index: 1
---
Content of chapter one.`
		writeTestFile(t, tmpDir, "Story/001_chapter-one.md", ch1)

		ch2 := `---
id: chapter-two
title: Second Chapter
index: 2
---
Content of chapter two.`
		writeTestFile(t, tmpDir, "Story/002_chapter-two.md", ch2)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Read output file
		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Find positions of chapters
		pos1 := strings.Index(contentStr, "Content of chapter one")
		pos2 := strings.Index(contentStr, "Content of chapter two")
		pos3 := strings.Index(contentStr, "Content of chapter three")

		if pos1 == -1 || pos2 == -1 || pos3 == -1 {
			t.Fatal("all chapters should appear in output")
		}

		// Verify chapters appear in correct order
		if pos1 > pos2 {
			t.Error("chapter one should appear before chapter two")
		}
		if pos2 > pos3 {
			t.Error("chapter two should appear before chapter three")
		}
	})

	t.Run("chapter numbers match index", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Verify chapter headers have correct numbers
		if !strings.Contains(contentStr, "Chapter 1:") {
			t.Error("first chapter should be numbered as 1")
		}
		if !strings.Contains(contentStr, "Chapter 2:") {
			t.Error("second chapter should be numbered as 2")
		}
	})
}

func TestExportCmd_Run_EmptyProject(t *testing.T) {
	t.Run("export with no chapters", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Empty Novel
system_prompt: You are a writer.
---
A world without chapters yet.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create empty directories
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command should succeed with no chapters: %v", err)
		}

		// Read output file
		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Should contain project name
		if !strings.Contains(contentStr, "Empty Novel") {
			t.Error("output should contain project name even with no chapters")
		}

		// Should not contain any chapter content
		if strings.Contains(contentStr, "Chapter ") {
			t.Error("output should not contain chapter markers when no chapters exist")
		}
	})

	t.Run("export with single chapter", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Single Chapter Novel
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		ch1 := `---
id: only-chapter
title: The Only Chapter
index: 1
---
This is the only chapter content.`
		writeTestFile(t, tmpDir, "Story/001_only.md", ch1)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Verify project and chapter
		if !strings.Contains(contentStr, "Single Chapter Novel") {
			t.Error("output should contain project name")
		}
		if !strings.Contains(contentStr, "Chapter 1: The Only Chapter") {
			t.Error("output should contain chapter title")
		}
		if !strings.Contains(contentStr, "only chapter content") {
			t.Error("output should contain chapter content")
		}
	})
}

func TestExportCmd_Run_ContentFormatting(t *testing.T) {
	t.Run("chapters separated by blank lines", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Check for blank line separation (double newlines)
		if !strings.Contains(contentStr, "\n\n") {
			t.Error("chapters should be separated by blank lines")
		}
	})

	t.Run("project name appears at top", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")

		// Project name should be in the first few lines
		found := false
		for i := 0; i < 5 && i < len(lines); i++ {
			if strings.Contains(lines[i], "Test Novel Project") {
				found = true
				break
			}
		}

		if !found {
			t.Error("project name should appear at the top of the file")
		}
	})

	t.Run("chapter format includes title and content", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Chapter 1 format
		if !strings.Contains(contentStr, "Chapter 1: The Beginning") {
			t.Error("chapter should have 'Chapter N: Title' format")
		}

		// Content should follow title
		chapterIndex := strings.Index(contentStr, "Chapter 1: The Beginning")
		contentIndex := strings.Index(contentStr, "dark and stormy night")

		if chapterIndex == -1 || contentIndex == -1 {
			t.Fatal("chapter header and content should both be present")
		}

		if contentIndex < chapterIndex {
			t.Error("chapter content should appear after chapter title")
		}
	})
}

func TestExportCmd_Run_ErrorCases(t *testing.T) {
	t.Run("error with unsupported file type", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.pdf")
		exportCmd.fileType = "pdf"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error for unsupported file type, got nil")
		}

		if !strings.Contains(err.Error(), "unsupported export file type") {
			t.Errorf("error should mention 'unsupported export file type', got: %v", err)
		}
		if !strings.Contains(err.Error(), "pdf") {
			t.Errorf("error should mention the unsupported type 'pdf', got: %v", err)
		}
	})

	t.Run("error when project not found", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create directories but no project file
		os.MkdirAll(filepath.Join(tmpDir, "Config"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when project not found, got nil")
		}

		if !strings.Contains(err.Error(), "failed to load project") {
			t.Errorf("error should mention 'failed to load project', got: %v", err)
		}
	})

	t.Run("error when vault cannot be created", func(t *testing.T) {
		// Create a temporary directory and then remove it
		tmpDir := t.TempDir()
		invalidPath := filepath.Join(tmpDir, "nonexistent", "vault")

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		// Don't change to invalid directory, just use it in the command
		os.Chdir(invalidPath)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when vault cannot be opened")
		}
	})

	t.Run("error when output directory does not exist", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Try to write to a non-existent directory without creating it
		nonexistentDir := filepath.Join(tmpDir, "nonexistent", "nested", "dir")
		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(nonexistentDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when output directory does not exist")
		}

		if !strings.Contains(err.Error(), "failed to create output file") {
			t.Errorf("error should mention 'failed to create output file', got: %v", err)
		}
	})

	t.Run("invalid chapters are skipped with warning", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Test Project
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create Story directory with invalid chapter file
		invalidChapter := `invalid frontmatter
this is not valid YAML`
		writeTestFile(t, tmpDir, "Story/001_invalid.md", invalidChapter)

		// Create a valid chapter too
		validChapter := `---
id: valid
title: Valid Chapter
index: 1
---
Valid content.`
		writeTestFile(t, tmpDir, "Story/002_valid.md", validChapter)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		// Export should succeed even with invalid chapter (it's just skipped)
		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export should succeed with invalid chapters skipped: %v", err)
		}

		// Verify output file contains valid chapter
		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "Valid content") {
			t.Error("output should contain valid chapter content")
		}
	})
}

func TestExportCmd_Run_FileHandling(t *testing.T) {
	t.Run("overwrites existing file", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		outputPath := filepath.Join(tmpDir, "output.txt")

		// Create existing file with content
		existingContent := "This is old content that should be overwritten."
		err := os.WriteFile(outputPath, []byte(existingContent), 0644)
		if err != nil {
			t.Fatalf("failed to create existing file: %v", err)
		}

		exportCmd := NewExportCmd()
		exportCmd.outFile = outputPath
		exportCmd.fileType = "txt"

		err = exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Read the file and verify it was overwritten
		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Should not contain old content
		if strings.Contains(contentStr, "old content that should be overwritten") {
			t.Error("output file should have been overwritten")
		}

		// Should contain new content
		if !strings.Contains(contentStr, "Test Novel Project") {
			t.Error("output file should contain new export content")
		}
	})

	t.Run("creates file with appropriate permissions", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		outputPath := filepath.Join(tmpDir, "output.txt")

		exportCmd := NewExportCmd()
		exportCmd.outFile = outputPath
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Check file permissions
		info, err := os.Stat(outputPath)
		if err != nil {
			t.Fatalf("failed to stat output file: %v", err)
		}

		// File should be readable and writable by owner
		mode := info.Mode()
		if mode&0600 != 0600 {
			t.Errorf("output file should have read/write permissions for owner, got %v", mode)
		}

		// Should be a regular file
		if !mode.IsRegular() {
			t.Error("output should be a regular file")
		}
	})

	t.Run("file is not empty", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		outputPath := filepath.Join(tmpDir, "output.txt")

		exportCmd := NewExportCmd()
		exportCmd.outFile = outputPath
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Check file size
		info, err := os.Stat(outputPath)
		if err != nil {
			t.Fatalf("failed to stat output file: %v", err)
		}

		if info.Size() == 0 {
			t.Error("output file should not be empty")
		}
	})
}

func TestExportCmd_Run_LargeProjects(t *testing.T) {
	t.Run("export project with many chapters", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Large Novel Project
---
A novel with many chapters.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create 10 chapters
		for i := 1; i <= 10; i++ {
			filename := filepath.Join("Story", fmt.Sprintf("%03d_chapter-%d.md", i, i))
			chapterContent := fmt.Sprintf(`---
id: chapter-%d
title: Chapter Number %d
index: %d
---
This is the content of chapter number %d.`, i, i, i, i)
			writeTestFile(t, tmpDir, filename, chapterContent)
		}

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed with many chapters: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Verify all chapters appear
		for i := 1; i <= 10; i++ {
			chapterMarker := fmt.Sprintf("chapter number %d", i)
			if !strings.Contains(contentStr, chapterMarker) {
				t.Errorf("output should contain chapter %d", i)
			}
		}
	})

	t.Run("export project with long chapter content", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Long Content Project
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create chapter with long content (simulate real novel length)
		longContent := strings.Repeat("This is a paragraph with many words to simulate real novel content. ", 100)
		ch1 := `---
id: long-chapter
title: A Very Long Chapter
index: 1
---
` + longContent
		writeTestFile(t, tmpDir, "Story/001_long.md", ch1)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Verify long content is preserved (check for a portion of it)
		if !strings.Contains(contentStr, "This is a paragraph with many words to simulate real novel content.") {
			t.Error("output should contain chapter content")
		}

		// Verify the content appears multiple times (since we repeated it 100 times)
		count := strings.Count(contentStr, "This is a paragraph with many words")
		if count < 50 {
			t.Errorf("output should contain most of the repeated content, found only %d occurrences", count)
		}

		// Verify output file size is reasonable for the content
		info, _ := os.Stat(exportCmd.outFile)
		if info.Size() < 5000 {
			t.Errorf("output file should be large for long content, got %d bytes", info.Size())
		}
	})
}

func TestExportCmd_Run_SpecialCharacters(t *testing.T) {
	t.Run("handles unicode in chapter content", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: Unicode Novel 小说
---
World with unicode.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		ch1 := `---
id: unicode-chapter
title: Chapter with 日本語
index: 1
---
This chapter contains unicode: ñoño, 你好世界, émojis 🚀🎉, and symbols: ∑∫∂.`
		writeTestFile(t, tmpDir, "Story/001_unicode.md", ch1)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Verify unicode is preserved
		if !strings.Contains(contentStr, "小说") {
			t.Error("output should contain Chinese characters")
		}
		if !strings.Contains(contentStr, "日本語") {
			t.Error("output should contain Japanese characters")
		}
		if !strings.Contains(contentStr, "你好世界") {
			t.Error("output should contain Chinese text")
		}
		if !strings.Contains(contentStr, "🚀") {
			t.Error("output should contain emoji")
		}
		if !strings.Contains(contentStr, "∑") {
			t.Error("output should contain mathematical symbols")
		}
	})

	t.Run("handles special characters in project name", func(t *testing.T) {
		tmpDir := createTestVault(t)

		projectContent := `---
name: "Project: The \"Special\" Edition & More!"
---
World description.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		ch1 := `---
id: simple
title: Simple Chapter
index: 1
---
Content.`
		writeTestFile(t, tmpDir, "Story/001_simple.md", ch1)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		content, err := os.ReadFile(exportCmd.outFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		contentStr := string(content)

		// Verify special characters in project name are preserved
		if !strings.Contains(contentStr, "Project:") {
			t.Error("output should contain project name with special characters")
		}
	})
}

func TestExportCmd_Run_VaultLifecycle(t *testing.T) {
	t.Run("vault is properly closed", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		exportCmd := NewExportCmd()
		exportCmd.outFile = filepath.Join(tmpDir, "output.txt")
		exportCmd.fileType = "txt"

		err := exportCmd.run(exportCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}

		// Verify we can still open the vault after export (vault was closed properly)
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("should be able to open vault after export: %v", err)
		}
		vault.Close()
	})

	t.Run("can export multiple times", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// First export
		exportCmd1 := NewExportCmd()
		exportCmd1.outFile = filepath.Join(tmpDir, "output1.txt")
		exportCmd1.fileType = "txt"

		err := exportCmd1.run(exportCmd1.cmd, []string{})
		if err != nil {
			t.Fatalf("first export failed: %v", err)
		}

		// Second export
		exportCmd2 := NewExportCmd()
		exportCmd2.outFile = filepath.Join(tmpDir, "output2.txt")
		exportCmd2.fileType = "txt"

		err = exportCmd2.run(exportCmd2.cmd, []string{})
		if err != nil {
			t.Fatalf("second export failed: %v", err)
		}

		// Verify both files exist
		if _, err := os.Stat(exportCmd1.outFile); os.IsNotExist(err) {
			t.Error("first export file should exist")
		}
		if _, err := os.Stat(exportCmd2.outFile); os.IsNotExist(err) {
			t.Error("second export file should exist")
		}

		// Verify both files have same content
		content1, _ := os.ReadFile(exportCmd1.outFile)
		content2, _ := os.ReadFile(exportCmd2.outFile)

		if string(content1) != string(content2) {
			t.Error("multiple exports should produce identical content")
		}
	})
}
