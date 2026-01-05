package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdbackend "github.com/voilelab/gonovelmaker/cmd/novelmaker-obs/cmd-backend"
	cmdchapter "github.com/voilelab/gonovelmaker/cmd/novelmaker-obs/cmd-chapter"
	cmdcharacter "github.com/voilelab/gonovelmaker/cmd/novelmaker-obs/cmd-character"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "novelmaker-obs",
		Short: "A CLI tool for generating novels in Obsidian vaults",
		Long: `novelmaker-obs is a CLI tool that works with Obsidian vaults to manage 
novel projects, worldbooks, and chapters. It uses OpenAI API to generate new chapters 
based on your existing content.`,
	}

	initCmd := NewInitCmd()
	scanCmd := NewScanCmd()
	exportCmd := NewExportCmd()
	updatePluginCmd := NewUpdatePluginCmd()
	configTestCmd := NewConfigCheckCmd()
	versionCmd := NewVersionCmd()
	rewriteCmd := NewRewriteCmd(llmbackend.MakeOpenAI)

	backendCmd := cmdbackend.NewBackendCmd(llmbackend.MakeOpenAI)
	characterCmd := cmdcharacter.NewCharacterCmd(llmbackend.MakeOpenAI)
	chapterCmd := cmdchapter.NewChapterCmd(llmbackend.MakeOpenAI)

	rootCmd.AddCommand(initCmd.cmd)
	rootCmd.AddCommand(scanCmd.cmd)
	rootCmd.AddCommand(exportCmd.cmd)
	rootCmd.AddCommand(updatePluginCmd.cmd)
	rootCmd.AddCommand(configTestCmd.cmd)
	rootCmd.AddCommand(versionCmd.cmd)
	rootCmd.AddCommand(rewriteCmd.cmd)

	rootCmd.AddCommand(backendCmd.Command())
	rootCmd.AddCommand(characterCmd.Command())
	rootCmd.AddCommand(chapterCmd.Command())

	// Initialize config (creates empty config file if it doesn't exist)
	config.Load()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
