package cmd

import (
	"github.com/mathew-fleisch/bashbot/internal/slack"
	"github.com/spf13/cobra"
)

// installDependenciesCmd represents the installDependencies command
var installDependenciesCmd = &cobra.Command{
	Use:   "install-dependencies",
	Short: "Cycle through dependencies array in config file to install extra dependencies",
	Run: func(cmd *cobra.Command, _ []string) {
		configFile, _ := cmd.Flags().GetString("config-file")
		slackBotToken, _ := cmd.Flags().GetString("slack-bot-token")
		slackAppToken, _ := cmd.Flags().GetString("slack-app-token")
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")
		if logLevel != "" && logFormat != "" {
			slack.ConfigureLogger(logLevel, logFormat)
		}
		slackClient := slack.NewSlackClient(configFile, slackBotToken, slackAppToken)
		slackClient.InstallVendorDependencies()
	},
}

func init() {
	rootCmd.AddCommand(installDependenciesCmd)
}
