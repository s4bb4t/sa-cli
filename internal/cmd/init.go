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

The command scans the input directory for api/proto/ and api/openapi/
subdirectories to determine the project type (gRPC, REST, or hybrid).

Example:
  sac project init myapp
  sac project init myapp --module github.com/myorg/myapp
  sac project init myapp --output ./projects/myapp --input ./api`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

var (
	initModule string
	initOutput string
	initInput  string
)

func init() {
	initCmd.Flags().StringVarP(&initModule, "module", "m", "app", "Go module path (default: project name)")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "./", "Output directory (default: ./)")
	initCmd.Flags().StringVarP(&initInput, "input", "i", "api/", "Input directory (default: api/)")

	projectCmd.AddCommand(initCmd)
	rootCmd.AddCommand(projectCmd)
}

func runInit(_ *cobra.Command, args []string) error {
	name := "app"
	if len(args) > 0 {
		name = args[0]
	}

	//if _, err := os.Stat(initInput); err != nil {
	//	return fmt.Errorf("directory %q does not exist", initInput)
	//}
	//
	//if _, err := os.Stat(initOutput); err == nil {
	//	return fmt.Errorf("directory %q already exists", initOutput)
	//}

	hasGRPC := dirExists(filepath.Join(initInput, "proto"))
	hasOpenAPI := dirExists(filepath.Join(initInput, "openapi"))

	if !hasGRPC && !hasOpenAPI {
		return fmt.Errorf("no api specs found in input directory %q (expected proto/ and/or openapi/ subdirectories)", initInput)
	}

	var mode scaffold.ProjectMode
	if hasGRPC {
		mode |= scaffold.ModeGRPC
	}
	if hasOpenAPI {
		mode |= scaffold.ModeOpenAPI
	}

	absPath, err := filepath.Abs(initOutput)
	if err != nil {
		return err
	}

	modeDesc := "gRPC"
	if hasGRPC && hasOpenAPI {
		modeDesc = "gRPC + OpenAPI (hybrid)"
	} else if hasOpenAPI {
		modeDesc = "OpenAPI (REST)"
	}

	fmt.Printf("Creating %s project %q in %s\n\n", modeDesc, name, absPath)

	project := scaffold.New(name, initModule, initOutput, mode)
	if err := project.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	if hasGRPC {
		fmt.Println("\n  Scaffolding gRPC service...")
		grpcScaffold := scaffold.NewGRPC(name, initOutput)
		if err := grpcScaffold.Generate(); err != nil {
			return fmt.Errorf("failed to generate gRPC service: %w", err)
		}
	}

	if hasOpenAPI {
		fmt.Println("\n  Scaffolding OpenAPI service...")
		openapiScaffold := scaffold.NewOpenAPI(name, initModule, initOutput, initInput)
		if err := openapiScaffold.Generate(); err != nil {
			return fmt.Errorf("failed to generate OpenAPI service: %w", err)
		}
	}

	fmt.Printf("\nProject %q created successfully!\n", name)
	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
