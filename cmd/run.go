package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run bashbot",
	Run: func(cmd *cobra.Command, _ []string) {
		slackClient := initSlackClient(cmd)
		slackClient.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
