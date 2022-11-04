package cmd

import (
	"github.com/spf13/cobra"
)

// installDependenciesCmd represents the installDependencies command
var installDependenciesCmd = &cobra.Command{
	Use:   "install-dependencies",
	Short: "Cycle through dependencies array in config file to install extra dependencies",
	Run: func(cmd *cobra.Command, _ []string) {
		slackClient := initSlackClient(cmd)
		slackClient.InstallVendorDependencies()
	},
}

func init() {
	rootCmd.AddCommand(installDependenciesCmd)
}
