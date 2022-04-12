package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// Config holds all bashbot's configurations
type Config struct {
	Admins       []Admin      `json:"admins"`
	Messages     []Message    `json:"messages"`
	Tools        []Tool       `json:"tools"`
	Dependencies []Dependency `json:"dependencies"`
}

type Admin struct {
	Trigger          string   `json:"trigger"`
	AppName          string   `json:"appName"`
	UserIds          []string `json:"userIds"`
	PrivateChannelId string   `json:"privateChannelId"`
	LogChannelId     string   `json:"logChannelId"`
}

type Message struct {
	Active bool   `json:"active"`
	Name   string `json:"name"`
	Text   string `json:"text"`
}

type Tool struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Help         string      `json:"help"`
	Trigger      string      `json:"trigger"`
	Location     string      `json:"location"`
	Command      []string    `json:"command"`
	Permissions  []string    `json:"permissions"`
	Log          bool        `json:"log"`
	Ephemeral    bool        `json:"ephemeral"`
	Response     string      `json:"response"`
	Parameters   []Parameter `json:"parameters"`
	Envvars      []string    `json:"envvars"`
	Dependencies []string    `json:"dependencies"`
}

type Parameter struct {
	Name        string   `json:"name"`
	Allowed     []string `json:"allowed"`
	Description string   `json:"description,omitempty"`
	Source      []string `json:"source,omitempty"`
	Match       string   `json:"match,omitempty"`
}

type Dependency struct {
	Name    string   `json:"name"`
	Install []string `json:"install"`
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

var (
	Version = "development"
)

// loadConfigFile is a helper function for loading bashbot json
// configuration file into Config struct.
func loadConfigFile(filePath string) (*Config, error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(fileContents, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func RunSlackApp(cfg *Config, botToken, appToken, appTrigger string) {
	api := slack.New(botToken, slack.OptionAppLevelToken(appToken))
	client := socketmode.New(api)
	go client.Run()

	for event := range client.Events {
		switch event.Type {
		case socketmode.EventTypeEventsAPI:
			eventsAPIHandler(cfg, client, event)

		case socketmode.EventTypeConnected:
			log.Info("Bashbot is now connected to slack. Primary trigger: `" + appTrigger + "`")

		case socketmode.EventTypeConnectionError:
			log.Error("Slack socket connection error")

		case socketmode.EventTypeErrorBadMessage:
			log.Error("Bad message received")
		}
	}
}

// eventsAPIHandler is a slack socket event handler for handling
// events API event.
func eventsAPIHandler(cfg *Config, client *socketmode.Client, socketEvent socketmode.Event) error {
	event := socketEvent.Data.(slackevents.EventsAPIEvent)
	client.Ack(*socketEvent.Request)
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent

		switch event := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			if event.SubType == "bot_message" {
				break
			}
			processCommand(cfg, client, event)
		}
	default:
		return fmt.Errorf("unhandled event type: %s", event.Type)
	}
	return nil
}

// installVendorDependencies is a helper function for installing the
// vendor dependencies required by the current bashbot instance.
//
// In the process of installing the dependencies, the dependency installer
// executes the install command provided in the configuration file for each
// dependency.
func installVendorDependencies(cfg *Config) bool {
	log.Debug("installing vendor dependencies")
	for i := 0; i < len(cfg.Dependencies); i++ {
		log.Info(cfg.Dependencies[i].Name)
		words := strings.Fields(strings.Join(cfg.Dependencies[i].Install, " "))
		var tcmd []string
		for index, element := range words {
			log.Debugf("%d: %s", index, element)
			tcmd = append(tcmd, element)
		}
		cmd := []string{"bash", "-c", "pushd vendor && " + strings.Join(tcmd, " ") + " && popd"}
		log.Debug(strings.Join(cmd, " "))
		log.Info(runShellCommands(cmd))
	}
	return true
}

// runShellCommands is a helper function for executing shell commands on
// the bashbot host machine.
//
//  usage:
// 		runShellCommands([]string{"bash", "-c", "apt-get install git && echo hello"})
//
// The first value in the array should be the command name e.g bash, sh etc
// while the other values will be treated as arguments.
func runShellCommands(cmdArgs []string) string {
	cmdOut, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).CombinedOutput()
	if err != nil {
		return "error running command."
	}
	out := string(cmdOut)
	displayOut := regexp.MustCompile(`\s*\n`).ReplaceAllString(out, "\\n")
	log.Debug("Output from command: \n", displayOut)
	return out
}

