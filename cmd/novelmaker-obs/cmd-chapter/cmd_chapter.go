package cmdchapter

import (
	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type ChapterCmd struct {
	genNextCmd      *GenNextCmd
	genNextEmptyCmd *GenNextEmptyCmd
	genCurrCmd      *GenCurrCmd

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
	c.genNextCmd.cmd.Use = "gen-next"
	c.genNextCmd.cmd.Short = "Generate the next chapter using AI"
	c.genNextCmd.cmd.Long = `Generates a new chapter based on existing project configuration, worldbook, 
and previous chapters using the configured LLM backend.`

	c.genNextEmptyCmd = NewGenNextEmptyCmd()
	c.genNextEmptyCmd.cmd.Use = "gen-empty"
	c.genNextEmptyCmd.cmd.Short = "Generate an empty next chapter"
	c.genNextEmptyCmd.cmd.Long = `Creates a new empty chapter file with frontmatter but no content. 
This is useful for manually writing chapters or creating placeholders.`

	c.genCurrCmd = NewGenCurrCmd(llmBackendMaker)
	c.genCurrCmd.cmd.Use = "regen"
	c.genCurrCmd.cmd.Short = "Regenerate an existing chapter"
	c.genCurrCmd.cmd.Long = `Regenerates an existing chapter based on the prompt stored in its frontmatter.
The filepath should be relative to the vault root (e.g., "Story/001_ch1.md").`

	c.cmd.AddCommand(c.genNextCmd.cmd)
	c.cmd.AddCommand(c.genNextEmptyCmd.cmd)
	c.cmd.AddCommand(c.genCurrCmd.cmd)

	return c
}

func (c *ChapterCmd) Command() *cobra.Command {
	return c.cmd
}
