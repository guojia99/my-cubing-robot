package main

import (
	"github.com/spf13/cobra"

	"github.com/guojia99/my_cubing_robot/pkg/bot"
	"github.com/guojia99/my_cubing_robot/pkg/process"
)

func NewQQBotServerCmd() *cobra.Command {
	var config string
	cmd := &cobra.Command{
		Use:   "robot",
		Short: "魔方赛事系统Robot",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := bot.NewBots(config)
			if err != nil {
				return err
			}
			b.RegisterProcess(process.List...)
			return b.Run(cmd.Context())
		},
	}
	cmd.Flags().StringVarP(&config, "config", "c", "./etc/config.yml", "配置")
	return cmd
}

func main() {
	cmd := NewQQBotServerCmd()
	_ = cmd.Execute()
}