// sendConfigMessageToChannel sends a message to the slack channel based on the
// messages configured in the bashbot config file.
//
// usage:
// 		sendConfigMessageToChannel(cfg, client, "channelID", "processing_command", "try another command")
//
// The passalong parameter is an optional parameter because not all messages needs additional
// content(s) attached to the message sent.
func sendConfigMessageToChannel(cfg *Config, client *socketmode.Client, channel, message, passalong string) {
	isActive := true
	responseMessage := message
	for i := 0; i < len(cfg.Messages); i++ {
		if cfg.Messages[i].Name == message {
			log.Debug(cfg.Messages[i].Name)
			isActive = cfg.Messages[i].Active
			responseMessage = cfg.Messages[i].Text
			if passalong != "" {
				responseMessage = fmt.Sprintf(cfg.Messages[i].Text, passalong)
			}
		}
	}
	if isActive {
		sendMessageToChannel(cfg, client, channel, responseMessage)
		return
	}
	log.Warn("Message suppressed by configuration")
	log.Warn(responseMessage)
}

// sendMessageToChannel sends a message to the slack channel.
func sendMessageToChannel(cfg *Config, client *socketmode.Client, channel, msg string) {
	channelID, _, err := client.PostMessage(
		channel,
		slack.MsgOptionText(strings.Replace(msg, "\\n", "\n", -1), false),
		slack.MsgOptionUsername(cfg.Admins[0].AppName),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}),
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Sent slack message[Channel:%s]: %s", channelID, msg)
}

// sendMessageToUser sends to message to a slack user in a slack channel.
func sendMessageToUser(cfg *Config, client *socketmode.Client, channel, user, msg string) {
	_, err := client.PostEphemeral(channel,
		user,
		slack.MsgOptionText(strings.Replace(msg, "\\n", "\n", -1), false),
		slack.MsgOptionUsername(cfg.Admins[0].AppName),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}),
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Sent ephemeral slack message[Channel:" + channel + "]: " + msg)
}

func truncateString(str string, num int) string {
	res := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		res = str[0:num] + "..."
	}
	return res
}

// getChannelNames retreives the names of the channels monitored by bashbot
// using the channels id.
func getChannelNames(client *socketmode.Client, channelsID []string) []string {
	var names []string
	// extracting the private channels monitored by bashbot.
	channels, _, _ := client.GetConversations(&slack.GetConversationsParameters{
		Limit: 1000,
		Types: []string{"private_channel"},
	})
	log.Debugf("Number of private channels this bot is monitoring: %d", len(channels))
	for j := 0; j < len(channels); j++ {
		for i := range channelsID {
			if channelsID[i] == channels[j].ID {
				names = append(names, channels[j].Name)
			}
		}
	}
	// extracting the public channels monitored by bashbot.
	channels, _, _ = client.GetConversations(&slack.GetConversationsParameters{
		Limit: 1000,
		Types: []string{"public_channel"},
	})
	log.Debugf("Number of public channels this bot is monitoring: %d", len(channels))
	for j := 0; j < len(channels); j++ {
		for i := range channelsID {
			if channelsID[i] == channels[j].ID {
				names = append(names, channels[j].Name)
			}
		}
	}
	if len(names) > 0 {
		return names
	}
	return []string{"all"}
}

func processCommand(cfg *Config, client *socketmode.Client, event *slackevents.MessageEvent) bool {
	matchTrigger := fmt.Sprintf("(?i)^%s .", cfg.Admins[0].Trigger)
	cmdPattern := regexp.MustCompile(matchTrigger)
	if !cmdPattern.MatchString(event.Text) {
		return false
	}

	log.Info("command detected: `" + event.Text + "`")
	log.Debug(event)
	log.Info("Channel: " + event.Channel)
	log.Info("User: " + event.User)
	log.Info("Timestamp: " + event.TimeStamp)

	words := strings.Fields(event.Text)
	cmd := []string{}
	for index, element := range words {
		element = regexp.MustCompile(`<http(.*)>`).ReplaceAllString(element, "http$1")
		element = regexp.MustCompile(`“|”`).ReplaceAllString(element, "\"")
		element = regexp.MustCompile(`‘|’`).ReplaceAllString(element, "'")
		log.Infof("%d: %s", index, element)
		if index > 1 {
			cmd = append(cmd, element)
		}
	}

	tool := Tool{}
	for i := range cfg.Tools {
		if cfg.Tools[i].Trigger == words[1] {
			tool = cfg.Tools[i]
		}
	}

	switch words[1] {
	case tool.Trigger:
		sendConfigMessageToChannel(cfg, client, event.Channel, "processing_command", "")
		return processValidCommand(cfg, client, cmd, tool, event.Channel, event.User, event.TimeStamp)
	case "exit":
		if len(words) == 3 {
			switch words[2] {
			case "0":
				sendMessageToChannel(cfg, client, event.Channel, "exiting: success")
				os.Exit(0)
			default:
				sendMessageToChannel(cfg, client, event.Channel, "exiting: failure")
				os.Exit(1)
			}
		}
		sendMessageToChannel(cfg, client, event.Channel, "My battery is low and it's getting dark.")
		os.Exit(0)
		return false
	default:
		sendConfigMessageToChannel(cfg, client, event.Channel, "command_not_found", "")
		return false
	}
}

