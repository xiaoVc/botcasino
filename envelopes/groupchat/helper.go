package groupchat

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/zhangpanyi/basebot/logger"
	tg "github.com/zhangpanyi/basebot/telegram"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/botcasino/config"
	"github.com/zhangpanyi/botcasino/storage"
)

// 获取广告
func getAd(botID int64) string {
	handler := storage.AdStorage{}
	ads, err := handler.GetAds(botID)
	if err != nil || len(ads) == 0 {
		return ""
	}
	ad := ads[randx.Intn(len(ads))]
	return fmt.Sprintf("\n\n*[* %s *]*", tg.Pre(ad.Text))
}

// 语言翻译
func tr(userID int64, key string) string {
	return config.GetLanguge().Value("zh_CN", key)
}

// SendRedEnvelopeToGroup 发送红包到群组
func SendRedEnvelopeToGroup(bot *methods.BotExt, userID, chatID int64, id uint64) error {
	// 获取红包信息
	newHandler := storage.RedEnvelopeStorage{}
	redEnvelope, received, err := newHandler.GetRedEnvelope(id)
	if err != nil {
		logger.Errorf("Failed to get red envelope, %v", err)
		return err
	}

	// 红包身份验证
	if userID != redEnvelope.SenderID {
		return errors.New("auth failed")
	}

	// 检查重复激活
	if redEnvelope.Active {
		return errors.New("auth failed")
	}

	// 检查红包过期
	now := time.Now().UTC().Unix()
	dynamicCfg := config.GetDynamic()
	if now-redEnvelope.Timestamp >= dynamicCfg.RedEnvelopeExpire {
		return errors.New("already activated")
	}

	// 生成菜单列表
	data := strconv.FormatUint(redEnvelope.ID, 10)
	menus := [...]methods.InlineKeyboardButton{
		methods.InlineKeyboardButton{Text: tr(0, "lng_chat_receive"), CallbackData: data},
	}

	// 回复红包信息
	reply := tr(0, "lng_chat_welcome")
	typ := redEnvelopesTypeToString(redEnvelope.Lucky)
	amount := fmt.Sprintf("%.2f", float64(redEnvelope.Amount)/100.0)
	if !redEnvelope.Lucky {
		amount = fmt.Sprintf("%.2f", float64(redEnvelope.Amount*redEnvelope.Number)/100.0)
	}
	reply = fmt.Sprintf(reply, typ, received, redEnvelope.Number, amount,
		storage.GetAsset(redEnvelope.Asset), redEnvelope.SenderName,
		redEnvelope.Memo, getAd(bot.ID), bot.UserName, redEnvelope.ID, bot.UserName, redEnvelope.ID)
	markup := methods.MakeInlineKeyboardMarkup(menus[:], 1)
	message, err := bot.SendMessageDisableWebPagePreview(chatID, reply, true, markup)
	if err != nil {
		logger.Errorf("Failed to send red envelope info, %v", err)
		return err
	}

	// 激活红包
	err = newHandler.ActiveRedEnvelope(id, userID, message.Chat.ID, message.MessageID)
	if err != nil {
		logger.Errorf("Failed to active red envelope, %v", err)
		return err
	}
	return nil
}
