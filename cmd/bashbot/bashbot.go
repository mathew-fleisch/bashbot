package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/slack-go/slack"
)

var specials []func(event *slack.MessageEvent) bool

// Slacking off with global vars
var Version = "development"
var help bool
var getVersion bool
var configFile string
var slackToken string
var installVendorDependenciesFlag bool
var sendMessageChannel string
var sendMessageText string
var sendMessageEphemeral bool
var sendMessageUser string
var logLevel string
var logFormat string
var api *slack.Client
var rtm *slack.RTM
var channelsByName map[string]string
var countkey string
var emojiPattern *regexp.Regexp
var slackUserPattern *regexp.Regexp
var puncPattern *regexp.Regexp
var c *regexp.Regexp
var cmdPattern *regexp.Regexp

type Admins struct {
	Admins []Admin `json:"admins"`
}

type Admin struct {
	Trigger          string   `json:"trigger"`
	AppName          string   `json:"appName"`
	UserIds          []string `json:"userIds"`
	PrivateChannelId string   `json:"privateChannelId"`
	LogChannelId     string   `json:"logChannelId"`
}

var admin Admin

type Messages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Active bool   `json:"active"`
	Name   string `json:"name"`
	Text   string `json:"text"`
}

type Tools struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Help        string      `json:"help"`
	Trigger     string      `json:"trigger"`
	Location    string      `json:"location"`
	Setup       string      `json:"setup"`
	Command     string      `json:"command"`
	Permissions []string    `json:"permissions"`
	Log         bool        `json:"log"`
	Ephemeral   bool        `json:"ephemeral"`
	Response    string      `json:"response"`
	Parameters  []Parameter `json:"parameters"`
}

type Parameter struct {
	Name        string   `json:"name"`
	Allowed     []string `json:"allowed"`
	Description string   `json:"description,omitempty"`
	Source      string   `json:"source,omitempty"`
}

type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency struct {
	Name    string `json:"name"`
	Install string `json:"install"`
}

type Channel struct {
	Id                 string `json:"id"`
	Created            int    `json:"created"`
	IsOpen             bool   `json:"is_open"`
	IsGroup            bool   `json:"is_group"`
	IsShared           bool   `json:"is_shared"`
	IsIm               bool   `json:"is_im"`
	IsExtShared        bool   `json:"is_ext_shared"`
	IsOrgShared        bool   `json:"is_org_shared"`
	IsPendingExtShared bool   `json:"is_pending_ext_shared"`
	IsPrivate          bool   `json:"is_private"`
	IsMpim             bool   `json:"is_mpim"`
	Unlinked           int    `json:"unlinked"`
	NameNormalized     string `json:"name_normalized"`
	NumMembers         int    `json:"num_members"`
	Priority           int    `json:"priority"`
	User               string `json:"user"`
	Name               string `json:"name"`
	Creator            string `json:"creator"`
	IsArchived         bool   `json:"is_archived"`
	Members            string `json:"members"`
	Topic              Topic  `json:"topic"`
	Purpose            Topic  `json:"purpose"`
	IsChannel          bool   `json:"is_channel"`
	IsGeneral          bool   `json:"is_general"`
	IsMember           bool   `json:"is_member"`
	Local              string `json:"locale"`
}

type Topic struct {
	Value   string `json:"value"`
	Creator string `json:"creator"`
	LastSet int    `json:"last_set"`
}

func getChannelNames(channelIds []string) []string {
	log.Debug("getChannelNames()")
	var names []string
	pparams := &slack.GetConversationsParameters{Limit: 1000, Types: []string{"private_channel"}}
	pchannels, _, _ := api.GetConversations(pparams)
	pnumChannels := len(pchannels)
	log.Debug("Number of private channels this bot is monitoring: " + strconv.Itoa(pnumChannels))
	for j := 0; j < pnumChannels; j++ {
		pthisChannel, _ := json.Marshal(pchannels[j])
		var pthatChannel Channel
		json.Unmarshal([]byte(pthisChannel), &pthatChannel)
		for i := 0; i < len(channelIds); i++ {
			if channelIds[i] == pthatChannel.Id {
				names = append(names, pthatChannel.Name)
			}
		}
	}
	params := &slack.GetConversationsParameters{Limit: 1000, Types: []string{"public_channel"}}
	channels, _, _ := api.GetConversations(params)
	numChannels := len(channels)
	log.Debug("Number of public channels this bot is monitoring: " + strconv.Itoa(numChannels))
	for j := 0; j < numChannels; j++ {
		thisChannel, _ := json.Marshal(channels[j])
		var thatChannel Channel
		json.Unmarshal([]byte(thisChannel), &thatChannel)
		for i := 0; i < len(channelIds); i++ {
			if channelIds[i] == thatChannel.Id {
				names = append(names, thatChannel.Name)
			}
		}
	}
	if len(names) == 0 {
		names = append(names, "all")
	}
	return names
}

