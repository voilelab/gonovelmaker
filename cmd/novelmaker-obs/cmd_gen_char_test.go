package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

func TestGenCharCmd_Run_Success(t *testing.T) {
	t.Run("generate character with dummy backend", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Charlie"
		genCharCmd.prompt = "A wise old wizard with a mysterious past."

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify the new character was created
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		if len(characters) != 3 {
			t.Fatalf("expected 3 characters after gen-char, got %d", len(characters))
		}

		// Find the newly created character
		var newChar *novelmaker.Character
		for i, ch := range characters {
			if ch.Name == "Charlie" {
				newChar = &characters[i]
				break
			}
		}

		if newChar == nil {
			t.Fatal("newly created character not found")
		}

		// Verify character properties
		if newChar.ID != "charlie" {
			t.Errorf("new character ID = %s, want 'charlie'", newChar.ID)
		}
		if newChar.Name != "Charlie" {
			t.Errorf("new character name = %s, want 'Charlie'", newChar.Name)
		}
		if newChar.Main {
			t.Error("new character should not be main by default")
		}
		if !strings.Contains(newChar.Profile, "dummy response") {
			t.Error("new character should contain dummy backend response")
		}
	})

	t.Run("generate character without prompt", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Diana"
		genCharCmd.prompt = "" // Empty prompt

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify the character was created
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		if len(characters) != 3 {
			t.Fatalf("expected 3 characters, got %d", len(characters))
		}

		// Find the newly created character
		var newChar *novelmaker.Character
		for i, ch := range characters {
			if ch.Name == "Diana" {
				newChar = &characters[i]
				break
			}
		}

		if newChar == nil {
			t.Fatal("newly created character not found")
		}

		if newChar.ID != "diana" {
			t.Errorf("new character ID = %s, want 'diana'", newChar.ID)
		}
	})

	t.Run("generate character with custom prompt", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Emma"
		genCharCmd.prompt = "A fierce warrior queen with unmatched battle skills."

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify the character was created with the profile
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		// Find the newly created character
		var newChar *novelmaker.Character
		for i, ch := range characters {
			if ch.Name == "Emma" {
				newChar = &characters[i]
				break
			}
		}

		if newChar == nil {
			t.Fatal("newly created character not found")
		}

		// The profile should come from the dummy backend
		if !strings.Contains(newChar.Profile, "dummy response") {
			t.Error("character profile should contain dummy backend response")
		}
	})

	t.Run("generate character with special characters in name", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Sir John O'Brien III"
		genCharCmd.prompt = "A noble knight"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify the character was created with slugified ID
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		// Find the newly created character
		var newChar *novelmaker.Character
		for i, ch := range characters {
			if ch.Name == "Sir John O'Brien III" {
				newChar = &characters[i]
				break
			}
		}

		if newChar == nil {
			t.Fatal("newly created character not found")
		}

		// ID should be slugified
		if newChar.ID == "" {
			t.Error("character ID should not be empty")
		}
		// ID should not contain special characters
		if strings.ContainsAny(newChar.ID, " '") {
			t.Errorf("character ID should not contain special characters, got: %s", newChar.ID)
		}
	})
}

func TestGenCharCmd_Run_JSONOutput(t *testing.T) {
	t.Run("json output format", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Frank"
		genCharCmd.prompt = "A mysterious stranger"
		genCharCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := genCharCmd.run(genCharCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Parse JSON output
		var buf []byte
		buf, _ = io.ReadAll(r)
		output := string(buf)

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		// Verify JSON contains filepath
		filepath, ok := data["filepath"].(string)
		if !ok {
			t.Fatal("JSON output should contain 'filepath' field")
		}

		if filepath == "" {
			t.Error("filepath should not be empty")
		}

		// Verify the file exists
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Errorf("file should exist at path: %s", filepath)
		}
	})
}

func TestGenCharCmd_Run_ErrorCases(t *testing.T) {
	t.Run("error when project not found", func(t *testing.T) {
		tmpDir := createTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "TestChar"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when project not found, got nil")
		}

		if !strings.Contains(err.Error(), "failed to load project") {
			t.Errorf("error should mention 'failed to load project', got: %v", err)
		}
	})

	t.Run("error when config cannot be loaded", func(t *testing.T) {
		tmpDir := createTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		// Change to a directory that doesn't have a config
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "TestChar"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when config cannot be loaded")
		}
	})
}

