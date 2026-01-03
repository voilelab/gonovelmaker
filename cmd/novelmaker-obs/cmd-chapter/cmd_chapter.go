package cmdchapter

import (
	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type ChapterCmd struct {
	genNextCmd  *GenNextCmd
	genEmptyCmd *GenEmptyCmd
	regenCmd    *RegenCmd

	cmd *cobra.Command
}

func NewChapterCmd(llmBackendMaker llmbackend.LLMBackendMaker) *ChapterCmd {
	c := &ChapterCmd{}
	c.cmd = &cobra.Command{
		Use:     "chapter",
		Aliases: []string{"chap"},
		Short:   "Chapter management commands",
		Long:    `Manage chapters in your novel project. Generate new chapters, regenerate existing ones, and create empty chapter placeholders.`,
	}

	// Create subcommands
	c.genNextCmd = NewGenNextCmd(llmBackendMaker)
	c.genEmptyCmd = NewGenEmptyCmd()
	c.regenCmd = NewRegenCmd(llmBackendMaker)

	c.cmd.AddCommand(c.genNextCmd.cmd)
	c.cmd.AddCommand(c.genEmptyCmd.cmd)
	c.cmd.AddCommand(c.regenCmd.cmd)

	return c
}

func (c *ChapterCmd) Command() *cobra.Command {
	return c.cmd
}