func getAdmin() Admin {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Error("Could not open config file")
		log.Error(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var Admins Admins
	err = json.Unmarshal(byteValue, &Admins)
	if err != nil {
		log.Error(fmt.Printf("this error: %s", err.Error()))
		log.Error(fmt.Println(err))
	}
	switch err := err.(type) {
	case *json.SyntaxError:
		log.Error(fmt.Printf("Syntax error (at byte: %d) in config.json file:\n\n %s\n", err.Offset, err.Error()))
		os.Exit(1)
	default:
		log.Debug("config.json parsed successfully")
	}
	return Admins.Admins[0]
}

func installVendorDependencies() bool {
	log.Debug("installVendorDependencies()")
	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Error(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var Dependencies Dependencies
	json.Unmarshal(byteValue, &Dependencies)

	for i := 0; i < len(Dependencies.Dependencies); i++ {
		log.Info(Dependencies.Dependencies[i].Name)

		words := strings.Fields(Dependencies.Dependencies[i].Install)
		var tcmd []string

		for index, element := range words {
			log.Debug(strconv.Itoa(index) + ": " + element)
			tcmd = append(tcmd, element)
		}
		cmd := []string{"bash", "-c", "pushd vendor && " + strings.Join(tcmd, " ") + " && popd"}
		log.Debug(strings.Join(cmd, " "))
		ret := shellOut(cmd)
		log.Info(ret)
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

// Slack Command Processing
func processCommand(event *slack.MessageEvent) bool {
	if !cmdPattern.MatchString(event.Text) {
		return false
	}

	log.Info("command detected: `" + event.Text + "`")
	log.Debug(event)
	log.Info("Channel: " + event.Channel)
	log.Info("User: " + event.User)
	log.Info("Timestamp: " + event.Timestamp)

	words := strings.Fields(event.Text)
	var triggered string
	var thisTool Tool
	var cmd []string

	for index, element := range words {
		log.Info(strconv.Itoa(index) + ": " + element)
		if index > 1 {
			cmd = append(cmd, element)
		}
	}

	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Error(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var Tools Tools
	json.Unmarshal(byteValue, &Tools)

	for i := 0; i < len(Tools.Tools); i++ {
		if Tools.Tools[i].Trigger == words[1] {
			triggered = Tools.Tools[i].Trigger
			thisTool = Tools.Tools[i]
		}
	}

	switch words[1] {
	case triggered:
		reportToChannel(event.Channel, "processing_command", "")
		return processWhitelistedCommand(cmd, thisTool, event.Channel, event.User, event.Timestamp)
	case "cmd":
		reportToChannel(event.Channel, "processing_raw_command", "")
		return processRawCommand(cmd, event.Channel, event.User)
	default:
		reportToChannel(event.Channel, "command_not_found", "")
		return false

	}
}

// Whitelisted commands
func processWhitelistedCommand(cmds []string, thisTool Tool, channel string, user string, timestamp string) bool {
	validParams := make([]bool, len(thisTool.Parameters))
	var tmpHelp string
	authorized := false

	// inject email if exists in command
	reEmail := regexp.MustCompile(`\${email}`)
	thisUser, err := api.GetUserInfo(user)
	if err != nil {
		log.Info(fmt.Printf("%s\n", err))
		return true
	}
	thisTool.Command = reEmail.ReplaceAllLiteralString(thisTool.Command, thisUser.Profile.Email)

	log.Debug("Tool Name:        " + thisTool.Name)
	log.Debug("Tool Description: " + thisTool.Description)
	log.Debug("Tool Log:         " + strconv.FormatBool(thisTool.Log))
	log.Debug("Tool Help:        " + thisTool.Help)
	log.Debug("Tool Trigger:     " + thisTool.Trigger)
	log.Debug("Tool Location:    " + thisTool.Location)
	log.Debug("Tool Setup:       " + thisTool.Setup)
	log.Debug("Tool Command:     " + thisTool.Command)
	log.Debug("Tool Ephemeral:   " + strconv.FormatBool(thisTool.Ephemeral))
	log.Debug("Tool Response:    " + thisTool.Response)
	var allowedChannels []string = getChannelNames(thisTool.Permissions)
	if admin.PrivateChannelId == channel {
		authorized = true
	} else {
		for j := 0; j < len(thisTool.Permissions); j++ {
			log.Debug("Tool Permissions[" + strconv.Itoa(j) + "]: " + thisTool.Permissions[j])
			if thisTool.Permissions[j] == channel || thisTool.Permissions[j] == "all" {
				authorized = true
			}
		}
	}

	// Show help if the first parameter is "help"
	cmdHelp := fmt.Sprintf("``` ====> %s [Allowed In: %s] <====\n%s\n%s%s```", thisTool.Name, strings.Join(allowedChannels, ", "), thisTool.Description, thisTool.Help, tmpHelp)
	if len(cmds) > 0 {
		for j := 0; j < len(cmds); j++ {
			if cmds[j] == "help" {
				yell(channel, cmdHelp)
				return true
			}
		}
	}

	if !authorized {
		reportToChannel(channel, "unauthorized", strings.Join(allowedChannels, ", "))
		yell(channel, cmdHelp)
		chatOpsLog(channel, user, thisTool.Trigger+" "+strings.Join(cmds, " "))
		return true
	}

	if len(thisTool.Parameters) > 0 {
		log.Debug("Tool Parameters Count: " + strconv.Itoa(len(thisTool.Parameters)))
		for j := 0; j < len(thisTool.Parameters); j++ {
			log.Debug("Tool Parameters[" + strconv.Itoa(j) + "]: " + thisTool.Parameters[j].Name)
			derivedSource := thisTool.Parameters[j].Source
			tmpHelp = fmt.Sprintf("%s\n%s: [%s%s]", tmpHelp, thisTool.Parameters[j].Name, strings.Join(thisTool.Parameters[j].Allowed, "|"), thisTool.Parameters[j].Description)
			if len(derivedSource) > 0 {
				log.Debug("No hard-coded allowed values. Deriving source: " + derivedSource)
				allowedOut := shellOut([]string{"bash", "-c", "cd " + thisTool.Location + " && " + derivedSource})
				log.Debug("Derived: " + allowedOut)
				thisTool.Parameters[j].Allowed = strings.Split(allowedOut, "\n")
			}
			// tmpHelp = fmt.Sprintf("%s\n%s: [%s]", tmpHelp, thisTool.Parameters[j].Name, strings.Join(thisTool.Parameters[j].Allowed, "|"))
			// for h := 0; h < len(thisTool.Parameters[j].Allowed); h++ {
			//   fmt.Println("Tool Parameters[" + strconv.Itoa(j) + "].Allowed[" + strconv.Itoa(h) + "]: " + thisTool.Parameters[j].Allowed[h])
			// }
		}
	}

	if thisTool.Log {
		chatOpsLog(channel, user, thisTool.Trigger+" "+strings.Join(cmds, " "))
	}

	// Verify all required parameters are passed
	if len(cmds) != len(thisTool.Parameters) {
		reportToChannel(channel, "incorrect_parameters", "")
		yell(channel, cmdHelp)
		return true
	}

	// Validate parameters against whitelist
	if len(thisTool.Parameters) > 0 {
		for j := 0; j < len(thisTool.Parameters); j++ {
			validParams[j] = false
			for h := 0; h < len(thisTool.Parameters[j].Allowed); h++ {
				if thisTool.Parameters[j].Allowed[h] == cmds[j] {
					validParams[j] = true
				}
			}
		}
	}

	buildCmd := thisTool.Command
	for x := 0; x < len(cmds); x++ {
		if !validParams[x] {
			reportToChannel(channel, "invalid_parameter", thisTool.Parameters[x].Name)
			return false
		}
		re := regexp.MustCompile(`\${` + thisTool.Parameters[x].Name + `}`)
		buildCmd = re.ReplaceAllString(buildCmd, cmds[x])
	}
	buildCmd = getUserChannelInfo(user, thisUser.Name, channel, timestamp) + " && cd " + thisTool.Location + " && " + thisTool.Setup + " && " + buildCmd
	splitOn := regexp.MustCompile(`\s\&\&`)
	displayCmd := splitOn.ReplaceAllString(buildCmd, " \\\n        &&")
	log.Info("Triggered Command:")
	log.Info(displayCmd)

	tmpCmd := []string{"bash", "-c", buildCmd}

	ret := splitOut(shellOut(tmpCmd), thisTool.Response)

	if thisTool.Ephemeral {
		reportToChannel(channel, "ephemeral", "")
		whisper(channel, user, ret)
	} else {
		yell(channel, ret)
	}
	if thisTool.Log {
		chatOpsLog(channel, user, ret)
	}
	return true
}
func getUserChannelInfo(userid string, username string, channel string, timestamp string) string {
	return "export TRIGGERED_AT=" + timestamp + " && export TRIGGERED_USER_ID=" + userid + " && export TRIGGERED_USER_NAME=" + username + " && export TRIGGERED_CHANNEL_ID=" + channel + " && export TRIGGERED_CHANNEL_NAME=" + strings.Join(getChannelNames([]string{channel}), "")
}

// Raw commands
func processRawCommand(cmds []string, channel string, user string) bool {
	if stringInSlice(user, admin.UserIds) && channel == admin.PrivateChannelId {
		tmpCmd := html.UnescapeString(strings.Join(cmds, " "))
		// fmt.Println("Combined cmd: " + tmpCmd)
		tmpCmds := []string{"bash", "-c", tmpCmd}
		ret := shellOut(tmpCmds)
		yell(channel, ret)
		chatOpsLog(channel, user, ret)
		return true
	} else {
		reportToChannel(channel, "unauthorized", "")
		chatOpsLog(channel, user, strings.Join(cmds, " "))
		chatOpsLog(channel, user, "Unauthorized raw command!!!")
		return true
	}
}

func shellOut(cmdArgs []string) string {
	var (
		cmdOut []byte
		err    error
	)
	var cmdName string
	cmdName, cmdArgs = cmdArgs[0], cmdArgs[1:]
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		return "error running command."
	}
	out := string(cmdOut)
	log.Debug("Output from command:")
	condense := regexp.MustCompile(`\s*\n`)
	displayOut := condense.ReplaceAllString(out, "\\n")
	log.Debug(displayOut)
	return out
}
func reportToChannel(channel string, message string, passalong string) {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		log.Error(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var Messages Messages
	json.Unmarshal(byteValue, &Messages)

	isActive := true
	retMessage := message

	for i := 0; i < len(Messages.Messages); i++ {
		if Messages.Messages[i].Name == message {
			log.Debug(Messages.Messages[i].Name)
			isActive = Messages.Messages[i].Active
			if len(passalong) > 0 {
				retMessage = fmt.Sprintf(Messages.Messages[i].Text, passalong)
			} else {
				retMessage = Messages.Messages[i].Text
			}
		}
	}
	if isActive {
		log.Debug("Sending slack message[Channel:" + channel + "]: " + retMessage)
		yell(channel, retMessage)
	} else {
		log.Warn("Message suppressed by configuration")
		log.Warn(retMessage)
	}

}

func yell(channel string, msg string) {
	channelID, _, err := api.PostMessage(channel,
		slack.MsgOptionText(strings.Replace(msg, "\\n", "\n", -1), false),
		slack.MsgOptionUsername(admin.AppName),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}))

	if err != nil {
		log.Error(fmt.Printf("%s\n", err))
		return
	}
	log.Info("Send slack message[Channel:" + channelID + "]: " + msg)
}

func whisper(channel string, user string, msg string) {
	_, err := api.PostEphemeral(channel,
		user,
		slack.MsgOptionText(strings.Replace(msg, "\\n", "\n", -1), false),
		slack.MsgOptionUsername(admin.AppName),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}))

	if err != nil {
		log.Info(fmt.Printf("%s\n", err))
		return
	}
	log.Info("Send ephemeral slack message[Channel:" + channel + "]: " + msg)
}

func chatOpsLog(channel string, user string, msg string) {
	thisUser, err := api.GetUserInfo(user)
	if err != nil {
		log.Error("Couldn't get user")
		log.Error(err)
		return
	}
	thisChannel := getChannelNames([]string{channel})
	retacks := regexp.MustCompile("`")
	msg = retacks.ReplaceAllLiteralString(msg, "")
	ret := splitOut(admin.AppName+"["+thisUser.Profile.RealName+":"+thisChannel[0]+"]: "+truncateString(msg, 1000), "code")
	// Display message in chat-ops-log unless it came from admin channel
	if channel != admin.PrivateChannelId {
		channelID, _, err := api.PostMessage(admin.LogChannelId,
			slack.MsgOptionText(ret, false),
			slack.MsgOptionUsername(admin.AppName),
			slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
				UnfurlLinks: true,
				UnfurlMedia: true,
			}))

		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Channel ID: " + channelID)
	}
	log.Info(ret)
}

