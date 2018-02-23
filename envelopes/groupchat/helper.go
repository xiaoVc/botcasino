package groupchat

import (
	"fmt"

	"github.com/zhangpanyi/botcasino/config"
	"github.com/zhangpanyi/botcasino/storage"
	tg "github.com/zhangpanyi/basebot/telegram"
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
func tr(userID int32, key string) string {
	return config.GetLanguge().Value("zh_CN", key)
}
