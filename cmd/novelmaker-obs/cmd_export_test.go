package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/cmd/novelmaker-obs/testutil"
)

func TestExportCmd_Run_Success(t *testing.T) {
	t.Run("export to txt format", func(t *testing.T) {
		tmpDir := testutil.SetupCompleteVault(t)
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
		tmpDir := testutil.SetupCompleteVault(t)
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
		tmpDir := testutil.SetupCompleteVault(t)
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
		tmpDir := testutil.SetupCompleteVault(t)
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