func splitOut(output string, responseType string) string {
	var splitInterval int = 4000
	switch responseType {
	case "code":
		if len(output) < splitInterval {
			return "```" + output + "```"
		}
		var splitChar = '\n'
		var splitCount int = 1
		var lastSplitPosition int = 0

		resultBuffer := bytes.Buffer{}

		for i, char := range output {
			if i >= (splitInterval*splitCount) && (char == splitChar) {
				resultBuffer.WriteString(
					strings.TrimLeft("```"+output[lastSplitPosition:i]+"``` \n", "\r\n"))
				lastSplitPosition = i + 1
				splitCount++
			}
		}
		log.Debug(resultBuffer.String())
		return resultBuffer.String()
	default:
		return output
	}
}

func handleMessage(event *slack.MessageEvent) {
	// log.Debug("handleMessage()")
	// log.Debug(event)
	// To Do: By bypassing this next check for a bot_message, we can test the bot's functionaltiy in a test slack channel
	if event.SubType == "bot_message" {
		return
	}

	for _, handler := range specials {
		if handler(event) {
			break
		}
	}
}

func initLog(logLevel string, logFormat string) {
	log.SetOutput(os.Stdout)

	switch logLevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.Warn(fmt.Sprintf("Invalid log-level (setting to info level): %s", logLevel))
	}

	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

}

