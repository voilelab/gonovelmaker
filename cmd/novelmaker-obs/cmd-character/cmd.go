package cmdcharacter

import (
	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type CharacterCmd struct {
	genCmd    *genCmd
	regenCmd  *regenCmd
	genImgCmd *genImgCmd

	cmd *cobra.Command
}

func NewCharacterCmd(llmBackendMaker llmbackend.LLMBackendMaker) *CharacterCmd {
	c := &CharacterCmd{}
	c.cmd = &cobra.Command{
		Use:     "character",
		Aliases: []string{"char"},
		Short:   "Character management commands",
		Long:    `Manage characters in your novel project. Generate new characters, regenerate existing ones, and create character images.`,
	}

	// Create subcommands
	c.genCmd = newGenCmd(llmBackendMaker)
	c.regenCmd = newRegenCmd(llmBackendMaker)
	c.genImgCmd = newGenImgCmd(llmBackendMaker)

	c.cmd.AddCommand(c.genCmd.cmd)
	c.cmd.AddCommand(c.regenCmd.cmd)
	c.cmd.AddCommand(c.genImgCmd.cmd)

	return c
}

func (c *CharacterCmd) Command() *cobra.Command {
	return c.cmd
}
