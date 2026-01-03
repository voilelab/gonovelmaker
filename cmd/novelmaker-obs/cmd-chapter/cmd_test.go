package cmdchapter

import (
	"testing"

	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

func TestChapterCmd_Structure(t *testing.T) {
	t.Run("has correct aliases", func(t *testing.T) {
		cmd := NewChapterCmd(llmbackend.MakeDummy)

		if cmd.cmd.Use != "chapter" {
			t.Errorf("expected Use='chapter', got %s", cmd.cmd.Use)
		}

		aliases := cmd.cmd.Aliases
		if len(aliases) != 1 || aliases[0] != "chap" {
			t.Errorf("expected Aliases=['chap'], got %v", aliases)
		}
	})

	t.Run("has three subcommands", func(t *testing.T) {
		cmd := NewChapterCmd(llmbackend.MakeDummy)

		if !cmd.cmd.HasSubCommands() {
			t.Fatal("chapter command should have subcommands")
		}

		subcommands := cmd.cmd.Commands()
		if len(subcommands) != 3 {
			t.Fatalf("expected 3 subcommands, got %d", len(subcommands))
		}

		// Check subcommand names
		expectedNames := map[string]bool{
			"gen-next":  false,
			"gen-empty": false,
			"regen":     false,
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

func TestChapterCmd_SubcommandAliases(t *testing.T) {
	t.Run("chapter gen-next subcommand works", func(t *testing.T) {
		cmd := NewChapterCmd(llmbackend.MakeDummy)

		genNextCmd, _, err := cmd.cmd.Find([]string{"gen-next"})
		if err != nil {
			t.Fatalf("failed to find gen-next subcommand: %v", err)
		}

		if genNextCmd.Use != "gen-next" {
			t.Errorf("expected Use='gen-next', got %s", genNextCmd.Use)
		}
	})

	t.Run("chapter gen-empty subcommand works", func(t *testing.T) {
		cmd := NewChapterCmd(llmbackend.MakeDummy)

		genEmptyCmd, _, err := cmd.cmd.Find([]string{"gen-empty"})
		if err != nil {
			t.Fatalf("failed to find gen-empty subcommand: %v", err)
		}

		if genEmptyCmd.Use != "gen-empty" {
			t.Errorf("expected Use='gen-empty', got %s", genEmptyCmd.Use)
		}
	})

	t.Run("chapter regen subcommand works", func(t *testing.T) {
		cmd := NewChapterCmd(llmbackend.MakeDummy)

		regenCmd, _, err := cmd.cmd.Find([]string{"regen"})
		if err != nil {
			t.Fatalf("failed to find regen subcommand: %v", err)
		}

		if regenCmd.Use != "regen" {
			t.Errorf("expected Use='regen', got %s", regenCmd.Use)
		}
	})
}

func TestChapterCmd_AliasWorks(t *testing.T) {
	t.Run("chap alias works", func(t *testing.T) {
		rootCmd := NewChapterCmd(llmbackend.MakeDummy)

		// Test that the alias is registered
		aliases := rootCmd.cmd.Aliases
		found := false
		for _, alias := range aliases {
			if alias == "chap" {
				found = true
				break
			}
		}

		if !found {
			t.Error("expected 'chap' alias to be registered")
		}
	})
}
