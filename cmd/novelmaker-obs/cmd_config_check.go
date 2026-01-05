package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

type ConfigCheckCmd struct {
	cmd *cobra.Command
}

func NewConfigCheckCmd() *ConfigCheckCmd {
	c := &ConfigCheckCmd{}

	c.cmd = &cobra.Command{
		Use:   "config-check",
		Short: "Check if chapter, character, and rewrite prompt templates can be parsed successfully",
		Long: `Check if chapter, character, and rewrite prompt templates can be parsed successfully.

This command reads Config/chapter_prompt.md, Config/character_prompt.md, and 
Config/rewrite_prompt.md from the vault and attempts to parse them. It reports 
any errors in the frontmatter or template syntax.

Example:
  novelmaker-obs config-check
  novelmaker-obs config-check --vault /path/to/vault`,
		RunE: c.run,
	}

	c.cmd.Flags().StringP("vault", "v", ".", "Path to the Obsidian vault")

	return c
}

func (c *ConfigCheckCmd) run(cmd *cobra.Command, args []string) error {
	vaultPath, _ := cmd.Flags().GetString("vault")

	// Open the vault
	vault, err := obsidian.NewVault(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	fmt.Println("Checking prompt templates...")
	fmt.Println()

	// Test chapter prompt
	fmt.Println("📄 Checking Config/chapter_prompt.md...")
	chapterPrompt, err := vault.LoadChapterPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to parse chapter prompt: %v\n", err)
		return err
	}
	fmt.Println("✅ Chapter prompt parsed successfully")
	if chapterPrompt.System != "" {
		fmt.Printf("   System prompt: %.80s...\n", chapterPrompt.System)
	}
	fmt.Println()

	// Test character prompt
	fmt.Println("📄 Checking Config/character_prompt.md...")
	characterPrompt, err := vault.LoadCharacterPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to parse character prompt: %v\n", err)
		return err
	}
	fmt.Println("✅ Character prompt parsed successfully")
	if characterPrompt.System != "" {
		fmt.Printf("   System prompt: %.80s...\n", characterPrompt.System)
	}
	fmt.Println()

	// Test rewrite prompt
	fmt.Println("📄 Checking Config/rewrite_prompt.md...")
	rewritePrompt, err := vault.LoadRewritePrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to parse rewrite prompt: %v\n", err)
		return err
	}
	fmt.Println("✅ Rewrite prompt parsed successfully")
	if rewritePrompt.System != "" {
		fmt.Printf("   System prompt: %.80s...\n", rewritePrompt.System)
	}
	fmt.Println()

	fmt.Println("🎉 All prompt templates are valid!")
	return nil
}