func usage() {
	banner := ` ____            _     ____        _   
|  _ \          | |   |  _ \      | |  
| |_) | __ _ ___| |__ | |_) | ___ | |_ 
|  _ < / _' / __| '_ \|  _ < / _ \| __|
| |_) | (_| \__ \ | | | |_) | (_) | |_ 
|____/ \__,_|___/_| |_|____/ \___/ \__|
Bashbot is a slack bot, written in golang, that can be configured
to run bash commands or scripts based on a configuration file.
`
	fmt.Println(banner)
	fmt.Println("Usage: ./bashbot [flags]")
	fmt.Println("")
	flag.PrintDefaults()
}

func main() {
	flag.StringVar(&configFile, "config-file", "", "[REQUIRED] Filepath to config.json file (or environment variable BASHBOT_CONFIG_FILEPATH set)")
	flag.StringVar(&slackToken, "slack-token", "", "[REQUIRED] Slack token used to authenticate with api (or environment variable SLACK_TOKEN set)")
	flag.BoolVar(&installVendorDependenciesFlag, "install-vendor-dependencies", false, "Cycle through dependencies array in config file to install extra dependencies")
	flag.StringVar(&sendMessageChannel, "send-message-channel", "", "Send stand-alone slack message to this channel (requires -send-message-text)")
	flag.StringVar(&sendMessageText, "send-message-text", "", "Send stand-alone slack message (requires -send-message-channel)")
	flag.BoolVar(&sendMessageEphemeral, "send-message-ephemeral", false, "Send stand-alone ephemeral slack message to a specific user (requires -send-message-channel -send-message-text and -send-message-user)")
	flag.StringVar(&sendMessageUser, "send-message-user", "", "Send stand-alone ephemeral slack message to this slack user (requires -send-message-channel -send-message-text and -send-message-ephemeral)")
	flag.StringVar(&logLevel, "log-level", "info", "Log elevel to display (info,debug,warn,error)")
	flag.StringVar(&logFormat, "log-format", "text", "Display logs as json or text")
	flag.BoolVar(&help, "help", false, "Help/usage information")
	flag.BoolVar(&getVersion, "version", false, "Get current version")
	flag.Parse()
	if help {
		usage()
		os.Exit(0)
	}
	if getVersion {
		operatingSystem := runtime.GOOS
		systemArchitecture := runtime.GOARCH
		fmt.Println("bashbot-" + operatingSystem + "-" + systemArchitecture + "\t" + Version)
		os.Exit(0)
	}

	initLog(logLevel, logFormat)
	if configFile == "" && os.Getenv("BASHBOT_CONFIG_FILEPATH") != "" {
		configFile = os.Getenv("BASHBOT_CONFIG_FILEPATH")
	}
	if configFile == "" {
		usage()
		log.Error("Must define a config.json file")
		os.Exit(1)
	}
	if slackToken == "" && os.Getenv("SLACK_TOKEN") != "" {
		slackToken = os.Getenv("SLACK_TOKEN")
	}
	if slackToken == "" {
		usage()
		operatingSystem := runtime.GOOS
		systemArchitecture := runtime.GOARCH
		log.Error("Must define a slack token")
		log.Error("After logging into slack, visit https://api.slack.com/apps?new_classic_app=1")
		log.Error("to set up a new \"legacy bot user\" and \"Bot User OAuth Access Token\"")
		log.Error("Export the slack token as the environment variable SLACK_TOKEN")
		log.Error("export SLACK_TOKEN=xoxb-xxxxxxxxx-xxxxxxx")
		log.Error("bashbot-" + operatingSystem + "-" + systemArchitecture + " -config-file ./config.json -slack-token $SLACK_TOKEN")
		log.Error("See Read-me file for more detailed instructions: http://github.com/mathew-fleisch/bashbot")
		os.Exit(1)
	}

	if installVendorDependenciesFlag {
		if !installVendorDependencies() {
			log.Error("Failed to install dependencies")
			os.Exit(1)
		}
		os.Exit(0)
	}

	admin = getAdmin()
	api = slack.New(slackToken)

	// Send simple text message to slack
	if sendMessageChannel != "" && sendMessageText != "" {
		if sendMessageEphemeral && sendMessageUser != "" {
			whisper(sendMessageChannel, sendMessageUser, sendMessageText)
			os.Exit(0)
		}
		yell(sendMessageChannel, sendMessageText)
		os.Exit(0)
	}

	log.Info(admin.AppName + " Started: " + time.Now().String())

	var matchTrigger string = fmt.Sprintf("^%s .", admin.Trigger)

	// Regular expressions we'll use a whole lot.
	// Should probably be in an intialization function to the side.
	emojiPattern = regexp.MustCompile(`:[^\t\n\f\r ]+:`)
	slackUserPattern = regexp.MustCompile(`<@[^\t\n\f\r ]+>`)
	puncPattern = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cmdPattern = regexp.MustCompile(matchTrigger)

	// Our special handlers. If they handled a message, they return true.
	specials = []func(event *slack.MessageEvent) bool{processCommand}

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.Info("Bashbot is now connected to slack. Primary trigger: `" + admin.AppName + "`")

		case *slack.MessageEvent:
			handleMessage(ev)

		case *slack.PresenceChangeEvent:
			log.Info("Presence Change: " + ev.Presence)

		case *slack.RTMError:
			log.Error("Slack API RTM Error: " + ev.Error())

		case *slack.InvalidAuthEvent:
			log.Error("Invalid credentials (slack-token)")

		case *slack.ConnectionErrorEvent:
			log.Error("Can't connect to slack...")
			log.Error(msg)

		default:
			// Ignore other events..
			// log.Debug("Unhandled Event: " + msg.Type)
		}
	}
}
