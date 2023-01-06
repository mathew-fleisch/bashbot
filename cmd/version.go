package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version = "development"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get current version",
	Run: func(_ *cobra.Command, _ []string) {
		operatingSystem := runtime.GOOS
		systemArchitecture := runtime.GOARCH
		fmt.Printf("bashbot-%s-%s\t %s\n", operatingSystem, systemArchitecture, Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
