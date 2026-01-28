package cmd

import (
	"fmt"
	"runtime"

	"github.com/s4bb4t/sa-cli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print detailed version information including build date and git commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		short, _ := cmd.Flags().GetBool("short")
		if short {
			fmt.Println(version.Version)
			return
		}
		fmt.Printf("sa-cli %s\n", version.Version)
		fmt.Printf("  Commit:     %s\n", version.GitCommit)
		fmt.Printf("  Built:      %s\n", version.BuildDate)
		fmt.Printf("  Go version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	versionCmd.Flags().BoolP("short", "s", false, "print only the version number")
	rootCmd.AddCommand(versionCmd)
}
