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
	genNextCmd := NewGenNextCmd(openAIBackendMaker)
	genNextEmptyCmd := NewGenNextEmptyCmd()
	genCurrCmd := NewGenCurrCmd(openAIBackendMaker)
	genCharCmd := NewGenCharCmd(openAIBackendMaker)
	genCharCurrCmd := NewGenCharCurrCmd(openAIBackendMaker)
	genCharImgCmd := NewGenCharImgCmd(openAIBackendMaker)
	exportCmd := NewExportCmd()
	updatePluginCmd := NewUpdatePluginCmd()
	configTestCmd := NewConfigCheckCmd()

	// Backend management commands
	backendCmd := &cobra.Command{
		Use:   "backend",
		Short: "Manage LLM backend configurations",
		Long:  `Add, remove, list, and configure LLM backends for novel generation.`,
	}
	backendAddCmd := NewBackendAddCmd()
	backendRemoveCmd := NewBackendRemoveCmd()
	backendUseCmd := NewBackendUseCmd()
	backendListCmd := NewBackendListCmd()
	backendCheckCmd := NewBackendCheckCmd()

	backendCmd.AddCommand(backendAddCmd.cmd)
	backendCmd.AddCommand(backendRemoveCmd.cmd)
	backendCmd.AddCommand(backendUseCmd.cmd)
	backendCmd.AddCommand(backendListCmd.cmd)
	backendCmd.AddCommand(backendCheckCmd.cmd)

	rootCmd.AddCommand(initCmd.cmd)
	rootCmd.AddCommand(scanCmd.cmd)
	rootCmd.AddCommand(genNextCmd.cmd)
	rootCmd.AddCommand(genNextEmptyCmd.cmd)
	rootCmd.AddCommand(genCurrCmd.cmd)
	rootCmd.AddCommand(genCharCmd.cmd)
	rootCmd.AddCommand(genCharCurrCmd.cmd)
	rootCmd.AddCommand(genCharImgCmd.cmd)
	rootCmd.AddCommand(exportCmd.cmd)
	rootCmd.AddCommand(updatePluginCmd.cmd)
	rootCmd.AddCommand(configTestCmd.cmd)
	rootCmd.AddCommand(backendCmd)

	// Initialize config (creates empty config file if it doesn't exist)
	config.Load()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
