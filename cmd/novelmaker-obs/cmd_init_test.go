package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

func TestInitCmd_Run_Success(t *testing.T) {
	t.Run("initialize with plugin", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		initCmd.includePlugin = true

		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		// Verify Config directory was created
		configPath := filepath.Join(tmpDir, "Config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config directory should be created")
		}

		// Verify project.md exists
		projectPath := filepath.Join(tmpDir, "Config", "project.md")
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			t.Error("Config/project.md should be created")
		}

		// Verify World directory was created
		worldPath := filepath.Join(tmpDir, "World")
		if _, err := os.Stat(worldPath); os.IsNotExist(err) {
			t.Error("World directory should be created")
		}

		// Verify Character directory was created
		charPath := filepath.Join(tmpDir, "Character")
		if _, err := os.Stat(charPath); os.IsNotExist(err) {
			t.Error("Character directory should be created")
		}

		// Verify Story directory was created
		storyPath := filepath.Join(tmpDir, "Story")
		if _, err := os.Stat(storyPath); os.IsNotExist(err) {
			t.Error("Story directory should be created")
		}

		// Verify plugin files were created
		pluginPath := filepath.Join(tmpDir, ".obsidian", "plugins", "obsidian-novelmaker")
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			t.Error("plugin directory should be created")
		}

		mainJsPath := filepath.Join(pluginPath, "main.js")
		if _, err := os.Stat(mainJsPath); os.IsNotExist(err) {
			t.Error("plugin main.js should be created")
		}

		manifestPath := filepath.Join(pluginPath, "manifest.json")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			t.Error("plugin manifest.json should be created")
		}

		stylesPath := filepath.Join(pluginPath, "styles.css")
		if _, err := os.Stat(stylesPath); os.IsNotExist(err) {
			t.Error("plugin styles.css should be created")
		}
	})

	t.Run("initialize without plugin", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		initCmd.includePlugin = false

		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		// Verify basic directories were created
		configPath := filepath.Join(tmpDir, "Config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config directory should be created")
		}

		// Verify plugin files were NOT created
		pluginPath := filepath.Join(tmpDir, ".obsidian", "plugins", "obsidian-novelmaker")
		if _, err := os.Stat(pluginPath); !os.IsNotExist(err) {
			t.Error("plugin directory should not be created when include-plugin is false")
		}
	})
}

func TestInitCmd_Run_CreatesValidFiles(t *testing.T) {
	t.Run("created project file is valid", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		// Verify project.md has valid frontmatter
		projectPath := filepath.Join(tmpDir, "Config", "project.md")
		content, err := os.ReadFile(projectPath)
		if err != nil {
			t.Fatalf("failed to read project.md: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "---") {
			t.Error("project.md should have frontmatter delimiters")
		}
		if !strings.Contains(contentStr, "name:") {
			t.Error("project.md should have name field")
		}

		// Try to load the project to verify it's valid
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to create vault: %v", err)
		}
		defer vault.Close()

		project, err := vault.LoadProject()
		if err != nil {
			t.Fatalf("failed to load initialized project: %v", err)
		}

		if project.Name == "" {
			t.Error("initialized project should have a name")
		}
	})

	t.Run("created sample files are valid", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to create vault: %v", err)
		}
		defer vault.Close()

		// Try to load worldbooks
		worldbooks, err := vault.LoadWorldbooks()
		if err != nil {
			t.Fatalf("failed to load worldbooks: %v", err)
		}
		if len(worldbooks) == 0 {
			t.Error("expected at least one sample worldbook")
		}

		// Try to load characters
		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}
		if len(characters) == 0 {
			t.Error("expected at least one sample character")
		}

		// Try to load chapters
		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("failed to load chapters: %v", err)
		}
		if len(chapters) == 0 {
			t.Error("expected at least one sample chapter")
		}

		// Verify sample chapter has required fields
		if chapters[0].ID == "" {
			t.Error("sample chapter should have ID")
		}
		if chapters[0].Title == "" {
			t.Error("sample chapter should have title")
		}
		if chapters[0].Index == 0 {
			t.Error("sample chapter should have valid index")
		}
	})

	t.Run("README is created", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		readmePath := filepath.Join(tmpDir, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			t.Error("README.md should be created")
		}

		content, err := os.ReadFile(readmePath)
		if err != nil {
			t.Fatalf("failed to read README.md: %v", err)
		}

		if len(content) == 0 {
			t.Error("README.md should not be empty")
		}
	})
}

func TestInitCmd_Run_ErrorCases(t *testing.T) {
	t.Run("error when Config already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Pre-create Config directory
		configPath := filepath.Join(tmpDir, "Config")
		if err := os.MkdirAll(configPath, 0755); err != nil {
			t.Fatalf("failed to create Config directory: %v", err)
		}

		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when Config directory already exists, got nil")
		}

		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("error should mention 'already exists', got: %v", err)
		}
	})

	t.Run("error when project.md already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Pre-create Config directory with project.md
		configPath := filepath.Join(tmpDir, "Config")
		if err := os.MkdirAll(configPath, 0755); err != nil {
			t.Fatalf("failed to create Config directory: %v", err)
		}
		projectPath := filepath.Join(configPath, "project.md")
		if err := os.WriteFile(projectPath, []byte("existing content"), 0644); err != nil {
			t.Fatalf("failed to create project.md: %v", err)
		}

		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when project already initialized, got nil")
		}
	})

	t.Run("error handling invalid directory", func(t *testing.T) {
		// Save current directory
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		// Make sure we're in a valid directory
		os.Chdir(oldWd)

		// Create the init command but don't run it from a valid temp dir
		initCmd := NewInitCmd()

		// This test verifies the command structure, actual directory errors
		// are handled by the vault initialization
		if initCmd.cmd.RunE == nil {
			t.Error("command should have RunE function")
		}
	})
}
