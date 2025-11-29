package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

type ExportCmd struct {
	outFile  string
	fileType string

	cmd *cobra.Command
}

func NewExportCmd() *ExportCmd {
	exportCmd := &ExportCmd{}
	exportCmd.cmd = &cobra.Command{
		Use:   "export",
		Short: "Export the novel project to the desired format",
		Long:  `Exports the novel project to formats such as PDF, ePub, or HTML.`,
		RunE:  exportCmd.run,
	}

	exportCmd.cmd.Flags().StringVarP(&exportCmd.outFile, "output", "o", "", "Output file for the exported novel")
	exportCmd.cmd.Flags().StringVarP(&exportCmd.fileType, "type", "t", "txt", "Export file type (txt)")
	exportCmd.cmd.MarkFlagRequired("output")

	return exportCmd
}

func (e *ExportCmd) run(cmd *cobra.Command, args []string) error {
	if e.fileType != "txt" {
		return fmt.Errorf("unsupported export file type: %s", e.fileType)
	}

	fmt.Printf("Exporting novel to %s as %s format...\n", e.outFile, e.fileType)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create vault
	vault, err := obsidian.NewVault(cwd)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}
	defer vault.Close()

	// Load project
	project, err := vault.LoadProject()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Load chapters
	chapters, err := vault.LoadChapters()
	if err != nil {
		return fmt.Errorf("failed to load chapters: %w", err)
	}

	// Export to text file
	outFile, err := os.Create(e.outFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Write project title and author
	_, err = fmt.Fprintf(outFile, "%s\n\n", project.Name)
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	// Write chapters
	for _, ch := range chapters {
		_, err = fmt.Fprintf(outFile, "Chapter %d: %s\n\n%s\n\n", ch.Index, ch.Title, ch.Content)
		if err != nil {
			return fmt.Errorf("failed to write chapter to output file: %w", err)
		}
	}

	fmt.Println("✓ Export completed successfully!")

	return nil
}
