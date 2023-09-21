package main

import (
	"github.com/spf13/cobra"

	"github.com/guojia99/my_cubing_robot/src"
)

func NewAPIServerCmd() *cobra.Command {
	var config string

	cmd := &cobra.Command{
		Use:   "robot",
		Short: "魔方赛事系统Robot",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := src.NewClient(config)
			if err != nil {
				return err
			}
			go c.Listen()
			return c.Run()
		},
	}
	cmd.Flags().StringVarP(&config, "config", "c", "./etc/configs.json", "配置")
	return cmd
}

func main() {
	cmd := NewAPIServerCmd()
	_ = cmd.Execute()
}
