package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

var specials []func(event *slack.MessageEvent) bool

// Slacking off with global vars
var api *slack.Client
var rtm *slack.RTM
var channelsByName map[string]string

// var yellkey string
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

func makeChannelMap() {
	var admin Admin = getAdmin()
	log.Println("CONNECTED; ACQUIRING TARGETING DATA")
	channelsByName = make(map[string]string)
	channels, err := api.GetChannels(true)
	if err != nil {
		return
	}

	for _, v := range channels {
		channelsByName[v.Name] = v.ID
	}

	address, found := os.LookupEnv("WELCOME_CHANNEL")
	if found {
		reportToChannel(findChannelByName(address), "welcome", admin.AppName)
	}
	log.Println(admin.AppName + " IS NOW OPERATIONAL")
}

func findChannelByName(name string) string {
	// This feels unidiomatic.
	val, ok := channelsByName[name]
	if ok {
		return val
	}
	return ""
}

func getChannelNames(channelIds []string) []string {
	fmt.Printf("getChannelNames()")
	var names []string
	pparams := &slack.GetConversationsParameters{ExcludeArchived: "true", Limit: 1000, Types: []string{"private_channel"}}
	pchannels, _, _ := api.GetConversations(pparams)
	pnumChannels := len(pchannels)
	fmt.Printf("\nNumber of private channels:%s", pnumChannels)
	for j := 0; j < pnumChannels; j++ {
		pthisChannel, _ := json.Marshal(pchannels[j])
		var pthatChannel Channel
		json.Unmarshal([]byte(pthisChannel), &pthatChannel)
		// fmt.Printf("\n%s:%s", pthatChannel.Id, pthatChannel.Name)
		for i := 0; i < len(channelIds); i++ {
			// fmt.Printf("\ndoes %s match %s", pthatChannel.Id, channelIds[i])
			if channelIds[i] == pthatChannel.Id {
				names = append(names, pthatChannel.Name)
			}
		}
	}
	params := &slack.GetConversationsParameters{ExcludeArchived: "true", Limit: 1000, Types: []string{"public_channel"}}
	channels, _, _ := api.GetConversations(params)
	numChannels := len(channels)
	fmt.Printf("\nNumber of public channels:%s", numChannels)
	for j := 0; j < numChannels; j++ {
		thisChannel, _ := json.Marshal(channels[j])
		var thatChannel Channel
		json.Unmarshal([]byte(thisChannel), &thatChannel)
		// fmt.Printf("\n%s:%s", thatChannel.Id, thatChannel.Name)
		for i := 0; i < len(channelIds); i++ {
			// fmt.Printf("\ndoes %s match %s", thatChannel.Id, channelIds[i])
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
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var Admins Admins
	json.Unmarshal(byteValue, &Admins)
	return Admins.Admins[0]
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

	log.Println("command detected")
	log.Println(event.Text)
	fmt.Printf("%+v\n", event)
	fmt.Printf("Channel: %+v\n", event.Channel)
	fmt.Printf("User: %+v\n", event.User)
	fmt.Printf("Timestamp: %+v\n", event.Timestamp)

	words := strings.Fields(event.Text)
	var triggered string
	var thisTool Tool
	var cmd []string

	for index, element := range words {
		log.Printf(fmt.Sprintf("%s:%s", index, element))
		if index > 1 {
			cmd = append(cmd, element)
		}
	}

	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
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
	var admin Admin = getAdmin()
	validParams := make([]bool, len(thisTool.Parameters))
	var tmpHelp string
	authorized := false

	// inject email if exists in command
	reEmail := regexp.MustCompile(`\${email}`)
	thisUser, err := api.GetUserInfo(user)
	if err != nil {
		fmt.Printf("%s\n", err)
		return true
	}
	// fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", thisUser.ID, thisUser.Profile.RealName, thisUser.Profile.Email)
	thisTool.Command = reEmail.ReplaceAllLiteralString(thisTool.Command, thisUser.Profile.Email)

	// fmt.Println("Tool Name:        " + thisTool.Name)
	// fmt.Println("Tool Description: " + thisTool.Description)
	// fmt.Println("Tool Log:         " + strconv.FormatBool(thisTool.Log))
	// fmt.Println("Tool Help:        " + thisTool.Help)
	// fmt.Println("Tool Group:       " + thisTool.Group)
	// fmt.Println("Tool Trigger:     " + thisTool.Trigger)
	// fmt.Println("Tool Location:    " + thisTool.Location)
	// fmt.Println("Tool Setup:       " + thisTool.Setup)
	// fmt.Println("Tool Command:     " + thisTool.Command)
	// fmt.Println("Tool Ephemeral:   " + strconv.FormatBool(thisTool.Ephemeral))
	// fmt.Println("Tool Response:    " + thisTool.Response)
	var allowedChannels []string = getChannelNames(thisTool.Permissions)
	if admin.PrivateChannelId == channel {
		authorized = true
	} else {
		for j := 0; j < len(thisTool.Permissions); j++ {
			fmt.Println("Tool Permissions[" + strconv.Itoa(j) + "]: " + thisTool.Permissions[j])
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

	if authorized == false {
		reportToChannel(channel, "unauthorized", strings.Join(allowedChannels, ", "))
		yell(channel, cmdHelp)
		chatOpsLog(channel, user, thisTool.Trigger+" "+strings.Join(cmds, " "))
		return true
	}

	if len(thisTool.Parameters) > 0 {
		// fmt.Println("Tool Parameters Count: " + strconv.Itoa(len(thisTool.Parameters)))
		for j := 0; j < len(thisTool.Parameters); j++ {
			// fmt.Println("Tool Parameters[" + strconv.Itoa(j) + "]: " + thisTool.Parameters[j].Name)
			derivedSource := thisTool.Parameters[j].Source
			tmpHelp = fmt.Sprintf("%s\n%s: [%s%s]", tmpHelp, thisTool.Parameters[j].Name, strings.Join(thisTool.Parameters[j].Allowed, "|"), thisTool.Parameters[j].Description)
			if len(derivedSource) > 0 {
				// fmt.Println("No hard-coded allowed values. Deriving source: " + derivedSource)
				allowedOut := shellOut([]string{"bash", "-c", "cd " + thisTool.Location + " && " + derivedSource})
				// fmt.Println("Derived: " + allowedOut)
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
	if len(thisTool.Parameters) > 0 {
		// fmt.Println("Passed cmds: " + strconv.Itoa(len(cmds)) + " ?= " + strconv.Itoa(len(thisTool.Parameters)))
		if len(cmds) != len(thisTool.Parameters) {
			reportToChannel(channel, "incorrect_parameters", "")
			yell(channel, cmdHelp)
			return true
		}
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
		// fmt.Println("cmd[" + strconv.Itoa(x) + "]: " + cmds[x] + " -> " + strconv.FormatBool(validParams[x]))
		if validParams[x] == false {
			reportToChannel(channel, "invalid_parameter", thisTool.Parameters[x].Name)
			return false
		}
		re := regexp.MustCompile(`\${` + thisTool.Parameters[x].Name + `}`)
		buildCmd = re.ReplaceAllString(buildCmd, cmds[x])
		// fmt.Printf("%q\n", buildCmd)
	}
	buildCmd = getUserChannelInfo(user, thisUser.Name, channel, timestamp) + " && cd " + thisTool.Location + " && " + thisTool.Setup + " && " + buildCmd
	fmt.Printf("%q\n", buildCmd)

	tmpCmd := []string{"bash", "-c", buildCmd}

	ret := splitOut(shellOut(tmpCmd), thisTool.Response)

	if thisTool.Ephemeral == true {
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
	var admin Admin = getAdmin()
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
	fmt.Println(cmdName, out)
	return out
}
func reportToChannel(channel string, message string, passalong string) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var Messages Messages
	json.Unmarshal(byteValue, &Messages)

	isActive := true
	retMessage := message

	for i := 0; i < len(Messages.Messages); i++ {
		fmt.Printf(Messages.Messages[i].Name)
		if Messages.Messages[i].Name == message {
			isActive = Messages.Messages[i].Active
			if len(passalong) > 0 {
				retMessage = fmt.Sprintf(Messages.Messages[i].Text, passalong)
			} else {
				retMessage = Messages.Messages[i].Text
			}
		}
	}
	if isActive {
		yell(channel, retMessage)
	} else {
		log.Printf("Message suppressed by configuration:\n%s\n", retMessage)
	}

}

func yell(channel string, msg string) {
	var admin Admin = getAdmin()
	channelID, _, err := api.PostMessage(channel,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionUsername(admin.AppName),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}))

	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	log.Printf("Sent to %s: `%s`", channelID, msg)
}

func whisper(channel string, user string, msg string) {
	var admin Admin = getAdmin()
	_, err := api.PostEphemeral(channel,
		user,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionUsername(admin.AppName),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}))

	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	log.Printf("Sent to %s: `%s`", channel, msg)
}

func chatOpsLog(channel string, user string, msg string) {
	var admin Admin = getAdmin()
	thisUser, err := api.GetUserInfo(user)
	if err != nil {
		fmt.Printf("ERROR GETTING USER: %s\n", err)
		return
	}
	thisChannel := getChannelNames([]string{channel})
	// fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", thisUser.ID, thisUser.Profile.RealName, thisUser.Profile.Email)
	retacks := regexp.MustCompile("`")
	msg = retacks.ReplaceAllLiteralString(msg, "")
	ret := splitOut(admin.AppName+"["+thisUser.Profile.RealName+":"+thisChannel[0]+"]: "+truncateString(msg, 1000), "code")
	// Display message in chat-ops-log unless it came from matbots
	if channel != admin.PrivateChannelId {
		channelID, _, err := api.PostMessage(admin.LogChannelId,
			slack.MsgOptionText(ret, false),
			slack.MsgOptionUsername(admin.AppName),
			slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
				UnfurlLinks: true,
				UnfurlMedia: true,
			}))

		if err != nil {
			log.Printf("%s\n", err)
			return
		}
		log.Printf(channelID)
	}
	log.Printf(ret)
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
		// fmt.Printf("%s", resultBuffer.String())
		return resultBuffer.String()
	default:
		return output
	}
}

func handleMessage(event *slack.MessageEvent) {
	// fmt.Printf("Handle Event: %v\n", event)
	if event.SubType == "bot_message" {
		return
	}

	for _, handler := range specials {
		if handler(event) {
			break
		}
	}
}

func main() {
	var admin Admin = getAdmin()
	err := godotenv.Load(".env")
	log.Printf(admin.AppName+" Started: %s", time.Now())

	slacktoken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatal("You must provide an access token in SLACK_TOKEN")
	}

	api = slack.New(slacktoken)

	// log.Printf("Admin Config: %+v\n", admin)
	var matchTrigger string = fmt.Sprintf("^%s .", admin.Trigger)

	if err != nil {
		log.Fatal(err)
	}

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
			makeChannelMap()

		case *slack.MessageEvent:
			// fmt.Printf("Message: %v\n", ev)
			handleMessage(ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Printf(slacktoken)
			log.Fatal("Invalid credentials")

		case *slack.ConnectionErrorEvent:
			fmt.Printf("Event: %v\n", msg)
			log.Fatal("Can't connect")

		default:
			// Ignore other events..
			// fmt.Printf("Unhandled Event: %v\n", msg)
		}
	}
}
