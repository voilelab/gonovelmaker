package obsidian

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/voilelab/gonovelmaker/novelmaker"
)

func createTestVault(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return tmpDir
}

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

func TestNewVault(t *testing.T) {
	t.Run("valid directory", func(t *testing.T) {
		tmpDir := createTestVault(t)
		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		if vault.root == nil {
			t.Error("vault.root is nil")
		}
	})

	t.Run("invalid directory", func(t *testing.T) {
		_, err := NewVault("/nonexistent/path/that/does/not/exist")
		if err == nil {
			t.Error("expected error for nonexistent directory, got nil")
		}
	})
}

func TestVault_LoadProject(t *testing.T) {
	t.Run("valid project", func(t *testing.T) {
		tmpDir := createTestVault(t)
		projectContent := `---
id: test-project
name: Test Project
system_prompt: Test system prompt
system_prompt_char: Test character prompt
---
This is the world description.`

		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		project, err := vault.LoadProject()
		if err != nil {
			t.Fatalf("LoadProject failed: %v", err)
		}

		if project.ID != "test-project" {
			t.Errorf("project.ID = %s, want test-project", project.ID)
		}
		if project.Name != "Test Project" {
			t.Errorf("project.Name = %s, want Test Project", project.Name)
		}
		if project.World != "This is the world description." {
			t.Errorf("project.World = %s, want 'This is the world description.'", project.World)
		}
		if project.SystemPrompt != "Test system prompt" {
			t.Errorf("project.SystemPrompt = %s, want 'Test system prompt'", project.SystemPrompt)
		}
		if project.SystemPromptChar != "Test character prompt" {
			t.Errorf("project.SystemPromptChar = %s, want 'Test character prompt'", project.SystemPromptChar)
		}
	})

	t.Run("missing project file", func(t *testing.T) {
		tmpDir := createTestVault(t)
		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		_, err = vault.LoadProject()
		if err == nil {
			t.Error("expected error for missing project file, got nil")
		}
	})

	t.Run("missing id field", func(t *testing.T) {
		tmpDir := createTestVault(t)
		projectContent := `---
name: Test Project
---
World content`

		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		_, err = vault.LoadProject()
		if err == nil {
			t.Error("expected error for missing id field, got nil")
		}
	})

	t.Run("missing name field", func(t *testing.T) {
		tmpDir := createTestVault(t)
		projectContent := `---
id: test-project
---
World content`

		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		_, err = vault.LoadProject()
		if err == nil {
			t.Error("expected error for missing name field, got nil")
		}
	})

	t.Run("invalid frontmatter", func(t *testing.T) {
		tmpDir := createTestVault(t)
		projectContent := `This file has no frontmatter`

		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		_, err = vault.LoadProject()
		if err == nil {
			t.Error("expected error for invalid frontmatter, got nil")
		}
	})
}

