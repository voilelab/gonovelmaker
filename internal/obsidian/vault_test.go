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

	t.Run("missing name field", func(t *testing.T) {
		tmpDir := createTestVault(t)
		projectContent := `---
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
tags:
  - magic
  - rules
---
The magic system is based on elements.`

		wb2 := `---
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
			if worldbooks[i].ID == "World/001_magic.md" {
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
name: Alice
main: true
---
Alice is the protagonist.`

		char2 := `---
name: Bob
main: false
---
Bob is a side character.`

		char3 := `---
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
	t.Run("single chapter", func(t *testing.T) {
		tmpDir := createTestVault(t)
		chapterContent := `---
title: Chapter One
index: 1
prompt: test prompt
---
This is the content of chapter one.`
		writeTestFile(t, tmpDir, "Story/001_ch1.md", chapterContent)

		vault, err := NewVault(tmpDir)
		if err != nil {
			t.Fatalf("NewVault failed: %v", err)
		}
		defer vault.Close()

		chapters, err := vault.LoadChapters()
		if err != nil {
			t.Fatalf("LoadChapters failed: %v", err)
		}

		if len(chapters) != 1 {
			t.Fatalf("expected 1 chapter, got %d", len(chapters))
		}

		if chapters[0].Title != "Chapter One" {
			t.Errorf("expected chapter title 'Chapter One', got '%s'", chapters[0].Title)
		}

		if chapters[0].Index != 1 {
			t.Errorf("expected chapter index 1, got %d", chapters[0].Index)
		}

		if chapters[0].Prompt != "test prompt" {
			t.Errorf("expected chapter prompt 'test prompt', got '%s'", chapters[0].Prompt)
		}

		if chapters[0].Content != "This is the content of chapter one." {
			t.Errorf("expected chapter content 'This is the content of chapter one.', got '%s'", chapters[0].Content)
		}
	})

	t.Run("multiple chapters with sorting", func(t *testing.T) {
		tmpDir := createTestVault(t)

		ch1 := `---
title: Prologue
index: 1
---
This is the prologue.`

		ch2 := `---
title: Chapter One
index: 2
---
This is chapter one.`

		ch3 := `---
title: Chapter Two
index: 3
---
This is chapter two.`

		// Write in non-sorted order
		writeTestFile(t, tmpDir, "Story/002_ch2.md", ch2)
		writeTestFile(t, tmpDir, "Story/003_ch3.md", ch3)
		writeTestFile(t, tmpDir, "Story/001_ch1.md", ch1)

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
		if chapters[0].Index != 1 {
			t.Errorf("first chapter should be prologue with index 1, got index %d", chapters[0].Index)
		}
		if chapters[1].Index != 2 {
			t.Errorf("second chapter should be chapter-one with index 2, got index %d", chapters[1].Index)
		}
		if chapters[2].Index != 3 {
			t.Errorf("third chapter should be chapter-two with index 3, got index %d", chapters[2].Index)
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
			Name:      "Test Character",
			Main:      true,
			Profile:   "This is a test character profile.",
			UpdatedAt: time.Now(),
		}

		_, err = vault.AddCharacter(char)
		if err != nil {
			t.Fatalf("AddCharacter failed: %v", err)
		}

		// Verify file was created
		filePath := filepath.Join(tmpDir, "Character", "test_character.md")
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read created file: %v", err)
		}

		contentStr := string(content)
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
			Name:    "New Character",
			Main:    false,
			Profile: "Profile content.",
		}

		_, err = vault.AddCharacter(char)
		if err != nil {
			t.Fatalf("AddCharacter failed: %v", err)
		}

		// Verify directory and file were created
		dirPath := filepath.Join(tmpDir, "Character")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Error("Character directory should be created")
		}

		filePath := filepath.Join(dirPath, "new_character.md")
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
			Index:     5,
			Title:     "Test Chapter",
			Prompt:    "This is the chapter prompt.",
			Content:   "This is the chapter content.",
			UpdatedAt: time.Now(),
		}

		_, err = vault.AddChapter(chapter)
		if err != nil {
			t.Fatalf("AddChapter failed: %v", err)
		}

		// Verify file was created with correct name format
		filePath := filepath.Join(tmpDir, "Story", "005_ch5.md")
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read created file: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "title: Test Chapter") {
			t.Error("file should contain 'title: Test Chapter'")
		}
		if !strings.Contains(contentStr, "index: 5") {
			t.Error("file should contain 'index: 5'")
		}
		if !strings.Contains(contentStr, "This is the chapter prompt.") {
			t.Error("file should contain chapter prompt")
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
			Index:   1,
			Title:   "New Chapter",
			Content: "Content.",
		}

		_, err = vault.AddChapter(chapter)
		if err != nil {
			t.Fatalf("AddChapter failed: %v", err)
		}

		// Verify directory and file were created
		dirPath := filepath.Join(tmpDir, "Story")
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Error("Story directory should be created")
		}

		filePath := filepath.Join(dirPath, "001_ch1.md")
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

		storySamplePath := filepath.Join(tmpDir, "Story", "001_ch1.md")
		if _, err := os.Stat(storySamplePath); os.IsNotExist(err) {
			t.Error("Story/001_ch1.md should be created")
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
