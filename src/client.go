package src

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	core "github.com/guojia99/my-cubing-core"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/guojia99/my_cubing_robot/src/model"
	"github.com/guojia99/my_cubing_robot/src/process"
	"github.com/guojia99/my_cubing_robot/src/utils"
)

type Client struct {
	cfg  *Config
	e    *gin.Engine
	db   *gorm.DB
	core core.Core

	lock       sync.Mutex
	inCh       chan model.Message
	outCh      chan model.SendMessage
	processFns []process.ProcessFn
}

func NewClient(config string) (*Client, error) {
	c := &Client{
		e:          gin.Default(),
		inCh:       make(chan model.Message, 2048),
		outCh:      make(chan model.SendMessage, 255),
		processFns: process.ProcessDict,
	}
	if err := c.Load(config); err != nil {
		return nil, err
	}
	if err := c.initDB(); err != nil {
		return nil, err
	}
	c.core = core.NewCore(c.db, false, time.Minute*5)
	return c, nil
}

func (c *Client) Run() {
	c.e.NoRoute(func(ctx *gin.Context) {
		var r model.Message
		_ = ctx.Bind(&r)

		ctx.JSON(http.StatusOK, gin.H{})
		if r.GroupId == 0 || len(r.Message) == 0 || r.Message[0] != '*' {
			return
		}
		c.inCh <- r
	})

	go c.listenInputMessage()
	go c.listenOutPutMessage()
	_ = c.e.Run(fmt.Sprintf("127.0.0.1:%d", c.cfg.Port))
}

func (c *Client) initDB() error {
	var err error
	switch c.cfg.DB.Driver {
	case "sqlite":
		c.db, err = gorm.Open(sqlite.Open(c.cfg.DB.DSN), &gorm.Config{Logger: logger.Discard})
	case "mysql":
		c.db, err = gorm.Open(mysql.New(mysql.Config{DSN: c.cfg.DB.DSN}), &gorm.Config{Logger: logger.Discard})
	}
	return err
}

func (c *Client) listenInputMessage() {
	for {
		select {
		case data := <-c.inCh:
			msg := data.Message[1:]
			for _, fn := range c.processFns {
				ts := time.Now()
				if out := fn(c.db, c.core, msg, fmt.Sprintf("%d", data.UserId)); len(out) > 0 {

					useTime := fmt.Sprintf("\n(耗时: %s)", time.Now().Sub(ts).String())
					if time.Now().Sub(ts) > time.Second {
						out += useTime
					}

					c.outCh <- model.SendMessage{
						GroupId: data.GroupId,
						Message: out,
					}
					break
				}
			}
		}
	}
}

func (c *Client) listenOutPutMessage() {
	for {
		select {
		case data := <-c.outCh:
			err := c.sendMessage(data.GroupId, data.Message)
			if err == nil {
				time.Sleep(time.Second * 2)
				break
			}
		}
	}
}

func init() {
	//f, _ := os.Create("test.log")
	//log.SetOutput(f)
}

func (c *Client) sendMessage(groupId int, message string) error {
	if message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}

	if c.cfg.NotMessage {
		log.Printf("%s\n", message)
		return nil
	}

	_, err := utils.HTTPRequest("POST", "http://127.0.0.1:5700/send_group_msg", nil, nil, model.SendMessage{
		GroupId: groupId,
		Message: message,
	})
	return err
}