func TestGenCharCmd_Run_WithContext(t *testing.T) {
	t.Run("includes worldbook entries in generation", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Add additional worldbook entries
		wb3 := `---
id: character-lore
tags:
  - lore
---
This world has a rich history of heroic characters.`
		writeTestFile(t, tmpDir, "World/003_character-lore.md", wb3)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Gregory"
		genCharCmd.prompt = "A character based on the world lore"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify character was created (worldbook should be included in context)
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		if len(characters) != 3 {
			t.Fatalf("expected 3 characters, got %d", len(characters))
		}
	})

	t.Run("includes existing characters in generation", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Hannah"
		genCharCmd.prompt = "A character that interacts with existing characters"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify character was created (existing characters should be included in context)
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		if len(characters) != 3 {
			t.Fatalf("expected 3 characters, got %d", len(characters))
		}
	})

	t.Run("works with no existing characters", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create minimal project
		projectContent := `---
name: New Project
system_prompt: You are a writer.
system_prompt_char: You are a character development AI.
---
A new world.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create empty directories
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Ian"
		genCharCmd.prompt = "First character in the story"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify the character was created
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		if len(characters) != 1 {
			t.Fatalf("expected 1 character, got %d", len(characters))
		}

		if characters[0].Name != "Ian" {
			t.Errorf("character name = %s, want 'Ian'", characters[0].Name)
		}
	})

	t.Run("works with no worldbooks", func(t *testing.T) {
		tmpDir := createTestVault(t)

		// Create minimal project
		projectContent := `---
name: Simple Project
system_prompt: You are a writer.
system_prompt_char: You are a character development AI.
---
A simple world.`
		writeTestFile(t, tmpDir, "Config/project.md", projectContent)

		// Create empty directories
		os.MkdirAll(filepath.Join(tmpDir, "World"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Character"), 0755)
		os.MkdirAll(filepath.Join(tmpDir, "Story"), 0755)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Julia"
		genCharCmd.prompt = "A character without worldbook context"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Verify the character was created
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		if len(characters) != 1 {
			t.Fatalf("expected 1 character, got %d", len(characters))
		}
	})
}

func TestGenCharCmd_Run_FlagOverrides(t *testing.T) {
	t.Run("api-key flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Kevin"
		genCharCmd.apiKey = "custom-api-key"

		// Should not fail even with custom API key (dummy backend doesn't use it)
		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}
	})

	t.Run("model flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Laura"
		genCharCmd.model = "custom-model"

		// Should not fail even with custom model (dummy backend doesn't use it)
		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}
	})

	t.Run("base-url flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Michael"
		genCharCmd.baseURL = "https://custom.api.url"

		// Should not fail even with custom base URL (dummy backend doesn't use it)
		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}
	})

	t.Run("timeout flag overrides config", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Nancy"
		genCharCmd.timeout = 60

		// Should not fail with custom timeout
		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}
	})
}

func TestGenCharCmd_Run_FileCreation(t *testing.T) {
	t.Run("created file has valid format", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Oliver"
		genCharCmd.prompt = "Test character prompt"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Find the created file
		charDir := filepath.Join(tmpDir, "Character")
		newFile := filepath.Join(charDir, "oliver.md")

		// Verify the file exists
		if _, err := os.Stat(newFile); os.IsNotExist(err) {
			// List all files in Character to help debug
			entries, _ := os.ReadDir(charDir)
			t.Logf("Files in Character directory:")
			for _, e := range entries {
				t.Logf("  %s", e.Name())
			}
			t.Fatal("new character file not found at expected location")
		}

		// Read and verify file content
		content, err := os.ReadFile(newFile)
		if err != nil {
			t.Fatalf("failed to read new file: %v", err)
		}

		contentStr := string(content)

		// Verify frontmatter
		if !strings.Contains(contentStr, "---") {
			t.Error("file should have frontmatter delimiters")
		}
		if !strings.Contains(contentStr, "id:") {
			t.Error("file should have id field")
		}
		if !strings.Contains(contentStr, "name:") {
			t.Error("file should have name field")
		}
		if !strings.Contains(contentStr, "main:") {
			t.Error("file should have main field")
		}

		// Verify profile from dummy backend
		if !strings.Contains(contentStr, "dummy response") {
			t.Error("file should contain dummy backend response")
		}
	})

	t.Run("filename is properly slugified", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Princess Anne Marie"
		genCharCmd.prompt = "A royal character"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		// Check that file was created with slugified name
		charDir := filepath.Join(tmpDir, "Character")
		entries, err := os.ReadDir(charDir)
		if err != nil {
			t.Fatalf("failed to read Character directory: %v", err)
		}

		found := false
		for _, entry := range entries {
			name := entry.Name()
			// Check if the file was created
			if strings.Contains(strings.ToLower(name), "princess") {
				found = true
				// Filename should not contain spaces
				if strings.Contains(name, " ") {
					t.Errorf("filename should not contain spaces, got: %s", name)
				}
			}
		}

		if !found {
			t.Error("character file should be created")
		}
	})

	t.Run("character has correct default values", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Quinn"
		genCharCmd.prompt = "A mysterious figure"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		// Find the newly created character
		var newChar *novelmaker.Character
		for i, ch := range characters {
			if ch.Name == "Quinn" {
				newChar = &characters[i]
				break
			}
		}

		if newChar == nil {
			t.Fatal("newly created character not found")
		}

		// Verify default values
		if newChar.Main {
			t.Error("character should not be main by default")
		}
		if newChar.ID == "" {
			t.Error("character should have an ID")
		}
		if newChar.Profile == "" {
			t.Error("character should have a profile")
		}
	})
}

func TestGenCharCmd_Run_CharacterID(t *testing.T) {
	t.Run("character ID is slugified from name", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "Rachel Smith"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		// Find the newly created character
		var newChar *novelmaker.Character
		for i, ch := range characters {
			if ch.Name == "Rachel Smith" {
				newChar = &characters[i]
				break
			}
		}

		if newChar == nil {
			t.Fatal("newly created character not found")
		}

		// ID should be slugified
		if strings.Contains(newChar.ID, " ") {
			t.Errorf("character ID should not contain spaces, got: %s", newChar.ID)
		}
		if strings.Contains(newChar.ID, "RACHEL") {
			t.Error("character ID should be lowercase")
		}
	})

	t.Run("character ID handles unicode characters", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCmd := NewGenCharCmd(dummyBackendMaker)
		genCharCmd.name = "José García"

		err := genCharCmd.run(genCharCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char command failed: %v", err)
		}

		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		// Find the newly created character
		var newChar *novelmaker.Character
		for i, ch := range characters {
			if ch.Name == "José García" {
				newChar = &characters[i]
				break
			}
		}

		if newChar == nil {
			t.Fatal("newly created character not found")
		}

		// ID should be valid
		if newChar.ID == "" {
			t.Error("character ID should not be empty")
		}
	})
}

func TestGenCharCmd_Run_MultipleCharacters(t *testing.T) {
	t.Run("create multiple characters in sequence", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Create first character
		genCharCmd1 := NewGenCharCmd(dummyBackendMaker)
		genCharCmd1.name = "Sarah"
		err := genCharCmd1.run(genCharCmd1.cmd, []string{})
		if err != nil {
			t.Fatalf("first gen-char command failed: %v", err)
		}

		// Create second character
		genCharCmd2 := NewGenCharCmd(dummyBackendMaker)
		genCharCmd2.name = "Tom"
		err = genCharCmd2.run(genCharCmd2.cmd, []string{})
		if err != nil {
			t.Fatalf("second gen-char command failed: %v", err)
		}

		// Verify both characters exist
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}
		defer vault.Close()

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		// Should have original 2 + new 2 = 4 characters
		if len(characters) != 4 {
			t.Fatalf("expected 4 characters, got %d", len(characters))
		}

		// Verify both new characters exist
		foundSarah := false
		foundTom := false
		for _, ch := range characters {
			if ch.Name == "Sarah" {
				foundSarah = true
			}
			if ch.Name == "Tom" {
				foundTom = true
			}
		}

		if !foundSarah {
			t.Error("Sarah character not found")
		}
		if !foundTom {
			t.Error("Tom character not found")
		}
	})
}
