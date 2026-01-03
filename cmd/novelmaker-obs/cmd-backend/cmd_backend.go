package cmdbackend

import (
	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type BackendCmd struct {
	addCmd                 *addCmd
	removeCmd              *removeCmd
	listCmd                *listCmd
	useCmd                 *useCmd
	checkCmd               *checkCmd
	listAvailableModelsCmd *listAvailableModelsCmd

	cmd *cobra.Command
}

func NewBackendCmd(llmbackendMaker llmbackend.LLMBackendMaker) *BackendCmd {
	backendCmd := &BackendCmd{}
	backendCmd.cmd = &cobra.Command{
		Use:   "backend",
		Short: "Manage LLM backends",
		Long:  `Commands to add, remove, list, and configure LLM backends.`,
	}

	backendCmd.addCmd = NewBackendAddCmd()
	backendCmd.removeCmd = NewBackendRemoveCmd()
	backendCmd.listCmd = NewBackendListCmd()
	backendCmd.useCmd = NewBackendUseCmd()
	backendCmd.checkCmd = NewBackendCheckCmd(llmbackendMaker)
	backendCmd.listAvailableModelsCmd = NewBackendListAvailableModelsCmd(llmbackendMaker)

	backendCmd.cmd.AddCommand(backendCmd.addCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.removeCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.listCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.useCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.checkCmd.cmd)
	backendCmd.cmd.AddCommand(backendCmd.listAvailableModelsCmd.cmd)

	return backendCmd
}

func (b *BackendCmd) Command() *cobra.Command {
	return b.cmd
}
