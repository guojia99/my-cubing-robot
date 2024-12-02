package qq_bot

import (
	"context"
	"github.com/guojia99/my_cubing_robot/pkg/bot/process_tools"
	"log"
	"time"

	core "github.com/guojia99/my-cubing-core"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/process"
)

func NewQQBotClient(conf Configs, db *gorm.DB) *QQBotClient {
	return &QQBotClient{
		conf:       conf,
		imageCache: cache.New(time.Minute*5, time.Minute*5),
		ProcessClient: process_tools.ProcessClient{
			Ctx:      nil,
			Db:       db,
			Core:     core.NewCore(db, false, time.Second),
			InputCh:  make(chan process.InMessage),
			OutputCh: make(chan interface{}),
		},
	}
}

type QQBotClient struct {
	process_tools.ProcessClient

	conf Configs
	api  OpenAPI

	imageCache *cache.Cache
}

func (q *QQBotClient) Run(ctx context.Context) error {
	q.Ctx = ctx
	q.SendMsgFn = q.sendMsgFn()

	SetLogger(logger)

	tk := BotToken(q.conf.AppID, q.conf.Token, "Bot")
	q.api = NewOpenAPI(tk).WithTimeout(10 * time.Second)
	ws, err := q.api.WS(ctx, nil, "")
	if err != nil {
		return err
	}

	var intent Intent
	if q.conf.Group {
		intent = RegisterHandlers(
			q._atGroupMessageEventHandler(),
			q._groupMessageEventHandler(),
		)
	} else {
		intent = RegisterHandlers(
			q._atMessageEventHandler(),
			q._messageEventHandler(),
		)
	}

	// look message input, and doing process detail
	for i := 0; i < 4; i++ {
		go q.DoProcessLoop()
	}

	log.Printf("start qq bot")
	return NewSessionManager().Start(ws, tk, &intent)
}