func processValidCommand(cfg *Config, client *socketmode.Client, cmds []string, tool Tool, channel, user, timestamp string) bool {
	// checking if all required environment variables exist.
	for _, envvar := range tool.Envvars {
		if os.Getenv(envvar) == "" {
			sendConfigMessageToChannel(cfg, client, channel, "missingenvvar", envvar)
			return false
		}
	}
	// checking if all required dependencies are installed on the
	// machine / os.
	for _, dependency := range tool.Dependencies {
		if _, err := exec.LookPath(dependency); err != nil {
			sendConfigMessageToChannel(cfg, client, channel, "missingdependency", dependency)
			return false
		}
	}
	// inject email if exists in command
	thisUser, err := client.GetUserInfo(user)
	if err != nil {
		log.Info(fmt.Printf("%s\n", err))
		return true
	}
	reEmail := regexp.MustCompile(`\${email}`)
	commandJoined := reEmail.ReplaceAllLiteralString(strings.Join(tool.Command, " "), thisUser.Profile.Email)

	log.Debugf(" ----> Param Name:        %s", tool.Name)
	log.Debugf(" ----> Param Description: %s", tool.Description)
	log.Debugf(" ----> Param Log:         %s", strconv.FormatBool(tool.Log))
	log.Debugf(" ----> Param Help:        %s", tool.Help)
	log.Debugf(" ----> Param Trigger:     %s", tool.Trigger)
	log.Debugf(" ----> Param Location:    %s", tool.Location)
	log.Debugf(" ----> Param Command:     %s", commandJoined)
	log.Debugf(" ----> Param Ephemeral:   %s", strconv.FormatBool(tool.Ephemeral))
	log.Debugf(" ----> Param Response:    %s", tool.Response)
	validParams := make([]bool, len(tool.Parameters))
	var tmpHelp string
	authorized := false
	var allowedChannels []string = getChannelNames(client, tool.Permissions)
	if cfg.Admins[0].PrivateChannelId == channel {
		authorized = true
	} else {
		for j := 0; j < len(tool.Permissions); j++ {
			log.Debugf(" ----> Param Permissions[%d]: %s", j, tool.Permissions[j])
			if tool.Permissions[j] == channel || tool.Permissions[j] == "all" {
				authorized = true
			}
		}
	}

	// Show help if the first parameter is "help"
	cmdHelp := fmt.Sprintf("``` ====> %s [Allowed In: %s] <====\n%s\n%s%s```", tool.Name, strings.Join(allowedChannels, ", "), tool.Description, tool.Help, tmpHelp)
	if len(cmds) > 0 {
		for j := 0; j < len(cmds); j++ {
			if cmds[j] == "help" {
				sendMessageToChannel(cfg, client, channel, cmdHelp)
				return true
			}
		}
	}

	if !authorized {
		sendConfigMessageToChannel(cfg, client, channel, "unauthorized", strings.Join(allowedChannels, ", "))
		sendMessageToChannel(cfg, client, channel, cmdHelp)
		logToChannel(cfg, client, channel, user, tool.Trigger+" "+strings.Join(cmds, " "))
		return true
	}

	if len(tool.Parameters) > 0 {
		log.Debug(" ----> Param Parameters Count: " + strconv.Itoa(len(tool.Parameters)))
		for j := range tool.Parameters {
			log.Debug(" ----> Param Parameters[" + strconv.Itoa(j) + "]: " + tool.Parameters[j].Name)
			derivedSource := tool.Parameters[j].Source
			tmpHelp = fmt.Sprintf("%s\n%s: [%s%s]", tmpHelp, tool.Parameters[j].Name, strings.Join(tool.Parameters[j].Allowed, "|"), tool.Parameters[j].Description)
			if len(derivedSource) > 0 {
				log.Debug("Deriving allowed parameters: " + strings.Join(derivedSource, " "))
				allowedOut := strings.Split(runShellCommands([]string{"bash", "-c", "cd " + tool.Location + " && " + strings.Join(derivedSource, " ")}), "\n")
				tool.Parameters[j].Allowed = append(tool.Parameters[j].Allowed, allowedOut...)
			}
		}
	}

	if tool.Log {
		logToChannel(cfg, client, channel, user, tool.Trigger+" "+strings.Join(cmds, " "))
	}

	// Validate parameters against whitelist
	if len(tool.Parameters) > 0 {
		for j := 0; j < len(tool.Parameters); j++ {
			log.Debug(" ====> Param Name: " + tool.Parameters[j].Name)
			validParams[j] = false

			if len(tool.Parameters[j].Match) > 0 {
				log.Debug(" ====> Parameter[" + strconv.Itoa(j) + "].Regex: " + tool.Parameters[j].Match)
				restOfCommand := strings.Join(cmds[j:], " ")
				if regexp.MustCompile(tool.Parameters[j].Match).MatchString(restOfCommand) {
					log.Debug("Parameter(s): '" + restOfCommand + "' matches regex: '" + tool.Parameters[j].Match + "'")
					validParams[j] = true
				} else {
					log.Debug("Parameter: " + cmds[j] + " does not match regex: " + tool.Parameters[j].Match)
				}
			} else {
				for h := 0; h < len(tool.Parameters[j].Allowed); h++ {
					log.Debug(" ====> Parameter[" + strconv.Itoa(j) + "].Allowed[" + strconv.Itoa(h) + "]: " + tool.Parameters[j].Allowed[h])
					if tool.Parameters[j].Allowed[h] == cmds[j] {
						validParams[j] = true
					}
				}
			}
		}
	}

	buildCmd := commandJoined
	for x := 0; x < len(tool.Parameters); x++ {
		if !validParams[x] {
			sendConfigMessageToChannel(cfg, client, channel, "invalid_parameter", tool.Parameters[x].Name)
			return false
		}
		re := regexp.MustCompile(`\${` + tool.Parameters[x].Name + `}`)
		if len(tool.Parameters[x].Match) > 0 {
			buildCmd = re.ReplaceAllString(buildCmd, strings.Join(cmds[x:], " "))
		} else {
			buildCmd = re.ReplaceAllString(buildCmd, cmds[x])
		}
	}
	buildCmd = fmt.Sprintf(
		"export TRIGGERED_AT=%s && export TRIGGERED_USER_ID=%s && export TRIGGERED_USER_NAME=%s && export TRIGGERED_CHANNEL_ID=%s && export TRIGGERED_CHANNEL_NAME=%s",
		timestamp,
		user,
		thisUser.Name,
		channel,
		strings.Join(getChannelNames(client, []string{channel}), ""),
	)
	splitOn := regexp.MustCompile(`\s\&\&`)
	displayCmd := splitOn.ReplaceAllString(buildCmd, " \\\n        &&")
	log.Info("Triggered Command:")
	log.Info(displayCmd)

	tmpCmd := []string{"bash", "-c", buildCmd}
	ret := splitOut(runShellCommands(tmpCmd), tool.Response)
	if tool.Ephemeral {
		sendConfigMessageToChannel(cfg, client, channel, "ephemeral", "")
		sendMessageToUser(cfg, client, channel, user, ret)
	} else {
		sendMessageToChannel(cfg, client, channel, ret)
	}
	if tool.Log {
		logToChannel(cfg, client, channel, user, ret)
	}
	return true
}

