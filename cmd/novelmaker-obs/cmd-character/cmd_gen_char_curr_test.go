package cmdcharacter

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/internal/llmbackend"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

func TestGenCharCurrCmd_Run_Success(t *testing.T) {
	t.Run("regenerate character with dummy backend", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Load the vault to verify initial state
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}

		// Find Alice's file path (assuming she exists in the test vault)
		var aliceFilePath string
		for _, ch := range characters {
			if ch.Name == "Alice" {
				aliceFilePath = "Character/alice.md"
				break
			}
		}

		if aliceFilePath == "" {
			t.Skip("Alice character not found in test vault, skipping test")
		}

		// Store original profile to compare later
		originalChar, err := vault.LoadCharacterByPath(aliceFilePath)
		if err != nil {
			t.Fatalf("failed to load Alice character: %v", err)
		}
		originalProfile := originalChar.Profile

		vault.Close()

		// Run gen-char-curr command
		genCharCurrCmd := NewGenCharCurrCmd(llmbackend.MakeDummy)
		genCharCurrCmd.filepath = aliceFilePath

		err = genCharCurrCmd.run(genCharCurrCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char-curr command failed: %v", err)
		}

		// Verify the character was regenerated
		vault2, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to reopen vault: %v", err)
		}
		defer vault2.Close()

		updatedChar, err := vault2.LoadCharacterByPath(aliceFilePath)
		if err != nil {
			t.Fatalf("failed to load updated character: %v", err)
		}

		// Verify the profile was updated
		if updatedChar.Profile == originalProfile {
			t.Error("character profile should have been regenerated")
		}

		// Verify it contains dummy backend response
		if !strings.Contains(updatedChar.Profile, "dummy response") {
			t.Error("regenerated character should contain dummy backend response")
		}

		// Verify name and other metadata are preserved
		if updatedChar.Name != originalChar.Name {
			t.Errorf("character name should be preserved, got %s, want %s", updatedChar.Name, originalChar.Name)
		}
		if updatedChar.Main != originalChar.Main {
			t.Errorf("character main status should be preserved, got %v, want %v", updatedChar.Main, originalChar.Main)
		}
		if updatedChar.Prompt != originalChar.Prompt {
			t.Errorf("character prompt should be preserved, got %s, want %s", updatedChar.Prompt, originalChar.Prompt)
		}
	})
}

func TestGenCharCurrCmd_Run_JSONOutput(t *testing.T) {
	t.Run("json output format", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Find a character file
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		characters, err := vault.LoadCharacters()
		if err != nil {
			t.Fatalf("failed to load characters: %v", err)
		}
		vault.Close()

		if len(characters) == 0 {
			t.Skip("no characters in test vault, skipping test")
		}

		// Use the first character
		testFilePath := "Character/alice.md" // Assuming alice exists

		genCharCurrCmd := NewGenCharCurrCmd(llmbackend.MakeDummy)
		genCharCurrCmd.filepath = testFilePath
		genCharCurrCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err = genCharCurrCmd.run(genCharCurrCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("gen-char-curr command failed: %v", err)
		}

		// Parse JSON output
		var buf []byte
		buf, _ = io.ReadAll(r)
		output := string(buf)

		var data map[string]any
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		// Verify JSON contains expected fields
		filepath, ok := data["filepath"].(string)
		if !ok {
			t.Fatal("JSON output should contain 'filepath' field")
		}

		if filepath != testFilePath {
			t.Errorf("filepath = %s, want %s", filepath, testFilePath)
		}

		name, ok := data["name"].(string)
		if !ok {
			t.Fatal("JSON output should contain 'name' field")
		}

		if name == "" {
			t.Error("name should not be empty")
		}
	})
}

func TestGenCharCurrCmd_Run_ErrorCases(t *testing.T) {
	t.Run("error when character file not found", func(t *testing.T) {
		tmpDir := setupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCurrCmd := NewGenCharCurrCmd(llmbackend.MakeDummy)
		genCharCurrCmd.filepath = "Character/nonexistent.md"

		err := genCharCurrCmd.run(genCharCurrCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when character file not found, got nil")
		}

		if !strings.Contains(err.Error(), "failed to load target character") {
			t.Errorf("error should mention 'failed to load target character', got: %v", err)
		}
	})

	t.Run("error when config cannot be loaded", func(t *testing.T) {
		tmpDir := createTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharCurrCmd := NewGenCharCurrCmd(llmbackend.MakeDummy)
		genCharCurrCmd.filepath = "Character/test.md"

		err := genCharCurrCmd.run(genCharCurrCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when config cannot be loaded")
		}
	})

	t.Run("requires filepath flag", func(t *testing.T) {
		cmd := NewGenCharCurrCmd(llmbackend.MakeDummy)
		cmd.cmd.SetArgs([]string{})
		err := cmd.cmd.Execute()
		if err == nil {
			t.Fatal("expected error when filepath flag is missing")
		}
	})

	t.Run("accepts valid filepath", func(t *testing.T) {
		cmd := NewGenCharCurrCmd(llmbackend.MakeDummy)
		if err := cmd.cmd.ParseFlags([]string{"--filepath", "Character/alice.md"}); err != nil {
			t.Fatalf("failed to parse valid flags: %v", err)
		}
	})
}
