package process_tools

import (
	"context"
	core "github.com/guojia99/my-cubing-core"
	"github.com/guojia99/my_cubing_robot/pkg/process"
	"gorm.io/gorm"
	"log"
)

type ProcessClient struct {
	Ctx context.Context

	Db       *gorm.DB
	Core     core.Core
	InputCh  chan process.InMessage
	OutputCh chan interface{}
	Process  []process.Process

	SendMsgFn process.SendEventHandler
}

func (p *ProcessClient) RegisterProcess(process ...process.Process) {
	p.Process = append(p.Process, process...)
}

func (p *ProcessClient) DoProcessLoop() {
	mp := process.PrefixMap(p.Process...)

	for {
		select {
		case <-p.Ctx.Done():
			return
		case msg := <-p.InputCh:
			log.Printf("input msg with `%s`, by `%s` send `%s`\n", msg.GroupID, msg.UserID, msg.Content)
			func() {
				ctx, cancel := context.WithCancel(p.Ctx)
				defer cancel()

				prs, err := process.CheckPrefixPro(msg.Content, mp)
				if err != nil {
					log.Printf("%s%s\n", msg.Content, err)
					return
				}

				if err = prs.Do(ctx, p.Db, p.Core, msg, p.SendMsgFn); err != nil {
					log.Printf("[debug] do process error %s\n", err)
				}
			}()

		}
	}
}
