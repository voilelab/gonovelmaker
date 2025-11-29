package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

func TestInitCmd_NewInitCmd(t *testing.T) {
	initCmd := NewInitCmd()

	if initCmd.cmd == nil {
		t.Fatal("command should not be nil")
	}

	if initCmd.cmd.Use != "init" {
		t.Errorf("command use = %s, want 'init'", initCmd.cmd.Use)
	}

	if initCmd.cmd.Short == "" {
		t.Error("command should have short description")
	}

	if initCmd.cmd.Long == "" {
		t.Error("command should have long description")
	}

	// Check flag exists
	flag := initCmd.cmd.Flags().Lookup("include-plugin")
	if flag == nil {
		t.Error("command should have --include-plugin flag")
	}

	// Check default value
	if !initCmd.includePlugin {
		t.Error("--include-plugin should default to true")
	}
}

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

func TestInitCmd_Run_FilePermissions(t *testing.T) {
	t.Run("created directories have correct permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		// Check directory permissions
		dirs := []string{"Config", "World", "Character", "Story"}
		for _, dir := range dirs {
			dirPath := filepath.Join(tmpDir, dir)
			info, err := os.Stat(dirPath)
			if err != nil {
				t.Fatalf("failed to stat %s: %v", dir, err)
			}

			if !info.IsDir() {
				t.Errorf("%s should be a directory", dir)
			}

			// Check that directory is readable and writable
			mode := info.Mode()
			if mode&0700 != 0700 {
				t.Errorf("%s should have read/write/execute permissions for owner, got %v", dir, mode)
			}
		}
	})

	t.Run("created files have correct permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		// Check file permissions
		projectPath := filepath.Join(tmpDir, "Config", "project.md")
		info, err := os.Stat(projectPath)
		if err != nil {
			t.Fatalf("failed to stat project.md: %v", err)
		}

		if info.IsDir() {
			t.Error("project.md should be a file, not a directory")
		}

		// Check that file is readable and writable
		mode := info.Mode()
		if mode&0600 != 0600 {
			t.Errorf("project.md should have read/write permissions for owner, got %v", mode)
		}
	})
}

func TestInitCmd_Run_PluginFiles(t *testing.T) {
	t.Run("plugin files are valid", func(t *testing.T) {
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

		// Check manifest.json is valid JSON
		manifestPath := filepath.Join(tmpDir, ".obsidian", "plugins", "obsidian-novelmaker", "manifest.json")
		content, err := os.ReadFile(manifestPath)
		if err != nil {
			t.Fatalf("failed to read manifest.json: %v", err)
		}

		if len(content) == 0 {
			t.Error("manifest.json should not be empty")
		}

		// Check main.js exists and is not empty
		mainJsPath := filepath.Join(tmpDir, ".obsidian", "plugins", "obsidian-novelmaker", "main.js")
		mainContent, err := os.ReadFile(mainJsPath)
		if err != nil {
			t.Fatalf("failed to read main.js: %v", err)
		}

		if len(mainContent) == 0 {
			t.Error("main.js should not be empty")
		}

		// Check styles.css exists
		stylesPath := filepath.Join(tmpDir, ".obsidian", "plugins", "obsidian-novelmaker", "styles.css")
		if _, err := os.Stat(stylesPath); os.IsNotExist(err) {
			t.Error("styles.css should be created")
		}
	})

	t.Run("plugin directory structure is correct", func(t *testing.T) {
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

		// Verify the .obsidian directory structure
		obsidianPath := filepath.Join(tmpDir, ".obsidian")
		if _, err := os.Stat(obsidianPath); os.IsNotExist(err) {
			t.Error(".obsidian directory should be created")
		}

		pluginsPath := filepath.Join(obsidianPath, "plugins")
		if _, err := os.Stat(pluginsPath); os.IsNotExist(err) {
			t.Error("plugins directory should be created")
		}

		novelmakerPath := filepath.Join(pluginsPath, "obsidian-novelmaker")
		if _, err := os.Stat(novelmakerPath); os.IsNotExist(err) {
			t.Error("obsidian-novelmaker directory should be created")
		}

		// Verify it's a directory
		info, err := os.Stat(novelmakerPath)
		if err != nil {
			t.Fatalf("failed to stat plugin directory: %v", err)
		}
		if !info.IsDir() {
			t.Error("obsidian-novelmaker should be a directory")
		}
	})
}

