package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sendFileCmd represents the sendFile command
var sendFileCmd = &cobra.Command{
	Use:   "send-file",
	Short: "Send file to slack channel",
	Run: func(cmd *cobra.Command, _ []string) {
		channel, _ := cmd.Flags().GetString("channel")
		if channel == "" {
			log.Fatal("--channel flag is required")
		}
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			log.Fatal("--file flag is required")
		}
		slackClient := initSlackClient(cmd)
		err := slackClient.SendFileToChannel(channel, file)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendFileCmd)

	sendFileCmd.Flags().String("channel", "", "Slack channel to send file to")
	sendFileCmd.Flags().String("file", "", "Filepath/filename to send to slack channel")
}
