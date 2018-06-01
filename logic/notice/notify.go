package notice

import (
	"sync"

	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/botcasino/pusher"
)

var once sync.Once
var globalBot *methods.BotExt

// InitBotForOnce 初始化机器人
func InitBotForOnce(bot *methods.BotExt) {
	once.Do(func() {
		globalBot = bot
	})
}

// SendNotice 发送通知
func SendNotice(userID int64, message string) {
	if globalBot == nil {
		return
	}
	pusher.To(globalBot, userID, message, true, nil)
}
