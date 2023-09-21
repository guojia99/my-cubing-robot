package src

import (
	"os"

	json "github.com/json-iterator/go"
)

type Config struct {
	DB   DBConfig `json:"DB"`
	Port int      `json:"Port"`
}

type DBConfig struct {
	Driver string `json:"Driver"`
	DSN    string `json:"DSN"`
}

func (c *Client) Load(file string) error {
	configBody, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configBody, &c.cfg)
	return err
}
