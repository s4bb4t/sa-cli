package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/s4bb4t/sa-cli/internal/scaffold"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project management commands",
	Long:  `Commands for creating and managing Go projects.`,
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new project",
	Long: `Initialize a new Go project with a production-ready structure.

Example:
  sac project init myapp
  sac project init myapp --module github.com/myorg/myapp
  sac project init myapp --output ./projects/myapp`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

var (
	initModule string
	initOutput string
)

func init() {
	initCmd.Flags().StringVarP(&initModule, "module", "m", "", "Go module path (default: project name)")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "", "Output directory (default: project name)")

	projectCmd.AddCommand(initCmd)
	rootCmd.AddCommand(projectCmd)
}

func runInit(_ *cobra.Command, args []string) error {
	name := "app"
	if len(args) > 0 {
		name = args[0]
	}

	output := initOutput

	if _, err := os.Stat(output); err == nil {
		return fmt.Errorf("directory %q already exists", output)
	}

	absPath, err := filepath.Abs(output)
	if err != nil {
		return err
	}

	fmt.Printf("Creating project %q in %s\n\n", name, absPath)

	project := scaffold.New(name, initModule, output)
	if err := project.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	return nil
}
