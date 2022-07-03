package slack

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

type Config struct {
	Admins       []Admin      `json:"admins"`
	Messages     []Message    `json:"messages"`
	Tools        []Tool       `json:"tools"`
	Dependencies []Dependency `json:"dependencies"`
}

func (c *Config) GetTool(trigger string) Tool {
	for i := range c.Tools {
		if c.Tools[i].Trigger == trigger {
			return c.Tools[i]
		}
	}
	return Tool{}
}
