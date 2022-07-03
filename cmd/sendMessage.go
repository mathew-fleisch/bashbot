/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/mathew-fleisch/bashbot/internal/slack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sendMessageCmd represents the sendMessage command
var sendMessageCmd = &cobra.Command{
	Use:   "sendMessage",
	Short: "Send stand-alone slack message",
	Run: func(cmd *cobra.Command, args []string) {
		channel, _ := cmd.Flags().GetString("channel")
		if channel == "" {
			log.Fatal("--channel flag is required")
		}
		msg, _ := cmd.Flags().GetString("msg")
		user, _ := cmd.Flags().GetString("user")
		configFile, _ := cmd.Flags().GetString("config-file")
		slackToken, _ := cmd.Flags().GetString("slack-token")
		slackAppToken, _ := cmd.Flags().GetString("slack-app-token")
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")
		if logLevel != "" && logFormat != "" {
			slack.ConfigureLogger(logLevel, logFormat)
		}
		slackClient := slack.NewSlackClient(configFile, slackToken, slackAppToken)
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