func TestVault_LoadWorldbooks(t *testing.T) {
	t.Run("multiple worldbooks", func(t *testing.T) {
		tmpDir := createTestVault(t)

		wb1 := `---
id: magic-system
tags:
  - magic
  - rules
---
The magic system is based on elements.`

		wb2 := `---
id: geography
tags:
  - world
  - locations
---
The world consists of three continents.`

		writeTestFile(t, tmpDir, "World/001_magic.md", wb1)
		writeTestFile(t, tmpDir, "World/002_geography.md", wb2)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		worldbooks, err := vault.LoadWorldbooks()
		if err != nil {
			t.Fatalf("LoadWorldbooks failed: %v", err)
		}

		if len(worldbooks) != 2 {
			t.Fatalf("expected 2 worldbooks, got %d", len(worldbooks))
		}

		// Find magic-system worldbook
		var magicWB *novelmaker.Worldbook
		for i := range worldbooks {
			if worldbooks[i].ID == "magic-system" {
				magicWB = &worldbooks[i]
				break
			}
		}

		if magicWB == nil {
			t.Fatal("magic-system worldbook not found")
		}

		if magicWB.Content != "The magic system is based on elements." {
			t.Errorf("magicWB.Content = %s, want 'The magic system is based on elements.'", magicWB.Content)
		}

		if len(magicWB.Tags) != 2 {
			t.Errorf("magicWB.Tags length = %d, want 2", len(magicWB.Tags))
		}
	})

	t.Run("empty worldbook directory", func(t *testing.T) {
		tmpDir := createTestVault(t)
		if err := os.MkdirAll(filepath.Join(tmpDir, "World"), 0755); err != nil {
			t.Fatalf("failed to create World directory: %v", err)
		}

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		worldbooks, err := vault.LoadWorldbooks()
		if err != nil {
			t.Fatalf("LoadWorldbooks failed: %v", err)
		}

		if len(worldbooks) != 0 {
			t.Errorf("expected 0 worldbooks, got %d", len(worldbooks))
		}
	})

	t.Run("skip non-markdown files", func(t *testing.T) {
		tmpDir := createTestVault(t)

		wb := `---
id: magic-system
tags:
  - magic
---
Magic content.`

		writeTestFile(t, tmpDir, "World/magic.md", wb)
		writeTestFile(t, tmpDir, "World/notes.txt", "Some notes")
		writeTestFile(t, tmpDir, "World/README.pdf", "PDF content")

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		worldbooks, err := vault.LoadWorldbooks()
		if err != nil {
			t.Fatalf("LoadWorldbooks failed: %v", err)
		}

		if len(worldbooks) != 1 {
			t.Errorf("expected 1 worldbook, got %d", len(worldbooks))
		}
	})

	t.Run("missing worldbook directory", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		_, err = vault.LoadWorldbooks()
		if err == nil {
			t.Error("expected error for missing World directory, got nil")
		}
	})
}

func TestVault_LoadCharacters(t *testing.T) {
	t.Run("multiple characters with sorting", func(t *testing.T) {
		tmpDir := createTestVault(t)

		char1 := `---
id: alice
name: Alice
main: true
---
Alice is the protagonist.`

		char2 := `---
id: bob
name: Bob
main: false
---
Bob is a side character.`

		char3 := `---
id: charlie
name: Charlie
main: true
---
Charlie is a main character.`

		writeTestFile(t, tmpDir, "Character/alice.md", char1)
		writeTestFile(t, tmpDir, "Character/bob.md", char2)
		writeTestFile(t, tmpDir, "Character/charlie.md", char3)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("LoadCharacters failed: %v", err)
		}

		if len(characters) != 3 {
			t.Fatalf("expected 3 characters, got %d", len(characters))
		}

		// Check sorting: main characters first, then alphabetically
		if !characters[0].Main {
			t.Error("first character should be main")
		}
		if !characters[1].Main {
			t.Error("second character should be main")
		}
		if characters[2].Main {
			t.Error("third character should not be main")
		}

		// Main characters should be sorted alphabetically
		if characters[0].Name == "Charlie" && characters[1].Name != "Alice" {
			t.Error("main characters should be sorted alphabetically")
		}
	})

	t.Run("empty character directory", func(t *testing.T) {
		tmpDir := createTestVault(t)
		if err := os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755); err != nil {
			t.Fatalf("failed to create Character directory: %v", err)
		}

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("LoadCharacters failed: %v", err)
		}

		if len(characters) != 0 {
			t.Errorf("expected 0 characters, got %d", len(characters))
		}
	})

	t.Run("missing character directory", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("LoadCharacters failed: %v", err)
		}

		if len(characters) != 0 {
			t.Errorf("expected 0 characters for missing directory, got %d", len(characters))
		}
	})
}

