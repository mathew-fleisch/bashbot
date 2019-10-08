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
  "unicode"

  "github.com/joho/godotenv"
  "github.com/nlopes/slack"
)

var specials []func(event *slack.MessageEvent) bool

// Slacking off with global vars
// "github.com/go-redis/redis"
// strip "github.com/grokify/html-strip-tags-go"
// var db *redis.Client
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

// func makeRedis() (r *redis.Client) {
//   address, found := os.LookupEnv("REDIS_ADDRESS")
//   if !found {
//     address = "127.0.0.1:6379"
//   }
//   log.Printf("using redis @ %s to store our data", address)
//   client := redis.NewClient(&redis.Options{Addr: address})
//   return client
// }

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
    yell(findChannelByName(address), "Witness the power of "+admin.AppName)
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
  // fmt.Printf("getChannelNames()")
  var names []string
  params := &slack.GetConversationsParameters{ExcludeArchived: "true", Limit: 10000, Types: []string{"private_channel", "public_channel"}}
  channels, _, _ := api.GetConversations(params)
  for j := 0; j < len(channels); j++ {
    thisChannel, _ := json.Marshal(channels[j])
    var thatChannel Channel
    json.Unmarshal([]byte(thisChannel), &thatChannel)
    // fmt.Printf("\n%s:%s", thatChannel.Id, thatChannel.Name)
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

  words := strings.Fields(event.Text)
  var triggered string
  var thisTool Tool
  var cmd []string
  // fmt.Println(words, len(words))

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
  // fmt.Println("Successfully Opened config.json")
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
    yell(event.Channel, ":bender: Processing command...")
    return processWhitelistedCommand(cmd, thisTool, event.Channel, event.User)
  case "cmd":
    yell(event.Channel, ":cat_typing: Processing raw command...")
    return processRawCommand(cmd, event.Channel, event.User)
  default:
    yell(event.Channel, ":thinkspin: Command not found...")
    return false

  }
}

// Whitelisted commands
func processWhitelistedCommand(cmds []string, thisTool Tool, channel string, user string) bool {
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
        // fmt.Println(cmdHelp)
        yell(channel, cmdHelp)
        return true
      }
    }
  }

  if authorized == false {
    unauthorized := ":redalert: You are not authorized to use this command in this channel.\nAllowed in: [" + strings.Join(allowedChannels, ", ") + "]"
    yell(channel, unauthorized)
    chatOpsLog(channel, user, thisTool.Trigger+" "+strings.Join(cmds, " "))
    chatOpsLog(channel, user, unauthorized)
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
      yell(channel, "Incorrect number of parameters")
      // fmt.Println(cmdHelp)
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
      yell(channel, "Invalid parameter value: "+thisTool.Parameters[x].Name)
      return false
    }
    re := regexp.MustCompile(`\${` + thisTool.Parameters[x].Name + `}`)
    buildCmd = re.ReplaceAllString(buildCmd, cmds[x])
    // fmt.Printf("%q\n", buildCmd)
  }
  buildCmd = getUserChannelInfo(user, thisUser.Profile.DisplayName, channel) + " && cd " + thisTool.Location + " && " + thisTool.Setup + " && " + buildCmd
  fmt.Printf("%q\n", buildCmd)

  tmpCmd := []string{"bash", "-c", buildCmd}
  // chatOpsLog(channel, user, strings.Join(tmpCmd, " "))

  ret := splitOut(shellOut(tmpCmd), thisTool.Response)

  if thisTool.Ephemeral == true {
    yell(channel, "Message only shown to user who triggered it.")
    whisper(channel, user, ret)
  } else {
    yell(channel, ret)
  }
  if thisTool.Log {
    chatOpsLog(channel, user, ret)
  }
  return true
}
func getUserChannelInfo(userid string, username string, channel string) string {
  return "export TRIGGERED_USER_ID=" + userid + " && export TRIGGERED_USER_NAME=" + username + " && export TRIGGERED_CHANNEL_ID=" + channel + " && export TRIGGERED_CHANNEL_NAME=" + strings.Join(getChannelNames([]string{channel}), "")
}
func stringInSlice(a string, list []string) bool {
  for _, b := range list {
    if b == a {
      return true
    }
  }
  return false
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
    yell(channel, ":redalert: Unauthorized!")
    chatOpsLog(channel, user, strings.Join(cmds, " "))
    chatOpsLog(channel, user, ":redalert: Unauthorized!")
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
    fmt.Fprintln(os.Stderr, "There be an error running that command: ", err)
    return ":skull_and_crossbones: There be an error running that command..."
  }
  out := string(cmdOut)
  fmt.Println(cmdName, out)
  return out
}

