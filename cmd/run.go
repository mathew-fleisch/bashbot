package cmd

import (
	"github.com/mathew-fleisch/bashbot/internal/slack"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run bashbot",
	Run: func(cmd *cobra.Command, _ []string) {
		configFile, _ := cmd.Flags().GetString("config-file")
		slackBotToken, _ := cmd.Flags().GetString("slack-bot-token")
		slackAppToken, _ := cmd.Flags().GetString("slack-app-token")
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")
		slack.ConfigureLogger(logLevel, logFormat)
		slackClient := slack.NewSlackClient(configFile, slackBotToken, slackAppToken)
		slackClient.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
