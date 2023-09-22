package src

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guojia99/my-cubing/src/core"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	cfg  *Config
	e    *gin.Engine
	ch   chan Message
	db   *gorm.DB
	core core.Core

	processFns map[string]func(msg Message) error
}

func NewClient(config string) (*Client, error) {
	c := &Client{
		e:          gin.Default(),
		ch:         make(chan Message, 2048),
		processFns: make(map[string]func(msg Message) error),
	}
	if err := c.Load(config); err != nil {
		return nil, err
	}
	var err error
	switch c.cfg.DB.Driver {
	case "sqlite":
		c.db, err = gorm.Open(sqlite.Open(c.cfg.DB.DSN), &gorm.Config{})
	case "mysql":
		c.db, err = gorm.Open(mysql.New(mysql.Config{DSN: c.cfg.DB.DSN}), &gorm.Config{
			Logger: logger.Discard,
		})
	}
	c.initProcess()

	if err != nil {
		return nil, err
	}

	c.core = core.NewScoreCore(c.db, false)
	return c, nil
}

func (c *Client) Run() error {
	c.e.NoRoute(func(ctx *gin.Context) {
		var r Message
		_ = ctx.Bind(&r)
		ctx.JSON(http.StatusOK, gin.H{})
		if r.GroupId == 0 {
			return
		}
		c.setMessage(r)
	})
	return c.e.Run(fmt.Sprintf("127.0.0.1:%d", c.cfg.Port))
}

func (c *Client) Listen() {
	for {
		select {
		case data := <-c.ch:
			ts := time.Now()
			if err := c.messageProcess(data); err != nil {
				continue
			}
			if time.Now().Sub(ts) < time.Second*2 {
				time.Sleep(time.Now().Sub(ts))
			}
		}
	}
}
