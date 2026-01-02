package cmdbackend

import (
	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type BackendCmd struct {
	backendAddCmd    *BackendAddCmd
	backendRemoveCmd *BackendRemoveCmd
	backendListCmd   *BackendListCmd
	backendUseCmd    *BackendUseCmd
	backendCheckCmd  *BackendCheckCmd

	cmd *cobra.Command
}

func NewBackendCmd(llmbackendMaker llmbackend.LLMBackendMaker) *BackendCmd {
	backendCmd := &BackendCmd{}
	backendCmd.cmd = &cobra.Command{
		Use:   "backend",
		Short: "Manage LLM backends",
		Long:  `Commands to add, remove, list, and configure LLM backends.`,
	}

	backendCmd.backendAddCmd = NewBackendAddCmd()
	backendCmd.backendRemoveCmd = NewBackendRemoveCmd()
	backendCmd.backendListCmd = NewBackendListCmd()
	backendCmd.backendUseCmd = NewBackendUseCmd()
	backendCmd.backendCheckCmd = NewBackendCheckCmd(llmbackendMaker)

	backendCmd.cmd.AddCommand(backendCmd.backendAddCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.backendRemoveCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.backendListCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.backendUseCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.backendCheckCmd.cmd)

	return backendCmd
}

func (b *BackendCmd) Command() *cobra.Command {
	return b.cmd
}
