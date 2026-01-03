package cmdcharacter

import (
	"testing"

	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

func TestCharacterCmd_Structure(t *testing.T) {
	t.Run("has correct aliases", func(t *testing.T) {
		cmd := NewCharacterCmd(llmbackend.MakeDummy)

		if cmd.cmd.Use != "character" {
			t.Errorf("expected Use='character', got %s", cmd.cmd.Use)
		}

		aliases := cmd.cmd.Aliases
		if len(aliases) != 1 || aliases[0] != "char" {
			t.Errorf("expected Aliases=['char'], got %v", aliases)
		}
	})

	t.Run("has three subcommands", func(t *testing.T) {
		cmd := NewCharacterCmd(llmbackend.MakeDummy)

		if !cmd.cmd.HasSubCommands() {
			t.Fatal("character command should have subcommands")
		}

		subcommands := cmd.cmd.Commands()
		if len(subcommands) != 3 {
			t.Fatalf("expected 3 subcommands, got %d", len(subcommands))
		}

		// Check subcommand names
		expectedNames := map[string]bool{
			"gen":     false,
			"regen":   false,
			"gen-img": false,
		}

		for _, sub := range subcommands {
			if _, ok := expectedNames[sub.Use]; ok {
				expectedNames[sub.Use] = true
			}
		}

		for name, found := range expectedNames {
			if !found {
				t.Errorf("expected subcommand %s not found", name)
			}
		}
	})
}

func TestCharacterCmd_SubcommandAliases(t *testing.T) {
	t.Run("character gen subcommand works", func(t *testing.T) {
		cmd := NewCharacterCmd(llmbackend.MakeDummy)

		genCmd, _, err := cmd.cmd.Find([]string{"gen"})
		if err != nil {
			t.Fatalf("failed to find gen subcommand: %v", err)
		}

		if genCmd.Use != "gen" {
			t.Errorf("expected Use='gen', got %s", genCmd.Use)
		}
	})

	t.Run("character regen subcommand works", func(t *testing.T) {
		cmd := NewCharacterCmd(llmbackend.MakeDummy)

		regenCmd, _, err := cmd.cmd.Find([]string{"regen"})
		if err != nil {
			t.Fatalf("failed to find regen subcommand: %v", err)
		}

		if regenCmd.Use != "regen" {
			t.Errorf("expected Use='regen', got %s", regenCmd.Use)
		}
	})

	t.Run("character gen-img subcommand works", func(t *testing.T) {
		cmd := NewCharacterCmd(llmbackend.MakeDummy)

		genImgCmd, _, err := cmd.cmd.Find([]string{"gen-img"})
		if err != nil {
			t.Fatalf("failed to find gen-img subcommand: %v", err)
		}

		if genImgCmd.Use != "gen-img" {
			t.Errorf("expected Use='gen-img', got %s", genImgCmd.Use)
		}
	})
}

func TestCharacterCmd_AliasWorks(t *testing.T) {
	t.Run("char alias works", func(t *testing.T) {
		rootCmd := NewCharacterCmd(llmbackend.MakeDummy)

		// Test that the alias is registered
		aliases := rootCmd.cmd.Aliases
		found := false
		for _, alias := range aliases {
			if alias == "char" {
				found = true
				break
			}
		}

		if !found {
			t.Error("expected 'char' alias to be registered")
		}
	})
}
