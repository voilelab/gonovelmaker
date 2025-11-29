package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "novelmaker-obs",
		Short: "A CLI tool for generating novels in Obsidian vaults",
		Long: `novelmaker-obs is a CLI tool that works with Obsidian vaults to manage 
novel projects, worldbooks, and chapters. It uses OpenAI API to generate new chapters 
based on your existing content.

Configuration:
  Config file: ~/.novelmaker/config.toml
  Example:
    openai_api_key = "sk-xxx"
    model = "gpt-4o"

  You can also set OPENAI_API_KEY as an environment variable.`,
	}

	initCmd := NewInitCmd()
	scanCmd := NewScanCmd()
	genNextCmd := NewGenNextCmd()
	genCharCmd := NewGenCharCmd()
	exportCmd := NewExportCmd()
	updatePluginCmd := NewUpdatePluginCmd()

	rootCmd.AddCommand(initCmd.cmd)
	rootCmd.AddCommand(scanCmd.cmd)
	rootCmd.AddCommand(genNextCmd.cmd)
	rootCmd.AddCommand(genCharCmd.cmd)
	rootCmd.AddCommand(exportCmd.cmd)
	rootCmd.AddCommand(updatePluginCmd.cmd)

	// Initialize config (creates empty config file if it doesn't exist)
	config.Load()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
