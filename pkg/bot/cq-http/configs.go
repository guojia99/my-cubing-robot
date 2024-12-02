package cq_http

type Config struct {
	Enable   bool   `json:"enable" yaml:"enable"`
	Address  string `json:"address" yaml:"address"`
	Port     int    `json:"port" yaml:"port"`
	SendPort int    `json:"send_port" yaml:"send_port"`
}
