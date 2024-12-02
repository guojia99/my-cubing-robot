package bot

import (
	cq_http "github.com/guojia99/my_cubing_robot/pkg/bot/cq-http"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot"
)

type Config struct {
	DSN    string           `json:"DSN" yaml:"dsn"`
	QQBot  []qq_bot.Configs `json:"QQBot" yaml:"QQBot"`
	CQHttp []cq_http.Config `json:"CQHttp" yaml:"CQHttp"`
}

func LoadConfig(file string) (*Config, error) {
	var out *Config
	body, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(body, &out)
	return out, err
}