// func yourBasicShout(event *slack.MessageEvent) bool {
//   // log.Printf("event.Text: ")
//   // log.Printf(event.Text)
//   if !isLoud(event.Text) {
//     return false
//   }

//   // Your basic shout.
//   rejoinder, err := db.SRandMember(yellkey).Result()
//   if err != nil {
//     log.Printf("error selecting array: %s", err)
//     return false
//   }
//   yell(event.Channel, rejoinder)
//   // db.Incr(fmt.Sprintf("%s:count", countkey)).Result()
//   // db.SAdd(yellkey, event.Text).Result()
//   return true
// }

// End special handlers

func stripWhitespace(str string) string {
  var b strings.Builder
  b.Grow(len(str))
  for _, ch := range str {
    if !unicode.IsSpace(ch) {
      b.WriteRune(ch)
    }
  }
  return b.String()
}

// func isLoud(msg string) bool {
//   // log.Printf("msg: ")
//   // log.Printf(msg)
//   // strip tags & emoji
//   input := stripWhitespace(msg)
//   input = emojiPattern.ReplaceAllLiteralString(input, "")
//   input = slackUserPattern.ReplaceAllLiteralString(input, "")
//   input = puncPattern.ReplaceAllLiteralString(input, "")
//   input = strip.StripTags(input)
//   // log.Printf("input: ")
//   // log.Printf(input)
//   if len(input) < 3 {
//     return false
//   }
//   return strings.Contains(input, "pirate")
// }

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

func getAdmin() Admin {
  jsonFile, err := os.Open("admin.json")
  if err != nil {
    fmt.Println(err)
  }
  defer jsonFile.Close()
  byteValue, _ := ioutil.ReadAll(jsonFile)
  var Admins Admins
  json.Unmarshal(byteValue, &Admins)
  return Admins.Admins[0]
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
  var matchTrigger string = fmt.Sprintf("^%s", admin.Trigger)

  // rprefix, found := os.LookupEnv("REDIS_PREFIX")
  // if !found {
  //   rprefix = "PB"
  // }

  // yellkey = fmt.Sprintf("%s:YELLS", rprefix)
  // countkey = fmt.Sprintf("%s:COUNT", rprefix)

  // db = makeRedis()
  // card, err := db.SCard(yellkey).Result()
  if err != nil {
    // We fail NOW if we can't find our DB.
    log.Fatal(err)
  }
  // log.Printf(admin.AppName+" has %d things to say", card)

  // Regular expressions we'll use a whole lot.
  // Should probably be in an intialization function to the side.
  emojiPattern = regexp.MustCompile(`:[^\t\n\f\r ]+:`)
  slackUserPattern = regexp.MustCompile(`<@[^\t\n\f\r ]+>`)
  puncPattern = regexp.MustCompile(`[^a-zA-Z0-9]+`)
  cmdPattern = regexp.MustCompile(matchTrigger)

  // Our special handlers. If they handled a message, they return true.
  specials = []func(event *slack.MessageEvent) bool{processCommand}
  // processCommand,
  // yourBasicShout,
  // }

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
