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
	"fmt"

	"github.com/mathew-fleisch/bashbot/internal/slack"
	"github.com/spf13/cobra"
)

// installDependenciesCmd represents the installDependencies command
var installDependenciesCmd = &cobra.Command{
	Use:   "installDependencies",
	Short: "Cycle through dependencies array in config file to install extra dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("installDependencies called")
		configFile, _ := cmd.Flags().GetString("config-file")
		slackToken, _ := cmd.Flags().GetString("slack-token")
		slackAppToken, _ := cmd.Flags().GetString("slack-app-token")
		logLevel, _ := cmd.Flags().GetString("log-level")
		logFormat, _ := cmd.Flags().GetString("log-format")
		if logLevel != "" && logFormat != "" {
			slack.ConfigureLogger(logLevel, logFormat)
		}
		slackClient := slack.NewSlackClient(configFile, slackToken, slackAppToken)
		slackClient.InstallVendorDependencies()
	},
}

func init() {
	rootCmd.AddCommand(installDependenciesCmd)
}
