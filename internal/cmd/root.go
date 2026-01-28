package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/s4bb4t/sa-cli/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sac",
	Short: "SA CLI - A command line tool",
	Long:  `SA CLI is a powerful command line tool for managing your workflows.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		verbose, _ := cmd.Flags().GetBool("verbose")
		config.SetDebug(debug)
		config.SetVerbose(verbose)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
}

func Execute() int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		cmd, _, _ := rootCmd.Find(os.Args[1:])
		if cmd == nil || cmd.Name() == rootCmd.Name() {
			rootCmd.PrintErrln("Error:", err)
			rootCmd.PrintErrf("Run '%s --help' for usage.\n", rootCmd.Name())
		} else {
			cmd.PrintErrln("Error:", err)
			cmd.PrintErrf("Run '%s %s --help' for usage.\n", rootCmd.Name(), cmd.Name())
		}
		return 1
	}
	return 0
}
