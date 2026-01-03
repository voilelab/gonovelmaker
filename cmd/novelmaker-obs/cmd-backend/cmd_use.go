package cmdbackend

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
)

type useCmd struct {
	cmd *cobra.Command
}

func newUseCmd() *useCmd {
	useCmd := &useCmd{}
	useCmd.cmd = &cobra.Command{
		Use:   "use <name>",
		Short: "Set the default LLM backend to use",
		Long:  `Set the default LLM backend that will be used by generation commands.`,
		Args:  cobra.ExactArgs(1),
		RunE:  useCmd.run,
	}

	return useCmd
}

func (b *useCmd) run(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Set the default backend
	if err := cfg.SetDefaultBackend(name); err != nil {
		return err
	}

	// Save the config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Default backend set to '%s'\n", name)

	return nil
}
