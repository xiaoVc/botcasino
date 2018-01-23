package privatechat

import (
	"fmt"

	"github.com/zhangpanyi/basebot/history"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// RateBot 机器人评分
type RateBot struct {
}

// Handle 消息处理
func (rate *RateBot) Handle(bot *methods.BotExt, r *history.History, update *types.Update) {
	fromID := update.CallbackQuery.From.ID
	reply := fmt.Sprintf(tr(fromID, "lng_priv_rate_say"), bot.UserName, bot.UserName)
	menus := [...]methods.InlineKeyboardButton{
		methods.InlineKeyboardButton{
			Text:         tr(fromID, "lng_back_superior"),
			CallbackData: "/main/",
		},
	}
	markup := methods.MakeInlineKeyboardMarkupAuto(menus[:], 1)
	bot.EditMessageReplyMarkup(update.CallbackQuery.Message, reply, true, markup)
}

// 消息路由
func (rate *RateBot) route(bot *methods.BotExt, query *types.CallbackQuery) Handler {
	return nil
}
