package qq_bot

import (
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo"
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/dto"
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/event"
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/openapi"
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/token"
	"github.com/guojia99/my_cubing_robot/pkg/bot/qq-bot/botgo/websocket"
)

type MessageToCreate = dto.MessageToCreate
type OpenAPI = openapi.OpenAPI
type Intent = dto.Intent
type GroupMessageToCreate = dto.GroupMessageToCreate
type GroupRichMediaMessageToCreate = dto.GroupRichMediaMessageToCreate
type GroupAtMessageEventHandler = event.GroupAtMessageEventHandler
type GroupMessageEventHandler = event.GroupMessageEventHandler
type ATMessageEventHandler = event.ATMessageEventHandler
type MessageEventHandler = event.MessageEventHandler
type WSPayload = dto.WSPayload
type Media = dto.Media
type WSMessageData = dto.WSMessageData
type WSATMessageData = dto.WSATMessageData
type WSGroupMessageData = dto.WSGroupMessageData
type WSGroupATMessageData = dto.WSGroupATMessageData
type Message = dto.Message
type GroupMessage = dto.GroupMessage

var SetLogger = botgo.SetLogger
var BotToken = token.BotToken
var RegisterHandlers = websocket.RegisterHandlers
var NewSessionManager = botgo.NewSessionManager
var NewOpenAPI = botgo.NewOpenAPI