func logToChannel(cfg *Config, client *socketmode.Client, channelID, userID, msg string) {
	user, err := client.GetUserInfo(userID)
	if err != nil {
		log.Errorf("can't get user: %w", err)
		return
	}
	channel := getChannelNames(client, []string{channelID})
	retacks := regexp.MustCompile("`")
	msg = retacks.ReplaceAllLiteralString(msg, "")
	output := fmt.Sprintf(
		"%s[%s:%s]: %s",
		cfg.Admins[0].AppName,
		user.Profile.RealName,
		channel[0],
		truncateString(msg, 1000),
	)
	ret := splitOut(output, "code")
	// Display message in chat-ops-log unless it came from admin channel
	if channelID == cfg.Admins[0].PrivateChannelId {
		return
	}
	channelID, _, err = client.PostMessage(cfg.Admins[0].LogChannelId,
		slack.MsgOptionText(ret, false),
		slack.MsgOptionUsername(cfg.Admins[0].AppName),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}),
	)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Channel ID: " + channelID)
	log.Info(ret)
}

func splitOut(output string, responseType string) string {
	splitInterval := 4000
	switch responseType {
	case "code":
		if len(output) < splitInterval {
			return fmt.Sprintf("```%s```", output)
		}
		splitCount := 1
		lastSplitPosition := 0
		resultBuffer := bytes.Buffer{}
		for i, char := range output {
			if i >= (splitInterval*splitCount) && (char == '\n') {
				resultBuffer.WriteString(strings.TrimLeft("```"+output[lastSplitPosition:i]+"``` \n", "\r\n"))
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