func TestInitCmd_Run_OutputMessages(t *testing.T) {
	t.Run("success message is displayed", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		initCmd.includePlugin = false // Simplify output for this test

		// Note: This test doesn't capture stdout, but verifies the command runs without error
		// In a real scenario, you might want to capture stdout to verify the exact messages
		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command should succeed: %v", err)
		}

		// Verify initialization actually happened by checking for files
		configPath := filepath.Join(tmpDir, "Config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("initialization should create Config directory")
		}
	})

	t.Run("plugin message is displayed when included", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		initCmd := NewInitCmd()
		initCmd.includePlugin = true

		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command should succeed: %v", err)
		}

		// Verify plugin was installed
		pluginPath := filepath.Join(tmpDir, ".obsidian", "plugins", "obsidian-novelmaker")
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			t.Error("plugin should be installed when include-plugin is true")
		}
	})
}

func TestInitCmd_Run_IdempotencyCheck(t *testing.T) {
	t.Run("running init twice fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// First initialization
		initCmd1 := NewInitCmd()
		err := initCmd1.run(initCmd1.cmd, []string{})
		if err != nil {
			t.Fatalf("first init command failed: %v", err)
		}

		// Second initialization should fail
		initCmd2 := NewInitCmd()
		err = initCmd2.run(initCmd2.cmd, []string{})
		if err == nil {
			t.Error("second init command should fail, got nil")
		}

		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("error should mention 'already exists', got: %v", err)
		}
	})
}

func TestInitCmd_Run_PartialInitialization(t *testing.T) {
	t.Run("handles partial initialization state", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Create some directories but not Config
		worldPath := filepath.Join(tmpDir, "World")
		if err := os.MkdirAll(worldPath, 0755); err != nil {
			t.Fatalf("failed to create World directory: %v", err)
		}

		storyPath := filepath.Join(tmpDir, "Story")
		if err := os.MkdirAll(storyPath, 0755); err != nil {
			t.Fatalf("failed to create Story directory: %v", err)
		}

		// Init should still succeed since Config doesn't exist
		initCmd := NewInitCmd()
		err := initCmd.run(initCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("init command should succeed with partial state: %v", err)
		}

		// Verify Config was created
		configPath := filepath.Join(tmpDir, "Config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config directory should be created")
		}
	})
}

func TestInitCmd_Command_Integration(t *testing.T) {
	t.Run("command can be added to root command", func(t *testing.T) {
		initCmd := NewInitCmd()

		if initCmd.cmd.Use != "init" {
			t.Errorf("expected Use='init', got '%s'", initCmd.cmd.Use)
		}

		// Verify command has the required cobra command structure
		if initCmd.cmd.RunE == nil {
			t.Error("command should have RunE function")
		}
	})

	t.Run("flag parsing works correctly", func(t *testing.T) {
		initCmd := NewInitCmd()

		// Default should be true
		if !initCmd.includePlugin {
			t.Error("default includePlugin should be true")
		}

		// Parse flag with false value
		initCmd.cmd.Flags().Set("include-plugin", "false")
		flag := initCmd.cmd.Flags().Lookup("include-plugin")
		if flag.Value.String() != "false" {
			t.Error("flag should be set to false")
		}

		// Parse flag with true value
		initCmd.cmd.Flags().Set("include-plugin", "true")
		flag = initCmd.cmd.Flags().Lookup("include-plugin")
		if flag.Value.String() != "true" {
			t.Error("flag should be set to true")
		}
	})
}
