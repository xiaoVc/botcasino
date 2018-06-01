package privatechat

import (
	"fmt"
	"strings"

	"github.com/zhangpanyi/botcasino/storage"

	"github.com/zhangpanyi/basebot/history"
	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// Handler 消息处理器
type Handler interface {
	route(*methods.BotExt, *types.CallbackQuery) Handler
	Handle(*methods.BotExt, *history.History, *types.Update)
}

// MainMenuHandler 主菜单
type MainMenuHandler struct {
}

// Handle 消息处理
func (handler *MainMenuHandler) Handle(bot *methods.BotExt, r *history.History, update *types.Update) {
	if bot == nil || r == nil {
		return
	}

	// 处理消息
	if update.Message != nil {
		// 是否由子菜单处理
		var callback *types.Update
		r.Foreach(func(idx int, element *types.Update) bool {
			if element.CallbackQuery != nil {
				callback = element
				return false
			}
			return true
		})

		// 子菜单处理请求
		if update.Message.Text != "/start" && callback != nil {
			newHandler := handler.route(bot, callback.CallbackQuery)
			if newHandler == nil {
				r.Clear()
				return
			}
			newHandler.Handle(bot, r.Push(update), callback)
			return
		}

		// 发送菜单列表
		reply, menus := handler.replyMessage(update.Message.From.ID)
		markup := methods.MakeInlineKeyboardMarkup(menus, 2, 2, 2, 1)
		bot.SendMessage(update.Message.Chat.ID, reply, true, markup)
		r.Clear()
		return
	}

	if update.CallbackQuery == nil {
		return
	}

	// 回复主菜单
	if update.CallbackQuery.Data == "/main/" {
		bot.AnswerCallbackQuery(update.CallbackQuery, "", false, "", 0)
		reply, menus := handler.replyMessage(update.CallbackQuery.From.ID)
		markup := methods.MakeInlineKeyboardMarkup(menus, 2, 2, 2, 1)
		bot.EditMessageReplyMarkup(update.CallbackQuery.Message, reply, true, markup)
		return
	}

	// 路由到其它处理模块
	newHandler := handler.route(bot, update.CallbackQuery)
	if newHandler == nil {
		return
	}
	newHandler.Handle(bot, r, update)
}

// 消息路由
func (handler *MainMenuHandler) route(bot *methods.BotExt, query *types.CallbackQuery) Handler {
	// 发放红包
	if strings.HasPrefix(query.Data, "/give/") {
		return new(GiveHandler)
	}

	// 选择群组
	if strings.HasPrefix(query.Data, "/chatid/") {
		return new(SelectGroupHandler)
	}

	// 使用说明
	if strings.HasPrefix(query.Data, "/usage/") {
		return new(UsageHandler)
	}

	// 机器人评分
	if strings.HasPrefix(query.Data, "/rate/") {
		return new(RateBotHandler)
	}

	// 分享机器人
	if strings.HasPrefix(query.Data, "/share/") {
		return new(ShareBotHandler)
	}

	// 操作历史记录
	if strings.HasPrefix(query.Data, "/history/") {
		return new(HistoryHandler)
	}

	// 存款锚定资产
	if strings.HasPrefix(query.Data, "/deposit/") {
		return new(DepositHandler)
	}

	// 提现锚定资产
	if strings.HasPrefix(query.Data, "/withdraw/") {
		return new(WithdrawHandler)
	}
	return nil
}

// 获取用户资产数量
func getUserAssetAmount(userID int64, asset string) string {
	newHandler := storage.AssetStorage{}
	assetInfo, err := newHandler.GetAsset(userID, asset)
	if err != nil {
		if err != storage.ErrNoBucket && err != storage.ErrNoSuchTypeAsset {
			logger.Warnf("Failed to get user asset, %v, %v, %v", userID, asset, err)
		}
		return "0.00"
	}
	return fmt.Sprintf("%.2f", float64(assetInfo.Amount)/100.0)
}

// 获取回复消息
func (handler *MainMenuHandler) replyMessage(userID int64) (string, []methods.InlineKeyboardButton) {
	// 获取资产信息
	bitCNY := getUserAssetAmount(userID, storage.BitCNYSymbol)
	bitUSD := getUserAssetAmount(userID, storage.BitUSDSymbol)

	// 生成菜单列表
	menus := [...]methods.InlineKeyboardButton{
		methods.InlineKeyboardButton{Text: tr(userID, "lng_priv_give_lucky_money"), CallbackData: "/give/"},
		methods.InlineKeyboardButton{Text: tr(userID, "lng_priv_history"), CallbackData: "/history/"},
		methods.InlineKeyboardButton{Text: tr(userID, "lng_priv_deposit"), CallbackData: "/deposit/"},
		methods.InlineKeyboardButton{Text: tr(userID, "lng_priv_withdraw"), CallbackData: "/withdraw/"},
		methods.InlineKeyboardButton{Text: tr(userID, "lng_priv_rate"), CallbackData: "/rate/"},
		methods.InlineKeyboardButton{Text: tr(userID, "lng_priv_share"), CallbackData: "/share/"},
		methods.InlineKeyboardButton{Text: tr(userID, "lng_priv_help"), CallbackData: "/usage/"},
	}
	return fmt.Sprintf(tr(userID, "lng_priv_welcome"), bitCNY, bitUSD), menus[:]
}
