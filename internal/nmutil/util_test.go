package nmutil

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed test_embed_files
var testEmbedFS embed.FS

func TestCopyEmbedFS(t *testing.T) {
	t.Run("copy files and directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		err = CopyEmbedFS(testEmbedFS, "test_embed_files", root)
		if err != nil {
			t.Fatalf("copyEmbedFS failed: %v", err)
		}

		// Verify files were copied
		expectedFiles := []string{
			"file1.txt",
			"file2.md",
			"subdir/nested.txt",
			"subdir/deep/deeper.md",
		}

		for _, file := range expectedFiles {
			fullPath := filepath.Join(tmpDir, file)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("expected file %s to exist, but it doesn't", file)
			}
		}
	})

	t.Run("verify file contents", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		err = CopyEmbedFS(testEmbedFS, "test_embed_files", root)
		if err != nil {
			t.Fatalf("copyEmbedFS failed: %v", err)
		}

		// Read and verify content
		content, err := root.ReadFile("file1.txt")
		if err != nil {
			t.Fatalf("failed to read file1.txt: %v", err)
		}

		expectedContent := "This is test file 1"
		if strings.TrimSpace(string(content)) != expectedContent {
			t.Errorf("file1.txt content = %s, want %s", strings.TrimSpace(string(content)), expectedContent)
		}

		// Verify nested file content
		nestedContent, err := root.ReadFile("subdir/nested.txt")
		if err != nil {
			t.Fatalf("failed to read subdir/nested.txt: %v", err)
		}

		expectedNested := "Nested content"
		if strings.TrimSpace(string(nestedContent)) != expectedNested {
			t.Errorf("subdir/nested.txt content = %s, want %s", strings.TrimSpace(string(nestedContent)), expectedNested)
		}
	})

	t.Run("verify directory creation", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		err = CopyEmbedFS(testEmbedFS, "test_embed_files", root)
		if err != nil {
			t.Fatalf("copyEmbedFS failed: %v", err)
		}

		// Verify directories were created
		expectedDirs := []string{
			"subdir",
			"subdir/deep",
		}

		for _, dir := range expectedDirs {
			fullPath := filepath.Join(tmpDir, dir)
			info, err := os.Stat(fullPath)
			if err != nil {
				t.Errorf("expected directory %s to exist: %v", dir, err)
				continue
			}
			if !info.IsDir() {
				t.Errorf("%s should be a directory", dir)
			}
		}
	})

	t.Run("verify file permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		err = CopyEmbedFS(testEmbedFS, "test_embed_files", root)
		if err != nil {
			t.Fatalf("copyEmbedFS failed: %v", err)
		}

		// Verify file has correct permissions (0644)
		filePath := filepath.Join(tmpDir, "file1.txt")
		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}

		expectedPerm := os.FileMode(0644)
		if info.Mode().Perm() != expectedPerm {
			t.Errorf("file permissions = %v, want %v", info.Mode().Perm(), expectedPerm)
		}

		// Verify directory has correct permissions (0755)
		dirPath := filepath.Join(tmpDir, "subdir")
		dirInfo, err := os.Stat(dirPath)
		if err != nil {
			t.Fatalf("failed to stat directory: %v", err)
		}

		expectedDirPerm := os.FileMode(0755)
		if dirInfo.Mode().Perm() != expectedDirPerm {
			t.Errorf("directory permissions = %v, want %v", dirInfo.Mode().Perm(), expectedDirPerm)
		}
	})

	t.Run("empty source directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		// Create a minimal embed.FS for testing empty directory
		var emptyFS embed.FS
		err = CopyEmbedFS(emptyFS, "nonexistent", root)

		// Should handle gracefully - either succeed or return appropriate error
		// We don't enforce specific behavior for empty/nonexistent paths
		if err != nil {
			// Error is acceptable for nonexistent source
			if !strings.Contains(err.Error(), "file does not exist") &&
				!strings.Contains(err.Error(), "no such file") &&
				!strings.Contains(err.Error(), "not found") {
				t.Errorf("unexpected error type: %v", err)
			}
		}
	})

	t.Run("handles existing directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		// Pre-create a directory
		if err := root.MkdirAll("subdir", 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}

		// Copy should still work
		err = CopyEmbedFS(testEmbedFS, "test_embed_files", root)
		if err != nil {
			t.Fatalf("copyEmbedFS failed with existing directory: %v", err)
		}

		// Verify file in existing directory was still created
		content, err := root.ReadFile("subdir/nested.txt")
		if err != nil {
			t.Fatalf("failed to read file in existing directory: %v", err)
		}

		if len(content) == 0 {
			t.Error("file should have content")
		}
	})

	t.Run("multiple files in root", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		err = CopyEmbedFS(testEmbedFS, "test_embed_files", root)
		if err != nil {
			t.Fatalf("copyEmbedFS failed: %v", err)
		}

		// Check that multiple root-level files exist
		rootFiles := []string{"file1.txt", "file2.md"}
		for _, file := range rootFiles {
			if _, err := root.Stat(file); err != nil {
				t.Errorf("root file %s should exist: %v", file, err)
			}
		}
	})

	t.Run("preserves file extension", func(t *testing.T) {
		tmpDir := t.TempDir()
		root, err := os.OpenRoot(tmpDir)
		if err != nil {
			t.Fatalf("failed to open root: %v", err)
		}
		defer root.Close()

		err = CopyEmbedFS(testEmbedFS, "test_embed_files", root)
		if err != nil {
			t.Fatalf("copyEmbedFS failed: %v", err)
		}

		// Verify extensions are preserved
		txtPath := filepath.Join(tmpDir, "file1.txt")
		mdPath := filepath.Join(tmpDir, "file2.md")

		if !strings.HasSuffix(txtPath, ".txt") {
			t.Error("txt file should preserve .txt extension")
		}
		if !strings.HasSuffix(mdPath, ".md") {
			t.Error("md file should preserve .md extension")
		}

		// Verify files actually exist
		if _, err := os.Stat(txtPath); err != nil {
			t.Errorf("file1.txt should exist: %v", err)
		}
		if _, err := os.Stat(mdPath); err != nil {
			t.Errorf("file2.md should exist: %v", err)
		}
	})
}