func TestVault_LoadChapters(t *testing.T) {
	t.Run("multiple chapters with sorting", func(t *testing.T) {
		tmpDir := createTestVault(t)

		ch1 := `---
id: prologue
title: Prologue
index: 1
---
This is the prologue.`

		ch2 := `---
id: chapter-one
title: Chapter One
index: 2
---
This is chapter one.`

		ch3 := `---
id: chapter-two
title: Chapter Two
index: 3
---
This is chapter two.`

		// Write in non-sorted order
		writeTestFile(t, tmpDir, "Story/002_chapter-one.md", ch2)
		writeTestFile(t, tmpDir, "Story/003_chapter-two.md", ch3)
		writeTestFile(t, tmpDir, "Story/001_prologue.md", ch1)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("LoadChapters failed: %v", err)
		}

		if len(chapters) != 3 {
			t.Fatalf("expected 3 chapters, got %d", len(chapters))
		}

		// Check sorting by index
		if chapters[0].Index != 1 || chapters[0].ID != "prologue" {
			t.Errorf("first chapter should be prologue with index 1, got %s with index %d", chapters[0].ID, chapters[0].Index)
		}
		if chapters[1].Index != 2 || chapters[1].ID != "chapter-one" {
			t.Errorf("second chapter should be chapter-one with index 2, got %s with index %d", chapters[1].ID, chapters[1].Index)
		}
		if chapters[2].Index != 3 || chapters[2].ID != "chapter-two" {
			t.Errorf("third chapter should be chapter-two with index 3, got %s with index %d", chapters[2].ID, chapters[2].Index)
		}
	})

	t.Run("missing story directory", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		_, err = vault.LoadChapters()
		if err == nil {
			t.Error("expected error for missing Story directory, got nil")
		}
	})
}

func TestVault_AddCharacter(t *testing.T) {
	t.Run("add new character", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		char := &novelmaker.Character{
			ID:        "test-char",
			Name:      "Test Character",
			Main:      true,
			Profile:   "This is a test character profile.",
			UpdatedAt: time.Now(),
		}

		err = vault.AddCharacter(char)
		if err != nil {
			t.Fatalf("AddCharacter failed: %v", err)
		}

		// Verify file was created
		filePath := filepath.Join(tmpDir, "Character", "test-char.md")
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read created file: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "id: test-char") {
			t.Error("file should contain 'id: test-char'")
		}
		if !strings.Contains(contentStr, "name: Test Character") {
			t.Error("file should contain 'name: Test Character'")
		}
		if !strings.Contains(contentStr, "main: true") {
			t.Error("file should contain 'main: true'")
		}
		if !strings.Contains(contentStr, "This is a test character profile.") {
			t.Error("file should contain character profile")
		}
	})

	t.Run("creates character directory if missing", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		char := &novelmaker.Character{
			ID:      "new-char",
			Name:    "New Character",
			Main:    false,
			Profile: "Profile content.",
		}

		err = vault.AddCharacter(char)
		if err != nil {
			t.Fatalf("AddCharacter failed: %v", err)
		}

		// Verify directory and file were created
		dirPath := filepath.Join(tmpDir, "Character")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Error("Character directory should be created")
		}

		filePath := filepath.Join(dirPath, "new-char.md")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Error("character file should be created")
		}
	})
}

func TestVault_AddChapter(t *testing.T) {
	t.Run("add new chapter", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		chapter := &novelmaker.Chapter{
			ID:        "test-chapter",
			Index:     5,
			Title:     "Test Chapter",
			Content:   "This is the chapter content.",
			UpdatedAt: time.Now(),
		}

		err = vault.AddChapter(chapter)
		if err != nil {
			t.Fatalf("AddChapter failed: %v", err)
		}

		// Verify file was created with correct name format
		filePath := filepath.Join(tmpDir, "Story", "005_test-chapter.md")
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read created file: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "id: test-chapter") {
			t.Error("file should contain 'id: test-chapter'")
		}
		if !strings.Contains(contentStr, "title: Test Chapter") {
			t.Error("file should contain 'title: Test Chapter'")
		}
		if !strings.Contains(contentStr, "index: 5") {
			t.Error("file should contain 'index: 5'")
		}
		if !strings.Contains(contentStr, "This is the chapter content.") {
			t.Error("file should contain chapter content")
		}
	})

	t.Run("creates story directory if missing", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		chapter := &novelmaker.Chapter{
			ID:      "new-chapter",
			Index:   1,
			Title:   "New Chapter",
			Content: "Content.",
		}

		err = vault.AddChapter(chapter)
		if err != nil {
			t.Fatalf("AddChapter failed: %v", err)
		}

		// Verify directory and file were created
		dirPath := filepath.Join(tmpDir, "Story")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Error("Story directory should be created")
		}

		filePath := filepath.Join(dirPath, "001_new-chapter.md")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Error("chapter file should be created")
		}
	})
}

