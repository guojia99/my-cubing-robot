package qq_bot

import (
	"context"
	"log"
	"time"

	core "github.com/guojia99/my-cubing-core"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"

	"github.com/guojia99/my_cubing_robot/pkg/process"
)

func NewQQBotClient(conf Configs, db *gorm.DB) *QQBotClient {
	return &QQBotClient{
		db:         db,
		core:       core.NewCore(db, false, time.Second),
		conf:       conf,
		inputCh:    make(chan process.InMessage, 255),
		outputCh:   make(chan MessageToCreate, 255),
		imageCache: cache.New(time.Minute*5, time.Minute*5),
	}
}

type QQBotClient struct {
	ctx context.Context

	db   *gorm.DB
	core core.Core

	conf     Configs
	api      OpenAPI
	inputCh  chan process.InMessage
	outputCh chan MessageToCreate
	process  []process.Process

	imageCache *cache.Cache
}

func (q *QQBotClient) RegisterProcess(process ...process.Process) {
	q.process = append(q.process, process...)
}

func (q *QQBotClient) Run(ctx context.Context) error {
	q.ctx = ctx

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
		go q.doProcessLoop()
	}

	log.Printf("start qq bot")
	return NewSessionManager().Start(ws, tk, &intent)
}
