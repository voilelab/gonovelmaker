package obsidian

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// copyEmbedFS recursively copies files from an embedded FS to a destination directory
func copyEmbedFS(embedFS embed.FS, srcRoot string, dstRoot *os.Root) error {
	return fs.WalkDir(embedFS, srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from srcRoot
		relPath, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		if d.IsDir() {
			return dstRoot.MkdirAll(relPath, 0755)
		}

		// Read file content from embedded FS
		content, err := fs.ReadFile(embedFS, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Ensure parent directory exists
		if err := dstRoot.MkdirAll(filepath.Dir(relPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", relPath, err)
		}

		// Write file to destination
		if err := dstRoot.WriteFile(relPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", relPath, err)
		}

		return nil
	})
}