func TestVault_Initialize(t *testing.T) {
	t.Run("initialize new vault", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		err = vault.Initialize()
		if err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}

		// Verify Config directory was created
		configPath := filepath.Join(tmpDir, "Config")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config directory should be created")
		}

		// Verify project.md exists in Config
		projectPath := filepath.Join(tmpDir, "Config", "project.md")
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			t.Error("Config/project.md should be created")
		}

		// Verify other directories were created
		worldPath := filepath.Join(tmpDir, "World")
		if _, err := os.Stat(worldPath); os.IsNotExist(err) {
			t.Error("World directory should be created")
		}

		charPath := filepath.Join(tmpDir, "Character")
		if _, err := os.Stat(charPath); os.IsNotExist(err) {
			t.Error("Character directory should be created")
		}

		storyPath := filepath.Join(tmpDir, "Story")
		if _, err := os.Stat(storyPath); os.IsNotExist(err) {
			t.Error("Story directory should be created")
		}

		// Verify sample files exist
		charSamplePath := filepath.Join(tmpDir, "Character", "character_sample.md")
		if _, err := os.Stat(charSamplePath); os.IsNotExist(err) {
			t.Error("Character/character_sample.md should be created")
		}

		worldSamplePath := filepath.Join(tmpDir, "World", "001_world_sample.md")
		if _, err := os.Stat(worldSamplePath); os.IsNotExist(err) {
			t.Error("World/001_world_sample.md should be created")
		}

		storySamplePath := filepath.Join(tmpDir, "Story", "001_prologue.md")
		if _, err := os.Stat(storySamplePath); os.IsNotExist(err) {
			t.Error("Story/001_prologue.md should be created")
		}

		// Verify README exists
		readmePath := filepath.Join(tmpDir, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			t.Error("README.md should be created")
		}

		// Verify project.md has valid frontmatter
		content, err := os.ReadFile(projectPath)
		if err != nil {
			t.Fatalf("failed to read project.md: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "---") {
			t.Error("project.md should have frontmatter")
		}
		if !strings.Contains(contentStr, "id:") {
			t.Error("project.md should have id field")
		}
		if !strings.Contains(contentStr, "name:") {
			t.Error("project.md should have name field")
		}
	})

	t.Run("error when Config already exists", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create Config directory beforehand
		configPath := filepath.Join(tmpDir, "Config")
		if err := os.MkdirAll(configPath, 0755); err != nil {
			t.Fatalf("failed to create Config directory: %v", err)
		}

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		err = vault.Initialize()
		if err == nil {
			t.Error("expected error when Config directory already exists, got nil")
		}

		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("error should mention 'already exists', got: %v", err)
		}
	})

	t.Run("initialize creates readable files", func(t *testing.T) {
		tmpDir := createTestVault(t)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		err = vault.Initialize()
		if err != nil {
			t.Fatalf("Initialize failed: %v", err)
		}

		// Try to load the project after initialization
		project, err := vault.LoadProject()
		if err != nil {
			t.Fatalf("LoadProject after Initialize failed: %v", err)
		}

		if project.ID == "" {
			t.Error("initialized project should have an ID")
		}
		if project.Name == "" {
			t.Error("initialized project should have a name")
		}

		// Try to load worldbooks
		worldbooks, err := vault.LoadWorldbooks()
		if err != nil {
			t.Fatalf("LoadWorldbooks after Initialize failed: %v", err)
		}

		if len(worldbooks) == 0 {
			t.Error("expected at least one sample worldbook after initialization")
		}

		// Try to load characters
		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("LoadCharacters after Initialize failed: %v", err)
		}

		if len(characters) == 0 {
			t.Error("expected at least one sample character after initialization")
		}

		// Try to load chapters
		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("LoadChapters after Initialize failed: %v", err)
		}

		if len(chapters) == 0 {
			t.Error("expected at least one sample chapter after initialization")
		}
	})
}

func TestVault_Close(t *testing.T) {
	tmpDir := createTestVault(t)
	vault, err := NewVault(tmpDir)
	if err != nil {
		t.Fatalf("NewVault failed: %v", err)
	}

	err = vault.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
