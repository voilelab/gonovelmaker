package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v0.0.6"

type VersionCmd struct {
	cmd *cobra.Command
}

func NewVersionCmd() *VersionCmd {
	v := &VersionCmd{}
	v.cmd = &cobra.Command{
		Use:   "version",
		Short: "Display the version of Novel Maker CLI",
		Long:  `Displays the current version of the Novel Maker CLI tool.`,
		RunE:  v.run,
	}
	return v
}

func (v *VersionCmd) run(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s\n", Version)
	return nil
}
