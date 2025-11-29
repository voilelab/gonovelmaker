package main

import (
	"encoding/json"
	"io"
	"os"
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
		if newChar.ID != "Character/charlie.md" {
			t.Errorf("new character ID = %s, want 'Character/charlie.md'", newChar.ID)
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

		if newChar.ID != "Character/diana.md" {
			t.Errorf("new character ID = %s, want 'Character/diana.md'", newChar.ID)
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
