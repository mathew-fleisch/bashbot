package cmd

import (
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
		slackClient := initSlackClient(cmd)
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
