package cmd

import (
	"fmt"

	"github.com/s4bb4t/sa-cli/internal/config"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "Run a task",
	Long:  `Run a specified task with optional configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		cfg := config.Get()

		if cfg.Verbose {
			fmt.Printf("Running task: %s\n", name)
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		if dryRun {
			fmt.Printf("Dry run: would execute task '%s'\n", name)
			return nil
		}

		fmt.Printf("Executing task: %s\n", name)
		return nil
	},
}

func init() {
	runCmd.Flags().Bool("dry-run", false, "show what would be executed without running")
	rootCmd.AddCommand(runCmd)
}
