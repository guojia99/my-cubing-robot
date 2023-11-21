package botgo

import (
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/log"
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/openapi"
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/websocket"
)

// SetLogger 设置 logger，需要实现 sdk 的 log.Logger 接口
func SetLogger(logger log.Logger) {
	log.DefaultLogger = logger
}

// SetSessionManager 注册自己实现的 session manager
func SetSessionManager(m SessionManager) {
	defaultSessionManager = m
}

// SetWebsocketClient 替换 websocket 实现
func SetWebsocketClient(c websocket.WebSocket) {
	websocket.Register(c)
}

// SetOpenAPIClient 注册 openapi 的不同实现，需要设置版本
func SetOpenAPIClient(v openapi.APIVersion, c openapi.OpenAPI) {
	openapi.Register(v, c)
}
