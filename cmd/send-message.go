package cmd

import (
	"github.com/mathew-fleisch/bashbot/internal/slack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sendMessageCmd represents the sendMessage command
var sendMessageCmd = &cobra.Command{
	Use:   "send-message",
	Short: "Send stand-alone slack message",
	Run: func(cmd *cobra.Command, _ []string) {
		channel, _ := cmd.Flags().GetString("channel")
		if channel == "" {
			log.Fatal("--channel flag is required")
		}
		msg, _ := cmd.Flags().GetString("msg")
		user, _ := cmd.Flags().GetString("user")
		configFile, _ := cmd.Flags().GetString("config-file")
		slackBotToken, _ := cmd.Flags().GetString("slack-bot-token")
		slackAppToken, _ := cmd.Flags().GetString("slack-app-token")
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")
		if logLevel != "" && logFormat != "" {
			slack.ConfigureLogger(logLevel, logFormat)
		}
		slackClient := slack.NewSlackClient(configFile, slackBotToken, slackAppToken)
		if user != "" {
			slackClient.SendMessageToUser(channel, user, msg)
			return
		}
		slackClient.SendMessageToChannel(channel, msg)
	},
}

func init() {
	rootCmd.AddCommand(sendMessageCmd)

	sendMessageCmd.Flags().String("channel", "", "Slack channel to send stand-alone message to")
	sendMessageCmd.Flags().String("user", "", "Slack user to send stand-alone message to")
	sendMessageCmd.Flags().String("msg", "", "Message to send to slack channel/user")
}
