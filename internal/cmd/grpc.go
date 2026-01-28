package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/s4bb4t/sa-cli/internal/scaffold"
	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc [service-name]",
	Short: "Generate gRPC service scaffold",
	Long: `Generate a gRPC service with proto contract and server implementation.

Example:
  sac project grpc user
  sac project grpc user --module github.com/myorg/user-service
  sac project grpc user --output ./services/user`,
	Args: cobra.ExactArgs(1),
	RunE: runGrpc,
}

var (
	grpcModule string
	grpcOutput string
)

func init() {
	grpcCmd.Flags().StringVarP(&grpcModule, "module", "m", "", "Go module path (default: service name)")
	grpcCmd.Flags().StringVarP(&grpcOutput, "output", "o", ".", "Output directory (default: current dir)")

	projectCmd.AddCommand(grpcCmd)
}

func runGrpc(_ *cobra.Command, args []string) error {
	serviceName := args[0]

	output := grpcOutput
	if output == "" {
		output = "."
	}

	absPath, err := filepath.Abs(output)
	if err != nil {
		return err
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		if err := os.MkdirAll(output, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	fmt.Printf("Generating gRPC service %q in %s\n\n", serviceName, absPath)

	grpcScaffold := scaffold.NewGRPC(serviceName, grpcModule, output)
	if err := grpcScaffold.Generate(); err != nil {
		return fmt.Errorf("failed to generate gRPC service: %w", err)
	}

	return nil
}
