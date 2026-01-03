package cmdchapter

import (
	"os"
	"path/filepath"
	"testing"
)

// createTestVault creates a temporary vault directory for testing
func createTestVault(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return tmpDir
}

// writeTestFile writes a file to the test vault
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

// setupCompleteVault creates a complete test vault with all required files
func setupCompleteVault(t *testing.T) string {
	t.Helper()
	tmpDir := createTestVault(t)

	// Create a .novelmaker config directory and config.toml file
	configDir := filepath.Join(tmpDir, ".novelmaker")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}
	configContent := `user_llm_backend = "test"

[llm_backend.test]
type = "openai"
api_key = "test-key"
model = "gpt-4o"
image_model = "dall-e-3"
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Set HOME environment variable to tmpDir so config.Load() finds our test config
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})

	// Create project file
	projectContent := `---
name: Test Novel Project
system_prompt: You are a helpful AI assistant for writing novels.
system_prompt_char: You are a character development AI assistant.
---
This is a fantasy world with magic and dragons.`
	writeTestFile(t, tmpDir, "Config/project.md", projectContent)

	// Create chapter prompt template
	chapterPromptContent := `---
system: |
    You are a professional novel writing assistant.
---
`
	writeTestFile(t, tmpDir, "Config/chapter_prompt.md", chapterPromptContent)

	// Create character prompt template
	characterPromptContent := `---
system: You are a character development AI assistant.
---
Create a detailed character profile.`
	writeTestFile(t, tmpDir, "Config/character_prompt.md", characterPromptContent)

	// Create worldbook entries
	wb1 := `---
tags:
  - magic
  - rules
---
Magic in this world is based on elemental manipulation.`
	writeTestFile(t, tmpDir, "World/001_magic-system.md", wb1)

	wb2 := `---
tags:
  - world
  - locations
---
The world consists of three continents separated by vast oceans.`
	writeTestFile(t, tmpDir, "World/002_geography.md", wb2)

	// Create characters
	char1 := `---
name: Alice
main: true
prompt: A brave knight
---
Alice is a brave knight who seeks to protect her kingdom.`
	writeTestFile(t, tmpDir, "Character/alice.md", char1)

	char2 := `---
name: Bob
main: false
prompt: A wise wizard
---
Bob is a wise old wizard who mentors the protagonist.`
	writeTestFile(t, tmpDir, "Character/bob.md", char2)

	// Create chapters
	ch1 := `---
title: The Beginning
index: 1
prompt: Write an engaging prologue that introduces the main character.
---
It was a dark and stormy night when Alice first discovered her destiny.`
	writeTestFile(t, tmpDir, "Story/001_prologue.md", ch1)

	ch2 := `---
title: Chapter One - The Journey Begins
index: 2
prompt: Continue the story as Alice leaves her village.
---
At dawn, Alice packed her belongings and set out on her journey.`
	writeTestFile(t, tmpDir, "Story/002_chapter-one.md", ch2)

	return tmpDir
}
