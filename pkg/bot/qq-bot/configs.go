package qq_bot

type Configs struct {
	Enable    bool     `json:"enable" yaml:"enable"`
	AppID     uint64   `json:"app_id" yaml:"appID"`
	Token     string   `json:"token" yaml:"token"`
	Group     bool     `json:"group" yaml:"group"`
	GroupList []string `json:"groupList" yaml:"groupList"`
}
