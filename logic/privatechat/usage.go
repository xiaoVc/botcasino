package privatechat

import (
	"github.com/zhangpanyi/basebot/history"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// UsageHandler 使用说明
type UsageHandler struct {
}

// Handle 消息处理
func (*UsageHandler) Handle(bot *methods.BotExt, r *history.History, update *types.Update) {
	fromID := update.CallbackQuery.From.ID
	menus := [...]methods.InlineKeyboardButton{
		methods.InlineKeyboardButton{
			Text:         tr(fromID, "lng_back_superior"),
			CallbackData: "/main/",
		},
	}
	markup := methods.MakeInlineKeyboardMarkupAuto(menus[:], 1)

	reply := tr(fromID, "lng_priv_usage")
	bot.AnswerCallbackQuery(update.CallbackQuery, "", false, "", 0)
	bot.EditMessageReplyMarkupDisableWebPagePreview(update.CallbackQuery.Message, reply, true, markup)
}

// 消息路由
func (*UsageHandler) route(bot *methods.BotExt, query *types.CallbackQuery) Handler {
	return nil
}
