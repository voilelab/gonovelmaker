package obsidian

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/voilelab/gonovelmaker/novelmaker"
)

type Vault struct {
	root *os.Root
	path string
}

func NewVault(root string) (*Vault, error) {
	rt, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	return &Vault{root: rt, path: root}, nil
}

func (v *Vault) Close() error {
	return v.root.Close()
}

func (v *Vault) LoadProject() (*novelmaker.Project, error) {
	return loadProjectFromRoot(v.root)
}

func (v *Vault) LoadWorldbooks() ([]novelmaker.Worldbook, error) {
	return loadWorldbooksFromRoot(v.root)
}

func (v *Vault) LoadCharacters() ([]novelmaker.Character, error) {
	return loadCharactersFromRoot(v.root)
}

func (v *Vault) AddCharacter(c *novelmaker.Character) error {
	// Ensure Character directory exists on disk
	charDir := filepath.Join(v.path, "Character")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		return fmt.Errorf("failed to create Character directory: %w", err)
	}

	// Prepare frontmatter similar to CLI behavior
	frontmatter := fmt.Sprintf(`---
id: %s
name: %s
main: %t
---
%s`, c.ID, c.Name, c.Main, c.Profile)

	// Destination file
	filename := fmt.Sprintf("%s.md", c.ID)
	dstPath := filepath.Join(charDir, filename)

	if err := os.WriteFile(dstPath, []byte(frontmatter), 0644); err != nil {
		return fmt.Errorf("failed to write character file %s: %w", dstPath, err)
	}

	return nil

}

func (v *Vault) LoadChapters() ([]novelmaker.Chapter, error) {
	return loadChaptersFromRoot(v.root)
}

func (v *Vault) AddChapter(c *novelmaker.Chapter) error {
	// Ensure Story directory exists on disk
	storyDir := filepath.Join(v.path, "Story")
	if err := os.MkdirAll(storyDir, 0755); err != nil {
		return fmt.Errorf("failed to create Story directory: %w", err)
	}

	// Prepare frontmatter similar to loader expectations
	frontmatter := fmt.Sprintf(`---
id: %s
title: %s
index: %d
---
%s`, c.ID, c.Title, c.Index, c.Content)

	// Destination file: include index for readability
	filename := fmt.Sprintf("%03d_%s.md", c.Index, c.ID)
	dstPath := filepath.Join(storyDir, filename)

	if err := os.WriteFile(dstPath, []byte(frontmatter), 0644); err != nil {
		return fmt.Errorf("failed to write chapter file %s: %w", dstPath, err)
	}

	return nil
}
