package bot

import (
	"context"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8s.io/klog"

	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot"
	"github.com/guojia99/my_cubing_robot/pkg/process"
)

type Bot interface {
	Run(ctx context.Context) error
	RegisterProcess(process ...process.Process)
}

type Bots struct {
	bots []Bot
}

func (b *Bots) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(len(b.bots))

	for _, bot := range b.bots {
		go func(bot Bot) {
			defer wg.Done()
			if err := bot.Run(ctx); err != nil {
				klog.Error(err)
				cancel()
			}
		}(bot)
	}

	wg.Wait()
	return nil
}

func (b *Bots) RegisterProcess(process ...process.Process) {
	for _, bot := range b.bots {
		bot.RegisterProcess(process...)
	}
}

func NewBots(cfgFile string) (Bot, error) {
	cfg, err := LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	var db *gorm.DB
	db, err = gorm.Open(mysql.New(mysql.Config{DSN: cfg.DSN}), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		return nil, err
	}

	var bots = &Bots{
		bots: make([]Bot, 0),
	}
	for _, val := range cfg.QQBot {
		bots.bots = append(bots.bots, qq_bot.NewQQBotClient(val, db))
	}

	return bots, nil
}
