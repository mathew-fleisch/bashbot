package slack

type Admin struct {
	Trigger          string   `yaml:"trigger"`
	AppName          string   `yaml:"appName"`
	UserIds          []string `yaml:"userIds"`
	PrivateChannelId string   `yaml:"privateChannelId"`
	LogChannelId     string   `yaml:"logChannelId"`
}

type Message struct {
	Active bool   `yaml:"active"`
	Name   string `yaml:"name"`
	Text   string `yaml:"text"`
}

type Tool struct {
	Name         string      `yaml:"name"`
	Description  string      `yaml:"description"`
	Help         string      `yaml:"help"`
	Trigger      string      `yaml:"trigger"`
	Location     string      `yaml:"location"`
	Command      []string    `yaml:"command"`
	Permissions  []string    `yaml:"permissions"`
	Log          bool        `yaml:"log"`
	Ephemeral    bool        `yaml:"ephemeral"`
	Response     string      `yaml:"response"`
	Parameters   []Parameter `yaml:"parameters"`
	Envvars      []string    `yaml:"envvars"`
	Dependencies []string    `yaml:"dependencies"`
}

type Parameter struct {
	Name        string   `yaml:"name"`
	Allowed     []string `yaml:"allowed"`
	Description string   `yaml:"description,omitempty"`
	Source      []string `yaml:"source,omitempty"`
	Match       string   `yaml:"match,omitempty"`
}

type Dependency struct {
	Name    string   `yaml:"name"`
	Install []string `yaml:"install"`
}

type Channel struct {
	Id                 string `yaml:"id"`
	Created            int    `yaml:"created"`
	IsOpen             bool   `yaml:"is_open"`
	IsGroup            bool   `yaml:"is_group"`
	IsShared           bool   `yaml:"is_shared"`
	IsIm               bool   `yaml:"is_im"`
	IsExtShared        bool   `yaml:"is_ext_shared"`
	IsOrgShared        bool   `yaml:"is_org_shared"`
	IsPendingExtShared bool   `yaml:"is_pending_ext_shared"`
	IsPrivate          bool   `yaml:"is_private"`
	IsMpim             bool   `yaml:"is_mpim"`
	Unlinked           int    `yaml:"unlinked"`
	NameNormalized     string `yaml:"name_normalized"`
	NumMembers         int    `yaml:"num_members"`
	Priority           int    `yaml:"priority"`
	User               string `yaml:"user"`
	Name               string `yaml:"name"`
	Creator            string `yaml:"creator"`
	IsArchived         bool   `yaml:"is_archived"`
	Members            string `yaml:"members"`
	Topic              Topic  `yaml:"topic"`
	Purpose            Topic  `yaml:"purpose"`
	IsChannel          bool   `yaml:"is_channel"`
	IsGeneral          bool   `yaml:"is_general"`
	IsMember           bool   `yaml:"is_member"`
	Local              string `yaml:"locale"`
}

type Topic struct {
	Value   string `yaml:"value"`
	Creator string `yaml:"creator"`
	LastSet int    `yaml:"last_set"`
}

type Config struct {
	Admins       []Admin      `yaml:"admins"`
	Messages     []Message    `yaml:"messages"`
	Tools        []Tool       `yaml:"tools"`
	Dependencies []Dependency `yaml:"dependencies"`
}

func (c *Config) GetTool(trigger string) Tool {
	for i := range c.Tools {
		if c.Tools[i].Trigger == trigger {
			return c.Tools[i]
		}
	}
	return Tool{}
}
