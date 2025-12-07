package main

import (
	"testing"

	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

func TestConfigTestCmd(t *testing.T) {
	t.Run("successfully parses valid templates", func(t *testing.T) {
		// Use the init_template directory which has valid templates
		tmpDir := "../../internal/obsidian/init_template"

		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to create vault: %v", err)
		}

		// Test chapter prompt
		chapterPrompt, err := vault.LoadChapterPrompt()
		if err != nil {
			t.Errorf("failed to load chapter prompt: %v", err)
		}
		if chapterPrompt == nil {
			t.Error("chapter prompt should not be nil")
		}
		if chapterPrompt != nil && chapterPrompt.System == "" {
			t.Error("chapter prompt system should not be empty")
		}

		// Test character prompt
		characterPrompt, err := vault.LoadCharacterPrompt()
		if err != nil {
			t.Errorf("failed to load character prompt: %v", err)
		}
		if characterPrompt == nil {
			t.Error("character prompt should not be nil")
		}
		if characterPrompt != nil && characterPrompt.System == "" {
			t.Error("character prompt system should not be empty")
		}
	})
}
