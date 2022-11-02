package cmd

import (
	"fmt"
	"os"

	"github.com/mathew-fleisch/bashbot/internal/slack"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bashbot",
	Short: "Bashbot Slack bot",
	Long: `
 ____            _     ____        _   
|  _ \          | |   |  _ \      | |  
| |_) | __ _ ___| |__ | |_) | ___ | |_ 
|  _ < / _' / __| '_ \|  _ < / _ \| __|
| |_) | (_| \__ \ | | | |_) | (_) | |_ 
|____/ \__,_|___/_| |_|____/ \___/ \__|
Bashbot is a slack bot, written in golang, that can be configured
to run bash commands or scripts based on a configuration file.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config-file", "", "[REQUIRED] Filepath to config.json file (or environment variable BASHBOT_CONFIG_FILEPATH set)")
	rootCmd.PersistentFlags().String("slack-bot-token", "", "[REQUIRED] Slack bot token used to authenticate with api (or environment variable SLACK_TOKEN set)")
	rootCmd.PersistentFlags().String("slack-app-token", "", "[REQUIRED] Slack app token used to authenticate with api (or environment variable SLACK_APP_TOKEN set)")
	rootCmd.PersistentFlags().String("log-level", "info", "Log elevel to display (info,debug,warn,error)")
	rootCmd.PersistentFlags().String("log-format", "text", "Display logs as json or text")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".bashbot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bashbot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initSlackClient(cmd *cobra.Command) *slack.Client {
	configFile, _ := cmd.Flags().GetString("config-file")
	slackBotToken, _ := cmd.Flags().GetString("slack-bot-token")
	slackAppToken, _ := cmd.Flags().GetString("slack-app-token")
	logLevel, _ := cmd.Flags().GetString("log-level")
	logFormat, _ := cmd.Flags().GetString("log-format")
	if logLevel != "" && logFormat != "" {
		slack.ConfigureLogger(logLevel, logFormat)
	}
	return slack.NewSlackClient(configFile, slackBotToken, slackAppToken)
}
